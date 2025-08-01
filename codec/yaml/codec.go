package yaml

import (
	"github.com/izern/zf/codec"
	"github.com/izern/zf/types"
	"gopkg.in/yaml.v3"
	"strings"
)

func init() {

}

type YamlCodec struct {
}

func (y *YamlCodec) Marshal(data interface{}) ([]byte, types.ZfError) {
	result, e := yaml.Marshal(data)
	if e != nil {
		return nil, types.NewFormatError(e.Error(), "yaml")
	}
	return result, nil
}

func (y *YamlCodec) Unmarshal(data []byte) (interface{}, types.ZfError) {
	if len(data) == 0 {
		return nil, nil
	}
	
	var result interface{}
	// Try to unmarshal as a map first
	var tmp = make(map[string]interface{})
	e := yaml.Unmarshal(data, &tmp)
	if e != nil {
		// If that fails, try as generic interface
		e = yaml.Unmarshal(data, &result)
		if e != nil {
			return nil, types.NewFormatError(e.Error(), "yaml")
		}
	} else {
		result = tmp
	}
	return result, nil
}

func (y *YamlCodec) GetInfo() codec.CodecInfo {
	return codec.CodecInfo{
		Name:           "yaml",
		FileExtensions: []string{".yaml", ".yml"},
		MimeTypes:      []string{"application/yaml", "text/yaml", "application/x-yaml"},
		Capabilities: codec.CodecCapabilities{
			SupportsObjects:    true,
			SupportsArrays:     true,
			SupportsPrimitives: true,
			SupportsComments:   true,
		},
	}
}

func (y *YamlCodec) CanHandle(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	
	content := strings.TrimSpace(string(data))
	if len(content) == 0 {
		return false
	}
	
	// YAML characteristics:
	// - Contains ":" for key-value pairs
	// - Starts with "---" document separator
	// - Contains "-" for arrays
	// - Has comments with "#"
	// - Uses indentation
	
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Document separator
		if strings.HasPrefix(line, "---") {
			return true
		}
		
		// Key-value pairs (but not URLs or JSON)
		if strings.Contains(line, ":") && !strings.HasPrefix(line, "http") && !strings.Contains(content, "{") {
			return true
		}
		
		// Array items
		if strings.HasPrefix(line, "- ") {
			return true
		}
		
		// Comments
		if strings.HasPrefix(line, "#") {
			return true
		}
	}
	
	// Check for indentation pattern (multiple lines with consistent spacing)
	indentedLines := 0
	for _, line := range lines {
		if len(line) > 0 && (line[0] == ' ' || line[0] == '\t') {
			indentedLines++
		}
	}
	
	// If we have multiple indented lines, it's likely YAML
	return indentedLines > 1
}
