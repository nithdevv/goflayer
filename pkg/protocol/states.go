// Package protocol handles the Minecraft protocol implementation.
//
// This package contains the protocol client, packet serialization/deserialization,
// compression, encryption, and protocol state management.
package protocol

// State represents the current state of the Minecraft protocol connection.
//
// The Minecraft protocol has different states that determine which packets
// can be sent and received. The connection flows through these states:
//
// Handshaking → Status or Login → Play
//
// Example flow for server ping:
//   Handshaking → Status (get server info) → disconnect
//
// Example flow for joining:
//   Handshaking → Login (authenticate) → Play (gameplay)
type State int

const (
	// Handshaking is the initial state when a client connects.
	// The client sends a handshake packet with the protocol version,
	// server address, port, and next state (Status or Login).
	Handshaking State = iota

	// Status state is used for querying server information.
	// Clients in this state can request server status (motd, player count, etc.)
	// without actually joining the game.
	// This is what server lists use to display server information.
	Status

	// Login state is used for authentication.
	// Clients in this state exchange login packets, authenticate with
	// Mojang/Microsoft servers (if online mode), and transition to Play state.
	Login

	// Play state is the main gameplay state.
	// Most game packets are exchanged in this state: movement, blocks,
	// entities, chat, inventory, etc.
	Play
)

// String returns the string representation of the protocol state.
func (s State) String() string {
	switch s {
	case Handshaking:
		return "Handshaking"
	case Status:
		return "Status"
	case Login:
		return "Login"
	case Play:
		return "Play"
	default:
		return "Unknown"
	}
}

// CanTransitionTo checks if a transition to another state is valid.
//
// Valid transitions:
//   - Handshaking → Status
//   - Handshaking → Login
//   - Status → disconnect (end of connection)
//   - Login → Play
//   - Play → disconnect (end of connection)
func (s State) CanTransitionTo(newState State) bool {
	switch s {
	case Handshaking:
		return newState == Status || newState == Login
	case Status:
		return newState == Handshaking // reconnect after status check
	case Login:
		return newState == Play
	case Play:
		return false // Play is terminal state, only disconnect
	default:
		return false
	}
}

// IsHandshaking returns true if the state is Handshaking.
func (s State) IsHandshaking() bool {
	return s == Handshaking
}

// IsStatus returns true if the state is Status.
func (s State) IsStatus() bool {
	return s == Status
}

// IsLogin returns true if the state is Login.
func (s State) IsLogin() bool {
	return s == Login
}

// IsPlay returns true if the state is Play.
func (s State) IsPlay() bool {
	return s == Play
}

// ConnectionState tracks the current connection state.
type ConnectionState struct {
	state      State
	previous   State
	compressionEnabled bool
	encryptionEnabled  bool
	version    string
	protocolVersion int
}

// NewConnectionState creates a new connection state in Handshaking.
func NewConnectionState() *ConnectionState {
	return &ConnectionState{
		state: Handshaking,
	}
}

// State returns the current protocol state.
func (cs *ConnectionState) State() State {
	return cs.state
}

// SetState changes the protocol state.
// Returns error if the transition is invalid.
func (cs *ConnectionState) SetState(newState State) error {
	if !cs.state.CanTransitionTo(newState) {
		return ErrInvalidState
	}
	cs.previous = cs.state
	cs.state = newState
	return nil
}

// PreviousState returns the previous protocol state before the last transition.
func (cs *ConnectionState) PreviousState() State {
	return cs.previous
}

// EnableCompression enables packet compression (Zlib).
func (cs *ConnectionState) EnableCompression() {
	cs.compressionEnabled = true
}

// DisableCompression disables packet compression.
func (cs *ConnectionState) DisableCompression() {
	cs.compressionEnabled = false
}

// CompressionEnabled returns true if compression is enabled.
func (cs *ConnectionState) CompressionEnabled() bool {
	return cs.compressionEnabled
}

// EnableEncryption enables packet encryption (AES).
func (cs *ConnectionState) EnableEncryption() {
	cs.encryptionEnabled = true
}

// DisableEncryption disables packet encryption.
func (cs *ConnectionState) DisableEncryption() {
	cs.encryptionEnabled = false
}

// EncryptionEnabled returns true if encryption is enabled.
func (cs *ConnectionState) EncryptionEnabled() bool {
	return cs.encryptionEnabled
}

// SetVersion sets the Minecraft version and protocol version.
func (cs *ConnectionState) SetVersion(version string, protocolVersion int) {
	cs.version = version
	cs.protocolVersion = protocolVersion
}

// Version returns the Minecraft version string.
func (cs *ConnectionState) Version() string {
	return cs.version
}

// ProtocolVersion returns the protocol version number.
func (cs *ConnectionState) ProtocolVersion() int {
	return cs.protocolVersion
}

// Reset resets the connection state to initial values.
func (cs *ConnectionState) Reset() {
	cs.state = Handshaking
	cs.previous = Handshaking
	cs.compressionEnabled = false
	cs.encryptionEnabled = false
	cs.version = ""
	cs.protocolVersion = 0
}

// ConnectionState represents the current connection status.
// This is separate from protocol state and tracks the actual connection.
type ConnectionStatus int

const (
	// Disconnected means the connection is closed or not established.
	Disconnected ConnectionStatus = iota

	// Connecting means the connection is being established (TCP handshake).
	Connecting

	// Connected means the TCP connection is established but protocol
	// handshake is not complete.
	Connected

	// Authenticating means the client is authenticating with the server.
	Authenticating

	//LoggedIn means the client has successfully logged in and is in Play state.
	LoggedIn

	// Disconnecting means the connection is being closed.
	Disconnecting
)

// String returns the string representation of the connection status.
func (s ConnectionStatus) String() string {
	switch s {
	case Disconnected:
		return "Disconnected"
	case Connecting:
		return "Connecting"
	case Connected:
		return "Connected"
	case Authenticating:
		return "Authenticating"
	case LoggedIn:
		return "LoggedIn"
	case Disconnecting:
		return "Disconnecting"
	default:
		return "Unknown"
	}
}
