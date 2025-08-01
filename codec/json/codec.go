package json

import (
	"encoding/json"
	"github.com/izern/zf/codec"
	"github.com/izern/zf/types"
	"strings"
)

func init() {

}

type JSONCodec struct {
}

func (j *JSONCodec) Marshal(data interface{}) ([]byte, types.ZfError) {
	result, e := json.Marshal(data)
	if e != nil {
		return nil, types.NewFormatError(e.Error(), "json")
	}
	return result, nil
}

func (j *JSONCodec) Unmarshal(data []byte) (interface{}, types.ZfError) {
	if len(data) == 0 {
		return nil, nil
	}
	
	var result interface{}
	e := json.Unmarshal(data, &result)
	if e != nil {
		return nil, types.NewFormatError(e.Error(), "json")
	}
	return result, nil
}

func (j *JSONCodec) GetInfo() codec.CodecInfo {
	return codec.CodecInfo{
		Name:           "json",
		FileExtensions: []string{".json", ".js"},
		MimeTypes:      []string{"application/json", "text/json"},
		Capabilities: codec.CodecCapabilities{
			SupportsObjects:    true,
			SupportsArrays:     true,
			SupportsPrimitives: true,
			SupportsComments:   false,
		},
	}
}

func (j *JSONCodec) CanHandle(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	
	// Trim whitespace
	content := strings.TrimSpace(string(data))
	if len(content) == 0 {
		return false
	}
	
	// JSON typically starts with { [ " or is a primitive value
	firstChar := content[0]
	switch firstChar {
	case '{', '[', '"':
		return true
	case 't', 'f': // true/false
		return strings.HasPrefix(content, "true") || strings.HasPrefix(content, "false")
	case 'n': // null
		return strings.HasPrefix(content, "null")
	default:
		// Try to parse as number
		var num json.Number
		return json.Unmarshal([]byte(content), &num) == nil
	}
}
