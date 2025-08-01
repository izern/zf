package codec

import "github.com/izern/zf/types"

// Marshaler defines the interface for encoding data to bytes
type Marshaler interface {
	// Marshal encodes the given data structure into bytes
	// Returns the encoded bytes or an error if encoding fails
	Marshal(data interface{}) ([]byte, types.ZfError)
}

// CodecCapabilities defines what features a codec supports
type CodecCapabilities struct {
	SupportsObjects bool // Can handle object/map structures
	SupportsArrays  bool // Can handle array/slice structures
	SupportsPrimitives bool // Can handle primitive types (string, number, bool)
	SupportsComments bool // Can handle comments in the format
}

// CodecInfo provides metadata about a codec
type CodecInfo struct {
	Name         string
	FileExtensions []string
	MimeTypes    []string
	Capabilities CodecCapabilities
}
