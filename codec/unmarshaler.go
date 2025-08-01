package codec

import "github.com/izern/zf/types"

// Unmarshaler defines the interface for decoding data from bytes
type Unmarshaler interface {
	// Unmarshal decodes bytes into a data structure
	// Returns the decoded data or an error if decoding fails
	Unmarshal(data []byte) (interface{}, types.ZfError)
}

// Codec combines both marshaling and unmarshaling capabilities
type Codec interface {
	Marshaler
	Unmarshaler
	
	// GetInfo returns metadata about this codec
	GetInfo() CodecInfo
	
	// CanHandle checks if this codec can handle the given input
	// This can be used for automatic format detection
	CanHandle(data []byte) bool
}
