// Package protocol реализует Minecraft protocol codec.
package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

var (
	ErrVarIntTooBig  = errors.New("VarInt too big")
	ErrVarLongTooBig = errors.New("VarLong too big")
)

// Reader reads Minecraft protocol data.
type Reader struct {
	r io.Reader
}

// NewReader creates a new protocol reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{r: r}
}

// ReadVarInt reads a VarInt.
func (r *Reader) ReadVarInt() (int32, error) {
	var result uint32
	var shift uint

	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

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

// ReadVarLong reads a VarLong.
func (r *Reader) ReadVarLong() (int64, error) {
	var result uint64
	var shift uint

	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

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

// ReadByte reads a single byte.
func (r *Reader) ReadByte() (byte, error) {
	buf := [1]byte{0}
	_, err := io.ReadFull(r.r, buf[:])
	return buf[0], err
}

// ReadBytes reads a byte array prefixed with VarInt length.
func (r *Reader) ReadBytes() ([]byte, error) {
	length, err := r.ReadVarInt()
	if err != nil {
		return nil, err
	}

	if length < 0 {
		return nil, fmt.Errorf("negative length: %d", length)
	}

	if length > 10_000_000 { // 10MB limit
		return nil, fmt.Errorf("length too large: %d", length)
	}

	data := make([]byte, length)
	_, err = io.ReadFull(r.r, data)
	return data, err
}

// ReadString reads a string prefixed with VarInt length.
func (r *Reader) ReadString() (string, error) {
	data, err := r.ReadBytes()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ReadUint8 reads a uint8.
func (r *Reader) ReadUint8() (uint8, error) {
	return r.ReadByte()
}

// ReadUint16 reads a big-endian uint16.
func (r *Reader) ReadUint16() (uint16, error) {
	buf := [2]byte{}
	_, err := io.ReadFull(r.r, buf[:])
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(buf[:]), nil
}

// ReadInt16 reads a big-endian int16.
func (r *Reader) ReadInt16() (int16, error) {
	v, err := r.ReadUint16()
	return int16(v), err
}

// ReadUint32 reads a big-endian uint32.
func (r *Reader) ReadUint32() (uint32, error) {
	buf := [4]byte{}
	_, err := io.ReadFull(r.r, buf[:])
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(buf[:]), nil
}

// ReadInt32 reads a big-endian int32.
func (r *Reader) ReadInt32() (int32, error) {
	v, err := r.ReadUint32()
	return int32(v), err
}

// ReadUint64 reads a big-endian uint64.
func (r *Reader) ReadUint64() (uint64, error) {
	buf := [8]byte{}
	_, err := io.ReadFull(r.r, buf[:])
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(buf[:]), nil
}

// ReadInt64 reads a big-endian int64.
func (r *Reader) ReadInt64() (int64, error) {
	v, err := r.ReadUint64()
	return int64(v), err
}

// ReadFloat32 reads a big-endian float32.
func (r *Reader) ReadFloat32() (float32, error) {
	bits, err := r.ReadUint32()
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(bits), nil
}

// ReadFloat64 reads a big-endian float64.
func (r *Reader) ReadFloat64() (float64, error) {
	bits, err := r.ReadUint64()
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(bits), nil
}

// ReadBool reads a boolean.
func (r *Reader) ReadBool() (bool, error) {
	b, err := r.ReadByte()
	return b != 0, err
}

// Writer writes Minecraft protocol data.
type Writer struct {
	buf []byte
}

// NewWriter creates a new protocol writer.
func NewWriter() *Writer {
	return &Writer{
		buf: make([]byte, 0, 1024),
	}
}

// Bytes returns the written bytes.
func (w *Writer) Bytes() []byte {
	return w.buf
}

// Reset resets the buffer.
func (w *Writer) Reset() {
	w.buf = w.buf[:0]
}

// WriteVarInt writes a VarInt.
func (w *Writer) WriteVarInt(v int32) {
	uv := uint32(v)

	for {
		b := byte(uv & 0x7F)
		uv >>= 7
		if uv != 0 {
			b |= 0x80
		}
		w.buf = append(w.buf, b)

		if uv == 0 {
			break
		}
	}
}

// WriteVarLong writes a VarLong.
func (w *Writer) WriteVarLong(v int64) {
	uv := uint64(v)

	for {
		b := byte(uv & 0x7F)
		uv >>= 7
		if uv != 0 {
			b |= 0x80
		}
		w.buf = append(w.buf, b)

		if uv == 0 {
			break
		}
	}
}

// WriteByte writes a byte.
func (w *Writer) WriteByte(b byte) {
	w.buf = append(w.buf, b)
}

// WriteBytes writes a byte array with VarInt length prefix.
func (w *Writer) WriteBytes(data []byte) {
	w.WriteVarInt(int32(len(data)))
	w.buf = append(w.buf, data...)
}

// WriteString writes a string with VarInt length prefix.
func (w *Writer) WriteString(s string) {
	w.WriteBytes([]byte(s))
}

// WriteUint8 writes a uint8.
func (w *Writer) WriteUint8(v uint8) {
	w.buf = append(w.buf, v)
}

// WriteUint16 writes a big-endian uint16.
func (w *Writer) WriteUint16(v uint16) {
	w.buf = append(w.buf, byte(v>>8), byte(v))
}

// WriteInt16 writes a big-endian int16.
func (w *Writer) WriteInt16(v int16) {
	w.WriteUint16(uint16(v))
}

// WriteUint32 writes a big-endian uint32.
func (w *Writer) WriteUint32(v uint32) {
	w.buf = append(w.buf,
		byte(v>>24),
		byte(v>>16),
		byte(v>>8),
		byte(v),
	)
}

// WriteInt32 writes a big-endian int32.
func (w *Writer) WriteInt32(v int32) {
	w.WriteUint32(uint32(v))
}

// WriteUint64 writes a big-endian uint64.
func (w *Writer) WriteUint64(v uint64) {
	w.buf = append(w.buf,
		byte(v>>56),
		byte(v>>48),
		byte(v>>40),
		byte(v>>32),
		byte(v>>24),
		byte(v>>16),
		byte(v>>8),
		byte(v),
	)
}

// WriteInt64 writes a big-endian int64.
func (w *Writer) WriteInt64(v int64) {
	w.WriteUint64(uint64(v))
}

// WriteFloat32 writes a big-endian float32.
func (w *Writer) WriteFloat32(v float32) {
	w.WriteUint32(math.Float32bits(v))
}

// WriteFloat64 writes a big-endian float64.
func (w *Writer) WriteFloat64(v float64) {
	w.WriteUint64(math.Float64bits(v))
}

// WriteBool writes a boolean.
func (w *Writer) WriteBool(v bool) {
	if v {
		w.buf = append(w.buf, 1)
	} else {
		w.buf = append(w.buf, 0)
	}
}

// WriteRaw writes raw bytes to the buffer.
func (w *Writer) WriteRaw(data []byte) {
	w.buf = append(w.buf, data...)
}

// ReadDouble reads a big-endian float64.
func (r *Reader) ReadDouble() (float64, error) {
	bits, err := r.ReadUint64()
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(bits), nil
}

// WriteDouble writes a big-endian float64.
func (w *Writer) WriteDouble(v float64) {
	w.WriteFloat64(v)
}

// ReadBlockPos reads a block position.
func (r *Reader) ReadBlockPos() (*BlockPos, error) {
	x, err := r.ReadInt64()
	if err != nil {
		return nil, err
	}

	return &BlockPos{
		x: int32(x >> 38),
		y: int32((x << 52) >> 52),
		z: int32((x << 26) >> 38),
	}, nil
}

// BlockPos represents a position in the world.
type BlockPos struct {
	x, y, z int32
}

// X returns the X coordinate.
func (b *BlockPos) X() int32 {
	return b.x
}

// Y returns the Y coordinate.
func (b *BlockPos) Y() int32 {
	return b.y
}

// Z returns the Z coordinate.
func (b *BlockPos) Z() int32 {
	return b.z
}

// VarIntSize returns the size of a VarInt in bytes.
func VarIntSize(v int32) int {
	uv := uint32(v)
	if uv == 0 {
		return 1
	}

	size := 0
	for {
		size++
		uv >>= 7
		if uv == 0 {
			break
		}
	}
	return size
}
