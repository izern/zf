package toml

import (
	"encoding/json"
	"github.com/izern/zf/codec"
	"github.com/izern/zf/types"
	"github.com/pelletier/go-toml/v2"
	"strings"
)

func init() {

}

type TomlCodec struct {
}

func (t *TomlCodec) Marshal(data interface{}) ([]byte, types.ZfError) {
	dataType, err := types.GetType(data)
	if err != nil {
		return nil, err
	}
	
	var result []byte
	var e error
	
	switch dataType {
	case types.Object:
		result, e = toml.Marshal(data)
	default:
		// TOML can't handle primitives at root level, use JSON fallback
		result, e = json.Marshal(data)
	}
	
	if e != nil {
		return nil, types.NewFormatError(e.Error(), "toml")
	}
	return result, nil
}

func (t *TomlCodec) Unmarshal(data []byte) (interface{}, types.ZfError) {
	if len(data) == 0 {
		return nil, nil
	}
	
	var result interface{}
	result = make(map[interface{}]interface{})
	e := toml.Unmarshal(data, &result)
	if e != nil {
		// Try as JSON fallback for primitives
		result = make([]interface{}, 0)
		e = json.Unmarshal(data, &result)
		if e != nil {
			// If both fail, return as string
			return string(data), nil
		}
	}
	return result, nil
}

func (t *TomlCodec) GetInfo() codec.CodecInfo {
	return codec.CodecInfo{
		Name:           "toml",
		FileExtensions: []string{".toml", ".tml"},
		MimeTypes:      []string{"application/toml", "text/toml"},
		Capabilities: codec.CodecCapabilities{
			SupportsObjects:    true,
			SupportsArrays:     true,
			SupportsPrimitives: false, // TOML requires objects at root
			SupportsComments:   true,
		},
	}
}

func (t *TomlCodec) CanHandle(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	
	content := strings.TrimSpace(string(data))
	if len(content) == 0 {
		return false
	}
	
	lines := strings.Split(content, "\n")
	tomlFeatures := 0
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// TOML section headers
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			tomlFeatures += 2
		}
		
		// Key-value assignments with =
		if strings.Contains(line, "=") && !strings.Contains(line, "==") {
			tomlFeatures++
		}
		
		// TOML array syntax
		if strings.Contains(line, "[[") && strings.Contains(line, "]]") {
			tomlFeatures += 2
		}
		
		// Comments
		if strings.HasPrefix(line, "#") {
			tomlFeatures++
		}
		
		// Multi-line strings
		if strings.Contains(line, `"""`) || strings.Contains(line, `'''`) {
			tomlFeatures += 2
		}
	}
	
	// If we found multiple TOML-specific features, it's likely TOML
	return tomlFeatures >= 2
}
