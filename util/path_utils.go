package util

import (
	"github.com/izern/zf/types"
	"math"
	"strconv"
	"strings"
)

func init() {

}

// ParsePath parses a JSONPath-style path string into Path objects
// Enhanced with better validation and error messages
func ParsePath(path string) ([]*types.Path, types.ZfError) {
	if err := validatePath(path); err != nil {
		return nil, err
	}

	splits, err := splitPath(path)
	if err != nil {
		return nil, err
	}

	paths := make([]*types.Path, 0, len(splits))
	for i, str := range splits {
		pathNode, parseErr := parsePathNode(str, i == 0)
		if parseErr != nil {
			return nil, parseErr
		}
		paths = append(paths, pathNode)
	}

	return paths, nil
}

// validatePath performs basic validation on the path string
func validatePath(path string) types.ZfError {
	if path == "" {
		return types.NewFormatError(path, "path: empty path not allowed")
	}
	
	// Check for invalid endings (except escaped dots)
	if strings.HasSuffix(path, ".") && !strings.HasSuffix(path, "\\.") {
		return types.NewFormatError(path, "path: cannot end with '.'")
	}
	
	// Check for consecutive dots
	if strings.Contains(path, "..") {
		return types.NewFormatError(path, "path: consecutive dots '..' not allowed")
	}
	
	return nil
}

// parsePathNode parses a single path node (e.g., "name", "items[0]", "data[1,5]")
func parsePathNode(str string, isRoot bool) (*types.Path, types.ZfError) {
	// Handle root node
	if isRoot && (str == "$" || str == "") {
		return &types.Path{
			Type:        types.RootNode,
			NodeKey:     "$",
			OriginValue: str,
		}, nil
	}

	// Find array/index notation
	rangeStart := strings.Index(str, "[")
	rangeEnd := strings.LastIndex(str, "]")
	
	// Validate bracket pairing
	if rangeStart != -1 && rangeEnd == -1 {
		return nil, types.NewFormatError(str, "path: missing closing bracket ']'")
	}
	if rangeStart == -1 && rangeEnd != -1 {
		return nil, types.NewFormatError(str, "path: missing opening bracket '['")
	}
	if rangeStart > rangeEnd {
		return nil, types.NewFormatError(str, "path: malformed brackets")
	}

	// No brackets - normal node
	if rangeStart == -1 {
		return &types.Path{
			Type:        types.NormalNode,
			NodeKey:     str,
			OriginValue: str,
		}, nil
	}

	// Parse node with brackets
	nodeKey := str[:rangeStart]
	rangeContent := str[rangeStart+1 : rangeEnd]

	// Empty brackets [] - range all
	if rangeContent == "" {
		return &types.Path{
			Type:        types.RangeNode,
			NodeKey:     nodeKey,
			From:        0,
			To:          math.MaxInt16,
			OriginValue: str,
		}, nil
	}

	// Parse range content
	if strings.Contains(rangeContent, ",") {
		return parseRangeNode(nodeKey, rangeContent, str)
	} else {
		return parseIndexNode(nodeKey, rangeContent, str)
	}
}

// parseRangeNode parses range notation like [1,5]
func parseRangeNode(nodeKey, rangeContent, originalStr string) (*types.Path, types.ZfError) {
	parts := strings.Split(rangeContent, ",")
	if len(parts) != 2 {
		return nil, types.NewFormatError(originalStr, "path: range must have exactly two parts separated by comma")
	}

	start, err := parseUint(strings.TrimSpace(parts[0]), originalStr)
	if err != nil {
		return nil, err
	}

	end, err := parseUint(strings.TrimSpace(parts[1]), originalStr)
	if err != nil {
		return nil, err
	}

	if start > end {
		return nil, types.NewFormatError(originalStr, "path: range start cannot be greater than end")
	}

	return &types.Path{
		Type:        types.RangeNode,
		NodeKey:     nodeKey,
		From:        start,
		To:          end,
		OriginValue: originalStr,
	}, nil
}

// parseIndexNode parses index notation like [3]
func parseIndexNode(nodeKey, indexContent, originalStr string) (*types.Path, types.ZfError) {
	index, err := parseUint(strings.TrimSpace(indexContent), originalStr)
	if err != nil {
		return nil, err
	}

	return &types.Path{
		Type:        types.IndexNode,
		NodeKey:     nodeKey,
		Index:       index,
		OriginValue: originalStr,
	}, nil
}

// parseUint safely parses an unsigned integer with better error messages
func parseUint(s, context string) (uint, types.ZfError) {
	if s == "" {
		return 0, types.NewFormatError(context, "path: empty number not allowed")
	}

	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, types.NewFormatError(context, "path: invalid number '"+s+"'")
	}

	return uint(val), nil
}

func splitPath(path string) ([]string, types.ZfError) {
	if path == "" {
		return nil, types.NewFormatError(path, "path: empty path")
	}
	
	if path == "." {
		return []string{"$"}, nil
	}

	// Split by dots, but handle escaped dots
	parts := strings.Split(path, ".")
	
	// Process escaped dots
	result := make([]string, 0, len(parts))
	i := 0
	for i < len(parts) {
		part := parts[i]
		if i == 0 && part == "" {
			// Leading dot is converted to root
			result = append(result, "$")
			i++
			continue
		}
		
		// Handle escaped dots by rejoining with next part
		if strings.HasSuffix(part, "\\") && i < len(parts)-1 {
			// This is an escaped dot, combine with next part
			combined := part[:len(part)-1] + "." + parts[i+1]
			result = append(result, combined)
			i += 2 // Skip next part as it's already processed
		} else {
			result = append(result, part)
			i++
		}
	}

	return result, nil
}
