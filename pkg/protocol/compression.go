package protocol

import (
	"bytes"
	"compress/zlib"
	"io"
)

// Compressor handles packet compression/decompression using Zlib.
//
// Minecraft uses Zlib compression for packets larger than a certain threshold.
// The compression threshold is set by the server and can vary.
type Compressor struct {
	threshold int // Packets larger than this will be compressed
}

// NewCompressor creates a new compressor with the given threshold.
// Threshold is the minimum packet size for compression.
// Packets smaller than this are sent uncompressed but with a "data length" prefix of 0.
func NewCompressor(threshold int) *Compressor {
	return &Compressor{
		threshold: threshold,
	}
}

// Threshold returns the compression threshold.
func (c *Compressor) Threshold() int {
	return c.threshold
}

// SetThreshold sets a new compression threshold.
func (c *Compressor) SetThreshold(threshold int) {
	c.threshold = threshold
}

// Compress compresses the data if it exceeds the threshold.
// Returns a buffer with VarInt-prefixed uncompressed data length, followed by compressed data.
// If data size <= threshold, returns VarInt(0) followed by uncompressed data.
func (c *Compressor) Compress(data []byte) ([]byte, error) {
	if len(data) <= c.threshold {
		// Don't compress, but still write length prefix
		buffer := NewPacketBufferFromBytes(data)
		varIntBuf := NewPacketBuffer()
		varIntBuf.WriteVarInt(0)
		varIntBuf.WriteBytes(data)
		return varIntBuf.Bytes(), nil
	}

	// Compress the data
	var compressed bytes.Buffer
	writer := zlib.NewWriter(&compressed)

	if _, err := writer.Write(data); err != nil {
		writer.Close()
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	// Write VarInt uncompressed data length, then compressed data
	buffer := NewPacketBuffer()
	buffer.WriteVarInt(len(data))
	buffer.WriteBytes(compressed.Bytes())

	return buffer.Bytes(), nil
}

// Decompress decompresses data that was compressed by Compress.
// Reads a VarInt prefix for uncompressed length, then decompresses the remaining data.
// If the prefix is 0, the data is not compressed and is returned as-is.
func (c *Compressor) Decompress(data []byte) ([]byte, error) {
	r := bytes.NewReader(data)

	// Read uncompressed data length (VarInt)
	uncompressedLength, err := ReadVarInt(r)
	if err != nil {
		return nil, err
	}

	if uncompressedLength == 0 {
		// Data is not compressed, return remaining bytes
		remaining, _ := io.ReadAll(r)
		return remaining, nil
	}

	// Read compressed data
	compressedData, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Decompress
	reader, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// Read decompressed data
	decompressed := bytes.NewBuffer(nil)
	decompressed.Grow(uncompressedLength)

	if _, err := io.CopyN(decompressed, reader, int64(uncompressedLength)); err != nil {
		return nil, err
	}

	return decompressed.Bytes(), nil
}

// CompressWriter wraps an io.Writer and compresses data written to it.
type CompressWriter struct {
	writer   io.Writer
	compressor *zlib.Writer
}

// NewCompressWriter creates a new compress writer.
func NewCompressWriter(writer io.Writer) (*CompressWriter, error) {
	compressor := zlib.NewWriter(writer)
	return &CompressWriter{
		writer:   writer,
		compressor: compressor,
	}, nil
}

// Write writes data to the compressed writer.
func (cw *CompressWriter) Write(data []byte) (int, error) {
	return cw.compressor.Write(data)
}

// Close closes the compressor and flushes any remaining data.
func (cw *CompressWriter) Close() error {
	return cw.compressor.Close()
}

// DecompressReader wraps an io.Reader and decompresses data read from it.
type DecompressReader struct {
	reader     io.Reader
	decompressor *zlib.Reader
}

// NewDecompressReader creates a new decompress reader.
func NewDecompressReader(reader io.Reader) (*DecompressReader, error) {
	decompressor, err := zlib.NewReader(reader)
	if err != nil {
		return nil, err
	}

	return &DecompressReader{
		reader:      reader,
		decompressor: decompressor,
	}, nil
}

// Read reads decompressed data from the reader.
func (dr *DecompressReader) Read(data []byte) (int, error) {
	return dr.decompressor.Read(data)
}

// Close closes the decompressor.
func (dr *DecompressReader) Close() error {
	return dr.decompressor.Close()
}

// IsCompressed checks if data is compressed by checking the Zlib header.
// Zlib data starts with byte 0x78 (deflate, 32K window).
func IsCompressed(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	// Zlib header: first byte is 0x78 for deflate compression
	// Second byte can vary (check bits 0-4 for checksum)
	return (data[0] & 0x0F) == 0x08 && (data[0]>>4) == 0x07
}

// GetCompressionLevel returns the compression level based on Zlib header.
func GetCompressionLevel(data []byte) int {
	if len(data) == 0 {
		return 0
	}

	// Zlib compression level is encoded in the second byte
	if len(data) >= 2 {
		level := data[1] & 0xC0 >> 6
		switch level {
		case 0:
			return zlib.NoCompression
		case 1:
			return zlib.BestSpeed
		case 2:
			return 0 // Default compression
		case 3:
			return zlib.BestCompression
		}
	}

	return zlib.DefaultCompression
}

// CompressFast compresses data with best speed.
func CompressFast(data []byte) ([]byte, error) {
	var compressed bytes.Buffer
	writer, err := zlib.NewWriterLevel(&compressed, zlib.BestSpeed)
	if err != nil {
		return nil, err
	}

	if _, err := writer.Write(data); err != nil {
		writer.Close()
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return compressed.Bytes(), nil
}

// CompressBest compresses data with best compression.
func CompressBest(data []byte) ([]byte, error) {
	var compressed bytes.Buffer
	writer, err := zlib.NewWriterLevel(&compressed, zlib.BestCompression)
	if err != nil {
		return nil, err
	}

	if _, err := writer.Write(data); err != nil {
		writer.Close()
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return compressed.Bytes(), nil
}

// DecompressToBytes decompresses data and returns as byte slice.
func DecompressToBytes(data []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	decompressed, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return decompressed, nil
}

// DecompressString decompresses data and returns as string.
func DecompressString(data []byte) (string, error) {
	decompressed, err := DecompressToBytes(data)
	if err != nil {
		return "", err
	}
	return string(decompressed), nil
}
