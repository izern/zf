package util

import (
	"github.com/izern/zf/types"
	"math"
	"strconv"
	"strings"
)

func init() {

}

// ParsePath 解析路径
func ParsePath(path string) ([]*types.Path, types.ZfError) {
	splits, err := splitPath(path)
	if err != nil {
		return nil, err
	}
	var paths []*types.Path
	for i, str := range splits {
		if i == 0 && (str == "$" || str == "") {
			path := &types.Path{
				Type:        types.RootNode,
				NodeKey:     "$",
				OriginValue: str,
			}
			paths = append(paths, path)
		} else {
			rangeStrStart := strings.Index(str, "[")
			rangeStrEnd := strings.Index(str, "]")
			// 包含[]，分三种情况，[]/[1]/[1,2]
			if rangeStrStart > -1 && rangeStrEnd > rangeStrStart {
				if rangeStrEnd-1 == rangeStrStart {
					path := &types.Path{
						Type:        types.RangeNode,
						NodeKey:     str[:rangeStrStart],
						From:        0,
						To:          math.MaxInt16,
						OriginValue: str,
					}
					paths = append(paths, path)
					continue
				}

				rangeNum := str[rangeStrStart+1 : rangeStrEnd]
				sepIndex := strings.Index(rangeNum, ",")
				// 包含逗号，为范围
				if sepIndex > -1 {
					numbers := strings.Split(rangeNum, ",")
					if len(numbers) != 2 {
						return nil, types.NewFormatError(path, "path")
					}
					start, err := strconv.ParseUint(numbers[0], 0, 0)
					if err != nil {
						return nil, types.NewFormatError(path, "path")
					}
					end, err := strconv.ParseUint(numbers[1], 0, 0)
					if err != nil {
						return nil, types.NewFormatError(path, "path")
					}

					if start > end {
						return nil, types.NewUnSupportError(str)
					}
					path := &types.Path{
						Type:        types.RangeNode,
						NodeKey:     str[:rangeStrStart],
						From:        uint(start),
						To:          uint(end),
						OriginValue: str,
					}
					paths = append(paths, path)
				} else {
					// 不包含逗号，为精准定位
					index, err := strconv.ParseUint(rangeNum, 0, 0)
					if err != nil {
						return nil, types.NewFormatError(path, "path")
					}
					path := &types.Path{
						Type:        types.IndexNode,
						NodeKey:     str[:rangeStrStart],
						Index:       uint(index),
						OriginValue: str,
					}
					paths = append(paths, path)
				}
			} else {
				path := &types.Path{
					Type:        types.NormalNode,
					NodeKey:     str,
					OriginValue: str,
				}
				paths = append(paths, path)
			}
		}
	}

	return paths, nil
}

func splitPath(path string) ([]string, *types.FormatError) {

	if path == "" {
		return nil, types.NewFormatError(path, "path")
	}
	if path == "." {
		return []string{"$"}, nil
	}
	// 不允许以点结尾,\.除外
	if path[:len(path)] == "." && path[len(path)-1:len(path)] == "\\" {
		return nil, types.NewFormatError(path, "path")
	}

	splits := strings.Split(path, ".")

	return splits, nil
}
