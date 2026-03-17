package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

// Serializer handles serialization of packets to binary format.
type Serializer struct {
	version string
	state   State
}

// NewSerializer creates a new packet serializer.
func NewSerializer(version string, state State) *Serializer {
	return &Serializer{
		version: version,
		state:   state,
	}
}

// Serialize serializes a packet to binary format.
// Returns the packet buffer ready to be sent (without length prefix).
func (s *Serializer) Serialize(packet *Packet) ([]byte, error) {
	buffer := NewPacketBuffer()

	// Write packet ID (will be added based on registry)
	// For now, write a placeholder
	if err := buffer.WriteVarInt(0); err != nil {
		return nil, err
	}

	// Write packet data
	if err := s.writePacketData(buffer, packet); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// writePacketData writes the packet-specific data to the buffer.
func (s *Serializer) writePacketData(buffer *PacketBuffer, packet *Packet) error {
	// This is a simplified implementation
	// In the full version, this would use protodef or similar
	// to serialize based on packet schema

	for key, value := range packet.Data {
		if err := s.writeValue(buffer, value); err != nil {
			return fmt.Errorf("failed to write field %s: %w", key, err)
		}
	}

	return nil
}

// writeValue writes a value to the buffer based on its type.
func (s *Serializer) writeValue(buffer *PacketBuffer, value interface{}) error {
	switch v := value.(type) {
	case bool:
		return buffer.WriteBool(v)
	case byte:
		return buffer.WriteByte(v)
	case int8:
		return buffer.WriteInt8(v)
	case int16:
		return buffer.WriteInt16(v)
	case int32:
		return buffer.WriteInt32(v)
	case int64:
		return buffer.WriteInt64(v)
	case int:
		return buffer.WriteVarInt(v)
	case uint8:
		return buffer.WriteUInt8(v)
	case uint16:
		return buffer.WriteUInt16(v)
	case float32:
		return buffer.WriteFloat32(v)
	case float64:
		return buffer.WriteFloat64(v)
	case string:
		return buffer.WriteString(v)
	case []byte:
		return buffer.WriteBytes(v)
	case []interface{}:
		return buffer.WriteArray(v, s.writeValue)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
}

// Deserializer handles deserialization of packets from binary format.
type Deserializer struct {
	version   string
	state     State
	registry  *PacketRegistry
}

// NewDeserializer creates a new packet deserializer.
func NewDeserializer(version string, state State, registry *PacketRegistry) *Deserializer {
	return &Deserializer{
		version:  version,
		state:    state,
		registry: registry,
	}
}

// Deserialize deserializes a packet from binary format.
func (d *Deserializer) Deserialize(buffer []byte) (*Packet, error) {
	r := bytes.NewReader(buffer)
	packet := &Packet{
		State: d.state,
	}

	// Read packet ID
	packetID, err := ReadVarInt(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read packet ID: %w", err)
	}
	packet.ID = packetID

	// Get packet name from registry
	packetName, ok := d.registry.GetPacketName(ClientBound, packetID)
	if !ok {
		return nil, fmt.Errorf("unknown packet ID: %d in state %s", packetID, d.state)
	}
	packet.Name = packetName

	// Read packet data
	packet.Data = make(map[string]interface{})
	if err := d.readPacketData(r, packet); err != nil {
		return nil, fmt.Errorf("failed to read packet data: %w", err)
	}

	packet.Buffer = buffer
	return packet, nil
}

// readPacketData reads packet-specific data from the reader.
func (d *Deserializer) readPacketData(r io.Reader, packet *Packet) error {
	// This is a simplified implementation
	// In the full version, this would use protodef or similar
	// to deserialize based on packet schema

	// For now, just read remaining data as bytes
	remaining, _ := io.ReadAll(r)
	if len(remaining) > 0 {
		packet.Data["raw"] = remaining
	}

	return nil
}

// PacketBuffer extends bytes.Buffer with Minecraft protocol-specific read/write methods.
type PacketBuffer struct {
	buffer *bytes.Buffer
	reader io.Reader
}

// WriteBool writes a boolean value.
func (pb *PacketBuffer) WriteBool(value bool) error {
	var b byte
	if value {
		b = 1
	}
	return pb.buffer.WriteByte(b)
}

// WriteInt8 writes an 8-bit signed integer.
func (pb *PacketBuffer) WriteInt8(value int8) error {
	return pb.buffer.WriteByte(byte(value))
}

// WriteUInt8 writes an 8-bit unsigned integer.
func (pb *PacketBuffer) WriteUInt8(value uint8) error {
	return pb.buffer.WriteByte(value)
}

// WriteInt16 writes a 16-bit signed integer (big-endian).
func (pb *PacketBuffer) WriteInt16(value int16) error {
	return binary.Write(pb.buffer, binary.BigEndian, value)
}

// WriteUInt16 writes a 16-bit unsigned integer (big-endian).
func (pb *PacketBuffer) WriteUInt16(value uint16) error {
	return binary.Write(pb.buffer, binary.BigEndian, value)
}

// WriteInt32 writes a 32-bit signed integer (big-endian).
func (pb *PacketBuffer) WriteInt32(value int32) error {
	return binary.Write(pb.buffer, binary.BigEndian, value)
}

// WriteInt64 writes a 64-bit signed integer (big-endian).
func (pb *PacketBuffer) WriteInt64(value int64) error {
	return binary.Write(pb.buffer, binary.BigEndian, value)
}

// WriteFloat32 writes a 32-bit float (big-endian).
func (pb *PacketBuffer) WriteFloat32(value float32) error {
	return binary.Write(pb.buffer, binary.BigEndian, value)
}

// WriteFloat64 writes a 64-bit float (big-endian).
func (pb *PacketBuffer) WriteFloat64(value float64) error {
	return binary.Write(pb.buffer, binary.BigEndian, value)
}

// WriteVarInt writes a variable-length integer.
func (pb *PacketBuffer) WriteVarInt(value int) error {
	for {
		temp := byte(value & 0x7F)
		value >>= 7
		if value != 0 {
			temp |= 0x80
		}
		if err := pb.buffer.WriteByte(temp); err != nil {
			return err
		}
		if value == 0 {
			break
		}
	}
	return nil
}

// WriteVarLong writes a variable-length long.
func (pb *PacketBuffer) WriteVarLong(value int64) error {
	for {
		temp := byte(value & 0x7F)
		value >>= 7
		if value != 0 {
			temp |= 0x80
		}
		if err := pb.buffer.WriteByte(temp); err != nil {
			return err
		}
		if value == 0 {
			break
		}
	}
	return nil
}

// WriteString writes a VarInt-prefixed string.
func (pb *PacketBuffer) WriteString(value string) error {
	if err := pb.WriteVarInt(len(value)); err != nil {
		return err
	}
	_, err := pb.buffer.WriteString(value)
	return err
}

// WriteBytes writes a length-prefixed byte array.
func (pb *PacketBuffer) WriteBytes(value []byte) error {
	if err := pb.WriteVarInt(len(value)); err != nil {
		return err
	}
	_, err := pb.buffer.Write(value)
	return err
}

// WriteArray writes an array with a length prefix.
func (pb *PacketBuffer) WriteArray(array []interface{}, writeFunc func(*PacketBuffer, interface{}) error) error {
	if err := pb.WriteVarInt(len(array)); err != nil {
		return err
	}
	for _, item := range array {
		if err := writeFunc(pb, item); err != nil {
			return err
		}
	}
	return nil
}

// Bytes returns the buffer contents.
func (pb *PacketBuffer) Bytes() []byte {
	return pb.buffer.Bytes()
}

// Len returns the buffer length.
func (pb *PacketBuffer) Len() int {
	return pb.buffer.Len()
}

// ReadVarInt reads a variable-length integer.
func ReadVarInt(r io.Reader) (int, error) {
	var result int
	var shift uint

	for {
		buf := make([]byte, 1)
		if _, err := io.ReadFull(r, buf); err != nil {
			return 0, err
		}

		b := buf[0]
		result |= int(b&0x7F) << shift

		if (b & 0x80) == 0 {
			return result, nil
		}

		shift += 7
		if shift >= 35 {
			return 0, errors.New("VarInt too big")
		}
	}
}

// ReadVarLong reads a variable-length long.
func ReadVarLong(r io.Reader) (int64, error) {
	var result int64
	var shift uint

	for {
		buf := make([]byte, 1)
		if _, err := io.ReadFull(r, buf); err != nil {
			return 0, err
		}

		b := buf[0]
		result |= int64(b&0x7F) << shift

		if (b & 0x80) == 0 {
			return result, nil
		}

		shift += 7
		if shift >= 70 {
			return 0, errors.New("VarLong too big")
		}
	}
}

// ReadString reads a VarInt-prefixed string.
func ReadString(r io.Reader) (string, error) {
	length, err := ReadVarInt(r)
	if err != nil {
		return "", err
	}

	if length <= 0 {
		return "", nil
	}

	buf := make([]byte, length)
	if _, err := io.ReadFull(r, buf); err != nil {
		return "", err
	}

	return string(buf), nil
}

// ReadBool reads a boolean value.
func ReadBool(r io.Reader) (bool, error) {
	buf := make([]byte, 1)
	if _, err := io.ReadFull(r, buf); err != nil {
		return false, err
	}
	return buf[0] != 0, nil
}

// ReadInt8 reads an 8-bit signed integer.
func ReadInt8(r io.Reader) (int8, error) {
	buf := make([]byte, 1)
	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	return int8(buf[0]), nil
}

// ReadUInt8 reads an 8-bit unsigned integer.
func ReadUInt8(r io.Reader) (uint8, error) {
	buf := make([]byte, 1)
	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	return buf[0], nil
}

// ReadInt16 reads a 16-bit signed integer.
func ReadInt16(r io.Reader) (int16, error) {
	var value int16
	if err := binary.Read(r, binary.BigEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

// ReadUInt16 reads a 16-bit unsigned integer.
func ReadUInt16(r io.Reader) (uint16, error) {
	var value uint16
	if err := binary.Read(r, binary.BigEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

// ReadInt32 reads a 32-bit signed integer.
func ReadInt32(r io.Reader) (int32, error) {
	var value int32
	if err := binary.Read(r, binary.BigEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

// ReadInt64 reads a 64-bit signed integer.
func ReadInt64(r io.Reader) (int64, error) {
	var value int64
	if err := binary.Read(r, binary.BigEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

// ReadFloat32 reads a 32-bit float.
func ReadFloat32(r io.Reader) (float32, error) {
	var value float32
	if err := binary.Read(r, binary.BigEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

// ReadFloat64 reads a 64-bit float.
func ReadFloat64(r io.Reader) (float64, error) {
	var value float64
	if err := binary.Read(r, binary.BigEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

// ReadPosition reads a position (x, y, z) encoded as a single 64-bit long.
func ReadPosition(r io.Reader) (x, y, z int64, err error) {
	val, err := ReadInt64(r)
	if err != nil {
		return 0, 0, 0, err
	}

	x = val >> 38
	y = (val >> 26) & 0xFFF
	z = val << 38 >> 38

	// Handle sign extension for Y
	if y >= 0x800 {
		y -= 0x1000
	}

	return x, y, z, nil
}

// WritePosition writes a position (x, y, z) as a single 64-bit long.
func WritePosition(x, y, z int64) int64 {
	// Encode position: x (26 bits), y (12 bits), z (26 bits)
	return ((x & 0x3FFFFFF) << 38) | ((y & 0xFFF) << 26) | (z & 0x3FFFFFF)
}

// ReadUUID reads a 128-bit UUID.
func ReadUUID(r io.Reader) (string, error) {
	buf := make([]byte, 16)
	if _, err := io.ReadFull(r, buf); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:16]), nil
}

// ReadAngle reads a rotation angle (byte, where 256 = 360 degrees).
func ReadAngle(r io.Reader) (float64, error) {
	b, err := ReadUInt8(r)
	if err != nil {
		return 0, err
	}
	return float64(b) * 360 / 256, nil
}

// WriteAngle writes a rotation angle.
func WriteAngle(angle float64) byte {
	return byte(angle * 256 / 360)
}

// FixedPointInt32 reads an int32 and converts it from fixed-point to float64.
func FixedPointInt32(r io.Reader) (float64, error) {
	val, err := ReadInt32(r)
	if err != nil {
		return 0, err
	}
	return float64(val) / 32, nil
}

// FixedPointInt64 reads an int64 and converts it from fixed-point to float64.
func FixedPointInt64(r io.Reader) (float64, error) {
	val, err := ReadInt64(r)
	if err != nil {
		return 0, err
	}
	return float64(val) / (32 * 32), nil
}

// EncodeFixedPointInt32 converts float64 to fixed-point int32.
func EncodeFixedPointInt32(value float64) int32 {
	return int32(math.Round(value * 32))
}

// EncodeFixedPointInt64 converts float64 to fixed-point int64.
func EncodeFixedPointInt64(value float64) int64 {
	return int64(math.Round(value * 32 * 32))
}
