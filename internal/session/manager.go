// Package session управляет сессией Minecraft (handshake, login, keepalive).
package session

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nithdevv/goflayer/internal/conn"
	"github.com/nithdevv/goflayer/internal/logger"
	"github.com/nithdevv/goflayer/internal/protocol"
	"github.com/nithdevv/goflayer/pkg/events"
)

// Manager manages the Minecraft session.
type Manager struct {
	mu       sync.RWMutex
	state    protocol.State
	conn     *conn.Conn
	events   *events.Bus
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup

	// Configuration
	username string
	protocol int

	// Logger
	log *logger.Logger
}

// New creates a new session manager.
func New(c *conn.Conn, username string, protocolVer int, ev *events.Bus) *Manager {
	log := logger.Default().With("session")
	ctx, cancel := context.WithCancel(context.Background())

	return &Manager{
		state:   protocol.Handshaking,
		conn:    c,
		events:  ev,
		ctx:     ctx,
		cancel:  cancel,
		username: username,
		protocol: protocolVer,
		log:     log,
	}
}

// Start starts the session (handshake + login).
func (m *Manager) Start(serverHost string, serverPort uint16) error {
	m.log.Info("Starting session...")

	// Handshake
	if err := m.handshake(serverHost, serverPort); err != nil {
		return fmt.Errorf("handshake failed: %w", err)
	}

	// Login
	if err := m.login(); err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	// Start keepalive handler
	m.wg.Add(1)
	go m.keepaliveHandler()

	m.log.Info("Session started successfully")
	m.events.Emit("session_started")

	return nil
}

// handshake performs the handshake.
func (m *Manager) handshake(host string, port uint16) error {
	m.log.Info("Performing handshake...")

	m.setState(protocol.Handshaking)

	pkt := protocol.NewHandshakePacket(m.protocol, host, port, 2) // 2 = login
	if err := m.writePacket(pkt); err != nil {
		return err
	}

	m.log.Debug("Handshake packet sent")
	return nil
}

// login performs the login.
func (m *Manager) login() error {
	m.log.Info("Logging in as %s...", m.username)

	m.setState(protocol.Login)

	pkt := protocol.NewLoginStartPacket(m.username)
	if err := m.writePacket(pkt); err != nil {
		return err
	}

	m.log.Debug("Login start packet sent")

	// Wait for login success
	loginCh := make(chan *loginResult, 1)
	sub := m.events.Subscribe("packet", m.makeLoginHandler(loginCh))
	defer sub.Unsubscribe()

	select {
	case result := <-loginCh:
		if result.err != nil {
			return result.err
		}
		m.log.Info("Logged in successfully! UUID: %s", result.uuid)
		m.setState(protocol.Play)
		m.events.Emit("logged_in", result.uuid, m.username)
		return nil

	case <-time.After(15 * time.Second):
		return fmt.Errorf("login timeout")

	case <-m.ctx.Done():
		return fmt.Errorf("context cancelled")
	}
}

type loginResult struct {
	uuid string
	err  error
}

func (m *Manager) makeLoginHandler(ch chan<- *loginResult) func(...interface{}) {
	return func(data ...interface{}) {
		pkt, ok := data[0].(*protocol.Packet)
		if !ok {
			return
		}

		// Only handle login state packets
		currentState := m.getState()
		if currentState != protocol.Login {
			return
		}

		switch pkt.ID {
		case protocol.LoginSuccessPacketID:
			m.log.Debug("Received login success")
			uuid, _, err := protocol.ParseLoginSuccess(pkt.Data)
			if err != nil {
				ch <- &loginResult{err: err}
				return
			}
			ch <- &loginResult{uuid: uuid}

		case protocol.LoginDisconnectPacketID:
			r := protocol.NewReader(bytes.NewReader(pkt.Data))
			reason, _ := r.ReadString()
			ch <- &loginResult{err: fmt.Errorf("disconnected: %s", reason)}
		}
	}
}

// keepaliveHandler handles keepalive packets.
func (m *Manager) keepaliveHandler() {
	defer m.wg.Done()
	m.log.Debug("Keepalive handler started")

	sub := m.events.Subscribe("packet", m.handleKeepalivePacket)
	defer sub.Unsubscribe()

	<-m.ctx.Done()
	m.log.Debug("Keepalive handler stopped")
}

func (m *Manager) handleKeepalivePacket(data ...interface{}) {
	pkt, ok := data[0].(*protocol.Packet)
	if !ok {
		return
	}

	// Only handle play state keepalives
	if m.getState() != protocol.Play {
		return
	}

	if pkt.ID == protocol.KeepAliveClientboundPacketID {
		id, err := protocol.ParseKeepAlive(pkt.Data)
		if err != nil {
			m.log.Error("Failed to parse keepalive: %v", err)
			return
		}

		m.log.Debug("Received keepalive, responding...")
		response := protocol.NewKeepAlivePacket(id)
		if err := m.writePacket(response); err != nil {
			m.log.Error("Failed to send keepalive response: %v", err)
		}
	}
}

// writePacket writes a packet to the connection.
func (m *Manager) writePacket(pkt *protocol.Packet) error {
	// Serialize packet
	w := protocol.NewWriter()
	w.WriteVarInt(pkt.ID)
	if pkt.Data != nil {
		w.WriteRaw(pkt.Data)
	}

	data := w.Bytes()

	// Write length prefix
	lengthBuf := make([]byte, protocol.VarIntSize(int32(len(data))))
	lengthWriter := protocol.NewWriter()
	lengthWriter.WriteVarInt(int32(len(data)))
	copy(lengthBuf, lengthWriter.Bytes())

	// Write length + data
	fullData := append(lengthBuf, data...)

	_, err := m.conn.Write(fullData)
	return err
}

// setState sets the current state.
func (m *Manager) setState(state protocol.State) {
	m.mu.Lock()
	defer m.mu.Unlock()

	oldState := m.state
	m.state = state
	m.log.Debug("State transition: %s -> %s", oldState, state)

	m.events.Emit("state_changed", oldState, state)
}

// GetState returns the current state.
func (m *Manager) GetState() protocol.State {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state
}

// getState returns the current state (internal).
func (m *Manager) getState() protocol.State {
	return m.GetState()
}

// Close closes the session manager.
func (m *Manager) Close() error {
	m.log.Info("Closing session manager...")
	m.cancel()
	m.wg.Wait()
	m.log.Info("Session manager closed")
	return nil
}

// HandlePacket handles an incoming packet (called by packet reader).
func (m *Manager) HandlePacket(pkt *protocol.Packet) {
	m.events.Emit("packet", pkt)
}
