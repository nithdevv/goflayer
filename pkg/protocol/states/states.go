// Package states defines Minecraft protocol states.
//
// The Minecraft protocol has multiple states that a connection progresses through:
// Handshaking -> Status/Login -> Play
package states

// State represents a protocol state.
type State int

const (
	// Handshaking is the initial state after connecting.
	Handshaking State = iota

	// Status is for server list ping.
	Status

	// Login is the authentication phase.
	Login

	// Play is the main game state.
	Play
)

// String returns the string representation of the state.
func (s State) String() string {
	switch s {
	case Handshaking:
		return "handshaking"
	case Status:
		return "status"
	case Login:
		return "login"
	case Play:
		return "play"
	default:
		return "unknown"
	}
}

// CanTransitionTo returns true if the state can transition to the target state.
func (s State) CanTransitionTo(target State) bool {
	switch s {
	case Handshaking:
		return target == Status || target == Login
	case Status:
		return target == Handshaking
	case Login:
		return target == Play
	case Play:
		return false
	default:
		return false
	}
}

// PacketDirection represents the direction of a packet (clientbound or serverbound).
type PacketDirection int

const (
	// Clientbound packets are sent from server to client.
	Clientbound PacketDirection = iota

	// Serverbound packets are sent from client to server.
	Serverbound
)

// String returns the string representation of the direction.
func (d PacketDirection) String() string {
	switch d {
	case Clientbound:
		return "clientbound"
	case Serverbound:
		return "serverbound"
	default:
		return "unknown"
	}
}
