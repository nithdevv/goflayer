// Package nbt implements the Named Binary Tag (NBT) format.
//
// NBT is a binary format used by Minecraft for data storage.
// It's tree-like structure similar to JSON but with typed tags.
package nbt

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

// TagType represents the type of an NBT tag.
type TagType byte

const (
	TagEnd       TagType = 0
	TagByte      TagType = 1
	TagShort     TagType = 2
	TagInt       TagType = 3
	TagLong      TagType = 4
	TagFloat     TagType = 5
	TagDouble    TagType = 6
	TagByteArray TagType = 7
	TagString    TagType = 8
	TagList      TagType = 9
	TagCompound TagType = 10
	TagIntArray  TagType = 11
	TagLongArray TagType = 12
)

// String returns the string representation of a tag type.
func (t TagType) String() string {
	switch t {
	case TagEnd:
		return "TAG_End"
	case TagByte:
		return "TAG_Byte"
	case TagShort:
		return "TAG_Short"
	case TagInt:
		return "TAG_Int"
	case TagLong:
		return "TAG_Long"
	case TagFloat:
		return "TAG_Float"
	case TagDouble:
		return "TAG_Double"
	case TagByteArray:
		return "TAG_Byte_Array"
	case TagString:
		return "TAG_String"
	case TagList:
		return "TAG_List"
	case TagCompound:
		return "TAG_Compound"
	case TagIntArray:
		return "TAG_Int_Array"
	case TagLongArray:
		return "TAG_Long_Array"
	default:
		return fmt.Sprintf("Unknown(%d)", t)
	}
}

// Tag represents an NBT tag with its value.
type Tag struct {
	Name  string
	Type  TagType
	Value interface{}
}

// NewTag creates a new named tag.
func NewTag(name string, typ TagType, value interface{}) Tag {
	return Tag{
		Name:  name,
		Type:  typ,
		Value: value,
	}
}

// Decoder reads and decodes NBT data from an input stream.
type Decoder struct {
	r io.Reader
}

// NewDecoder creates a new NBT decoder.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Decode reads and decodes the next NBT tag from the input.
// For unnamed tags (in lists), use DecodeUnnmamed.
func (d *Decoder) Decode() (Tag, error) {
	typ, err := d.readByte()
	if err != nil {
		return Tag{}, err
	}

	if typ == byte(TagEnd) {
		return Tag{Type: TagEnd}, nil
	}

	name, err := d.readString()
	if err != nil {
		return Tag{}, fmt.Errorf("failed to read tag name: %w", err)
	}

	value, err := d.readValue(TagType(typ))
	if err != nil {
		return Tag{}, fmt.Errorf("failed to read tag value: %w", err)
	}

	return Tag{
		Name:  name,
		Type:  TagType(typ),
		Value: value,
	}, nil
}

// DecodeUnnamed reads an unnamed tag (used in lists).
func (d *Decoder) DecodeUnnamed(listType TagType) (Tag, error) {
	value, err := d.readValue(listType)
	if err != nil {
		return Tag{}, err
	}

	return Tag{
		Type:  listType,
		Value: value,
	}, nil
}

// readValue reads the value for a given tag type.
func (d *Decoder) readValue(typ TagType) (interface{}, error) {
	switch typ {
	case TagEnd:
		return nil, nil
	case TagByte:
		v, err := d.readByte()
		if err != nil {
			return nil, err
		}
		return int8(v), nil
	case TagShort:
		v, err := d.readInt16()
		if err != nil {
			return nil, err
		}
		return v, nil
	case TagInt:
		v, err := d.readInt32()
		if err != nil {
			return nil, err
		}
		return v, nil
	case TagLong:
		v, err := d.readInt64()
		if err != nil {
			return nil, err
		}
		return v, nil
	case TagFloat:
		v, err := d.readFloat32()
		if err != nil {
			return nil, err
		}
		return v, nil
	case TagDouble:
		v, err := d.readFloat64()
		if err != nil {
			return nil, err
		}
		return v, nil
	case TagString:
		return d.readString()
	case TagByteArray:
		length, err := d.readInt32()
		if err != nil {
			return nil, err
		}
		if length < 0 {
			return nil, errors.New("negative byte array length")
		}
		data := make([]byte, length)
		_, err = io.ReadFull(d.r, data)
		if err != nil {
			return nil, err
		}
		return data, nil
	case TagIntArray:
		length, err := d.readInt32()
		if err != nil {
			return nil, err
		}
		if length < 0 {
			return nil, errors.New("negative int array length")
		}
		arr := make([]int32, length)
		for i := range arr {
			arr[i], err = d.readInt32()
			if err != nil {
				return nil, err
			}
		}
		return arr, nil
	case TagLongArray:
		length, err := d.readInt32()
		if err != nil {
			return nil, err
		}
		if length < 0 {
			return nil, errors.New("negative long array length")
		}
		arr := make([]int64, length)
		for i := range arr {
			arr[i], err = d.readInt64()
			if err != nil {
				return nil, err
			}
		}
		return arr, nil
	case TagList:
		listType, err := d.readByte()
		if err != nil {
			return nil, err
		}
		length, err := d.readInt32()
		if err != nil {
			return nil, err
		}
		if length < 0 {
			return nil, errors.New("negative list length")
		}
		list := make([]interface{}, length)
		for i := range list {
			tag, err := d.DecodeUnnamed(TagType(listType))
			if err != nil {
				return nil, err
			}
			list[i] = tag.Value
		}
		return list, nil
	case TagCompound:
		compound := make(map[string]interface{})
		for {
			tag, err := d.Decode()
			if err != nil {
				return nil, err
			}
			if tag.Type == TagEnd {
				break
			}
			compound[tag.Name] = tag.Value
		}
		return compound, nil
	default:
		return nil, fmt.Errorf("unknown tag type: %d", typ)
	}
}

// readByte reads a single byte.
func (d *Decoder) readByte() (byte, error) {
	var b [1]byte
	_, err := io.ReadFull(d.r, b[:])
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

// readInt16 reads a big-endian int16.
func (d *Decoder) readInt16() (int16, error) {
	var b [2]byte
	_, err := io.ReadFull(d.r, b[:])
	if err != nil {
		return 0, err
	}
	return int16(binary.BigEndian.Uint16(b[:])), nil
}

// readInt32 reads a big-endian int32.
func (d *Decoder) readInt32() (int32, error) {
	var b [4]byte
	_, err := io.ReadFull(d.r, b[:])
	if err != nil {
		return 0, err
	}
	return int32(binary.BigEndian.Uint32(b[:])), nil
}

// readInt64 reads a big-endian int64.
func (d *Decoder) readInt64() (int64, error) {
	var b [8]byte
	_, err := io.ReadFull(d.r, b[:])
	if err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(b[:])), nil
}

// readFloat32 reads a big-endian float32.
func (d *Decoder) readFloat32() (float32, error) {
	bits, err := d.readInt32()
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(uint32(bits)), nil
}

// readFloat64 reads a big-endian float64.
func (d *Decoder) readFloat64() (float64, error) {
	bits, err := d.readInt64()
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(uint64(bits)), nil
}

// readString reads a NBT string (length-prefixed UTF-8).
func (d *Decoder) readString() (string, error) {
	length, err := d.readInt16()
	if err != nil {
		return "", err
	}
	if length < 0 {
		return "", errors.New("negative string length")
	}
	if length == 0 {
		return "", nil
	}
	data := make([]byte, length)
	_, err = io.ReadFull(d.r, data)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Encoder writes NBT data to an output stream.
type Encoder struct {
	w       io.Writer
scratch [8]byte
}

// NewEncoder creates a new NBT encoder.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Encode writes a named tag to the output.
func (e *Encoder) Encode(tag Tag) error {
	if err := e.writeByte(byte(tag.Type)); err != nil {
		return err
	}

	if tag.Type != TagEnd {
		if err := e.writeString(tag.Name); err != nil {
			return err
		}
	}

	return e.writeValue(tag.Type, tag.Value)
}

// writeValue writes a tag value to the output.
func (e *Encoder) writeValue(typ TagType, value interface{}) error {
	switch typ {
	case TagEnd:
		return nil
	case TagByte:
		return e.writeByte(byte(value.(int8)))
	case TagShort:
		return e.writeInt16(int16(value.(int16)))
	case TagInt:
		return e.writeInt32(int32(value.(int32)))
	case TagLong:
		return e.writeInt64(int64(value.(int64)))
	case TagFloat:
		bits := math.Float32bits(value.(float32))
		return e.writeInt32(int32(bits))
	case TagDouble:
		bits := math.Float64bits(value.(float64))
		return e.writeInt64(int64(bits))
	case TagString:
		return e.writeString(value.(string))
	case TagByteArray:
		data := value.([]byte)
		if err := e.writeInt32(int32(len(data))); err != nil {
			return err
		}
		_, err := e.w.Write(data)
		return err
	case TagIntArray:
		arr := value.([]int32)
		if err := e.writeInt32(int32(len(arr))); err != nil {
			return err
		}
		for _, v := range arr {
			if err := e.writeInt32(v); err != nil {
				return err
			}
		}
		return nil
	case TagLongArray:
		arr := value.([]int64)
		if err := e.writeInt32(int32(len(arr))); err != nil {
			return err
		}
		for _, v := range arr {
			if err := e.writeInt64(v); err != nil {
				return err
			}
		}
		return nil
	case TagList:
		list := value.([]interface{})
		if len(list) == 0 {
			// Write empty list type as TAG_End
			if err := e.writeByte(byte(TagEnd)); err != nil {
				return err
			}
			return e.writeInt32(0)
		}
		// Get type from first element
		listType := TagByte // Default
		if len(list) > 0 {
			// Try to determine type from first element
			switch list[0].(type) {
			case int8:
				listType = TagByte
			case int16:
				listType = TagShort
			case int32:
				listType = TagInt
			case int64:
				listType = TagLong
			case float32:
				listType = TagFloat
			case float64:
				listType = TagDouble
			case string:
				listType = TagString
			case []byte:
				listType = TagByteArray
			case []int32:
				listType = TagIntArray
			case []int64:
				listType = TagLongArray
			case []interface{}:
				listType = TagList
			case map[string]interface{}:
				listType = TagCompound
			}
		}
		if err := e.writeByte(byte(listType)); err != nil {
			return err
		}
		if err := e.writeInt32(int32(len(list))); err != nil {
			return err
		}
		for _, v := range list {
			if err := e.writeValue(listType, v); err != nil {
				return err
			}
		}
		return nil
	case TagCompound:
		compound := value.(map[string]interface{})
		for name, val := range compound {
			tag := e.inferTag(name, val)
			if err := e.Encode(tag); err != nil {
				return err
			}
		}
		// Write TAG_End
		return e.writeByte(byte(TagEnd))
	default:
		return fmt.Errorf("unknown tag type: %d", typ)
	}
}

// inferTag infers the tag type from a Go value.
func (e *Encoder) inferTag(name string, value interface{}) Tag {
	switch v := value.(type) {
	case int8:
		return NewTag(name, TagByte, v)
	case int16:
		return NewTag(name, TagShort, v)
	case int32:
		return NewTag(name, TagInt, v)
	case int64:
		return NewTag(name, TagLong, v)
	case float32:
		return NewTag(name, TagFloat, v)
	case float64:
		return NewTag(name, TagDouble, v)
	case string:
		return NewTag(name, TagString, v)
	case []byte:
		return NewTag(name, TagByteArray, v)
	case []int32:
		return NewTag(name, TagIntArray, v)
	case []int64:
		return NewTag(name, TagLongArray, v)
	case []interface{}:
		return NewTag(name, TagList, v)
	case map[string]interface{}:
		return NewTag(name, TagCompound, v)
	default:
		// Try to convert int types
		switch v := value.(type) {
		case int:
			return NewTag(name, TagInt, int32(v))
		case uint:
			return NewTag(name, TagInt, int32(v))
		default:
			return NewTag(name, TagString, fmt.Sprintf("%v", v))
		}
	}
}

// writeByte writes a single byte.
func (e *Encoder) writeByte(b byte) error {
	e.scratch[0] = b
	_, err := e.w.Write(e.scratch[:1])
	return err
}

// writeInt16 writes a big-endian int16.
func (e *Encoder) writeInt16(v int16) error {
	binary.BigEndian.PutUint16(e.scratch[:2], uint16(v))
	_, err := e.w.Write(e.scratch[:2])
	return err
}

// writeInt32 writes a big-endian int32.
func (e *Encoder) writeInt32(v int32) error {
	binary.BigEndian.PutUint32(e.scratch[:4], uint32(v))
	_, err := e.w.Write(e.scratch[:4])
	return err
}

// writeInt64 writes a big-endian int64.
func (e *Encoder) writeInt64(v int64) error {
	binary.BigEndian.PutUint64(e.scratch[:8], uint64(v))
	_, err := e.w.Write(e.scratch[:8])
	return err
}

// writeString writes a NBT string.
func (e *Encoder) writeString(s string) error {
	data := []byte(s)
	if err := e.writeInt16(int16(len(data))); err != nil {
		return err
	}
	_, err := e.w.Write(data)
	return err
}

// Unmarshal parses NBT data from a byte slice.
func Unmarshal(data []byte) (Tag, error) {
	r := bytes.NewReader(data)
	decoder := NewDecoder(r)
	return decoder.Decode()
}

// Marshal serializes a tag to NBT format.
func Marshal(tag Tag) ([]byte, error) {
	var buf bytes.Buffer
	encoder := NewEncoder(&buf)
	if err := encoder.Encode(tag); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
