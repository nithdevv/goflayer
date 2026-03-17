// Package codec implements Minecraft protocol encoding and decoding.
//
// Minecraft uses a custom binary format with variable-length integers (VarInt).
package codec

import (
	"errors"
	"fmt"
	"io"
)

var (
	ErrVarIntTooBig  = errors.New("VarInt too big")
	ErrVarLongTooBig = errors.New("VarLong too big")
)

// ReadVarInt reads a variable-length integer from the reader.
// VarInt is a variable-length integer encoding used by Minecraft.
// It uses 7 bits per byte, with the most significant bit indicating
// whether more bytes follow.
func ReadVarInt(r io.Reader) (int32, error) {
	var result uint32
	var shift uint

	for {
		buf := []byte{0}
		_, err := io.ReadFull(r, buf)
		if err != nil {
			return 0, err
		}

		b := buf[0]
		result |= uint32(b&0x7F) << shift

		if (b & 0x80) == 0 {
			break
		}

		shift += 7
		if shift >= 35 {
			return 0, ErrVarIntTooBig
		}
	}

	return int32(result), nil
}

// ReadVarLong reads a variable-length long (64-bit) from the reader.
func ReadVarLong(r io.Reader) (int64, error) {
	var result uint64
	var shift uint

	for {
		buf := []byte{0}
		_, err := io.ReadFull(r, buf)
		if err != nil {
			return 0, err
		}

		b := buf[0]
		result |= uint64(b&0x7F) << shift

		if (b & 0x80) == 0 {
			break
		}

		shift += 7
		if shift >= 70 {
			return 0, ErrVarLongTooBig
		}
	}

	return int64(result), nil
}

// WriteVarInt writes a variable-length integer to the byte slice.
// Returns the number of bytes written.
func WriteVarInt(buf []byte, value int32) int {
	uValue := uint32(value)

	for {
		b := byte(uValue & 0x7F)
		uValue >>= 7
		if uValue != 0 {
			b |= 0x80
		}
		buf = append(buf, b)

		if uValue == 0 {
			break
		}
	}

	return len(buf)
}

// WriteVarLong writes a variable-length long to the byte slice.
// Returns the number of bytes written.
func WriteVarLong(buf []byte, value int64) int {
	uValue := uint64(value)

	for {
		b := byte(uValue & 0x7F)
		uValue >>= 7
		if uValue != 0 {
			b |= 0x80
		}
		buf = append(buf, b)

		if uValue == 0 {
			break
		}
	}

	return len(buf)
}

// VarIntSize returns the number of bytes needed to encode the value as VarInt.
func VarIntSize(value int32) int {
	uValue := uint32(value)
	if uValue == 0 {
		return 1
	}

	size := 0
	for {
		size++
		uValue >>= 7
		if uValue == 0 {
			break
		}
	}
	return size
}

// VarLongSize returns the number of bytes needed to encode the value as VarLong.
func VarLongSize(value int64) int {
	uValue := uint64(value)
	if uValue == 0 {
		return 1
	}

	size := 0
	for {
		size++
		uValue >>= 7
		if uValue == 0 {
			break
		}
	}
	return size
}

// Reader wraps an io.Reader with protocol decoding methods.
type Reader struct {
	r io.Reader
}

// NewReader creates a new protocol reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{r: r}
}

// ReadVarInt reads a VarInt from the underlying reader.
func (r *Reader) ReadVarInt() (int32, error) {
	return ReadVarInt(r.r)
}

// ReadVarLong reads a VarLong from the underlying reader.
func (r *Reader) ReadVarLong() (int64, error) {
	return ReadVarLong(r.r)
}

// ReadBool reads a boolean (1 byte) from the underlying reader.
func (r *Reader) ReadBool() (bool, error) {
	buf := []byte{0}
	_, err := io.ReadFull(r.r, buf)
	if err != nil {
		return false, err
	}
	return buf[0] != 0, nil
}

// ReadByte reads a single byte from the underlying reader.
func (r *Reader) ReadByte() (byte, error) {
	buf := []byte{0}
	_, err := io.ReadFull(r.r, buf)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

// ReadBytes reads a byte array prefixed with VarInt length.
func (r *Reader) ReadBytes() ([]byte, error) {
	length, err := r.ReadVarInt()
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, fmt.Errorf("negative byte array length: %d", length)
	}
	if length > 10000000 { // Sanity check: 10MB max
		return nil, fmt.Errorf("byte array too large: %d", length)
	}

	data := make([]byte, length)
	_, err = io.ReadFull(r.r, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ReadString reads a string prefixed with VarInt length.
func (r *Reader) ReadString() (string, error) {
	data, err := r.ReadBytes()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ReadUint8 reads an unsigned 8-bit integer.
func (r *Reader) ReadUint8() (uint8, error) {
	buf := []byte{0}
	_, err := io.ReadFull(r.r, buf)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

// ReadUint16 reads a big-endian unsigned 16-bit integer.
func (r *Reader) ReadUint16() (uint16, error) {
	buf := make([]byte, 2)
	_, err := io.ReadFull(r.r, buf)
	if err != nil {
		return 0, err
	}
	return uint16(buf[0])<<8 | uint16(buf[1]), nil
}

// ReadInt16 reads a big-endian signed 16-bit integer.
func (r *Reader) ReadInt16() (int16, error) {
	v, err := r.ReadUint16()
	return int16(v), err
}

// ReadUint32 reads a big-endian unsigned 32-bit integer.
func (r *Reader) ReadUint32() (uint32, error) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(r.r, buf)
	if err != nil {
		return 0, err
	}
	return uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3]), nil
}

// ReadInt32 reads a big-endian signed 32-bit integer.
func (r *Reader) ReadInt32() (int32, error) {
	v, err := r.ReadUint32()
	return int32(v), err
}

// ReadUint64 reads a big-endian unsigned 64-bit integer.
func (r *Reader) ReadUint64() (uint64, error) {
	buf := make([]byte, 8)
	_, err := io.ReadFull(r.r, buf)
	if err != nil {
		return 0, err
	}
	return uint64(buf[0])<<56 | uint64(buf[1])<<48 | uint64(buf[2])<<40 | uint64(buf[3])<<32 |
		uint64(buf[4])<<24 | uint64(buf[5])<<16 | uint64(buf[6])<<8 | uint64(buf[7]), nil
}

// ReadInt64 reads a big-endian signed 64-bit integer.
func (r *Reader) ReadInt64() (int64, error) {
	v, err := r.ReadUint64()
	return int64(v), err
}

// ReadFloat32 reads a big-endian 32-bit float.
func (r *Reader) ReadFloat32() (float32, error) {
	bits, err := r.ReadUint32()
	if err != nil {
		return 0, err
	}
	// Convert bits to float32
	return float32(bits), nil // Placeholder - need proper conversion
}

// ReadFloat64 reads a big-endian 64-bit float.
func (r *Reader) ReadFloat64() (float64, error) {
	bits, err := r.ReadUint64()
	if err != nil {
		return 0, err
	}
	// Convert bits to float64
	return float64(bits), nil // Placeholder - need proper conversion
}

// Writer wraps an io.Writer with protocol encoding methods.
type Writer struct {
	w   io.Writer
	buf []byte
}

// NewWriter creates a new protocol writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w:   w,
		buf: make([]byte, 0, 1024),
	}
}

// WriteVarInt writes a VarInt to the buffer.
func (w *Writer) WriteVarInt(value int32) error {
	w.buf = append(w.buf, encodeVarInt(value)...)
	return nil
}

// WriteVarLong writes a VarLong to the buffer.
func (w *Writer) WriteVarLong(value int64) error {
	w.buf = append(w.buf, encodeVarLong(value)...)
	return nil
}

// encodeVarInt encodes a VarInt to a byte slice.
func encodeVarInt(value int32) []byte {
	uValue := uint32(value)
	buf := make([]byte, 0)

	for {
		b := byte(uValue & 0x7F)
		uValue >>= 7
		if uValue != 0 {
			b |= 0x80
		}
		buf = append(buf, b)

		if uValue == 0 {
			break
		}
	}

	return buf
}

// encodeVarLong encodes a VarLong to a byte slice.
func encodeVarLong(value int64) []byte {
	uValue := uint64(value)
	buf := make([]byte, 0)

	for {
		b := byte(uValue & 0x7F)
		uValue >>= 7
		if uValue != 0 {
			b |= 0x80
		}
		buf = append(buf, b)

		if uValue == 0 {
			break
		}
	}

	return buf
}

// WriteBool writes a boolean (1 byte) to the buffer.
func (w *Writer) WriteBool(value bool) error {
	if value {
		w.buf = append(w.buf, 1)
	} else {
		w.buf = append(w.buf, 0)
	}
	return nil
}

// WriteByte writes a single byte to the buffer.
func (w *Writer) WriteByte(value byte) error {
	w.buf = append(w.buf, value)
	return nil
}

// WriteBytes writes a byte array prefixed with VarInt length to the buffer.
func (w *Writer) WriteBytes(data []byte) error {
	w.WriteVarInt(int32(len(data)))
	w.buf = append(w.buf, data...)
	return nil
}

// WriteString writes a string prefixed with VarInt length to the buffer.
func (w *Writer) WriteString(value string) error {
	return w.WriteBytes([]byte(value))
}

// WriteUint8 writes an unsigned 8-bit integer to the buffer.
func (w *Writer) WriteUint8(value uint8) error {
	w.buf = append(w.buf, value)
	return nil
}

// WriteUint16 writes a big-endian unsigned 16-bit integer to the buffer.
func (w *Writer) WriteUint16(value uint16) error {
	w.buf = append(w.buf, byte(value>>8), byte(value))
	return nil
}

// WriteInt16 writes a big-endian signed 16-bit integer to the buffer.
func (w *Writer) WriteInt16(value int16) error {
	return w.WriteUint16(uint16(value))
}

// WriteUint32 writes a big-endian unsigned 32-bit integer to the buffer.
func (w *Writer) WriteUint32(value uint32) error {
	w.buf = append(w.buf,
		byte(value>>24),
		byte(value>>16),
		byte(value>>8),
		byte(value))
	return nil
}

// WriteInt32 writes a big-endian signed 32-bit integer to the buffer.
func (w *Writer) WriteInt32(value int32) error {
	return w.WriteUint32(uint32(value))
}

// WriteUint64 writes a big-endian unsigned 64-bit integer to the buffer.
func (w *Writer) WriteUint64(value uint64) error {
	w.buf = append(w.buf,
		byte(value>>56),
		byte(value>>48),
		byte(value>>40),
		byte(value>>32),
		byte(value>>24),
		byte(value>>16),
		byte(value>>8),
		byte(value))
	return nil
}

// WriteInt64 writes a big-endian signed 64-bit integer to the buffer.
func (w *Writer) WriteInt64(value int64) error {
	return w.WriteUint64(uint64(value))
}

// WriteFloat32 writes a big-endian 32-bit float to the buffer.
func (w *Writer) WriteFloat32(value float32) error {
	return w.WriteUint32(uint32(value)) // Placeholder - need proper conversion
}

// WriteFloat64 writes a big-endian 64-bit float to the buffer.
func (w *Writer) WriteFloat64(value float64) error {
	return w.WriteUint64(uint64(value)) // Placeholder - need proper conversion
}

// Bytes returns the accumulated buffer.
func (w *Writer) Bytes() []byte {
	return w.buf
}

// Reset clears the buffer.
func (w *Writer) Reset() {
	w.buf = w.buf[:0]
}

// Flush writes the buffer to the underlying writer and resets.
func (w *Writer) Flush() error {
	if len(w.buf) == 0 {
		return nil
	}
	_, err := w.w.Write(w.buf)
	w.Reset()
	return err
}
