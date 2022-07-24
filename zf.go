package main

import (
	"fmt"
	"github.com/izern/zf/cmd"
	_ "github.com/izern/zf/cmd"
	"github.com/izern/zf/types"
	"github.com/izern/zf/util"
	"github.com/spf13/cobra"
	"math"
)

var pretty bool

func init() {
}

func main() {

	rootCmd := &cobra.Command{
		Use:     "zf",
		Short:   "zf用来解析格式化字符串文本",
		Example: "cat file.yml | zf yaml ",
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "版本号",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("zf version: v0.9")
		},
	})

	typeCmds := cmd.GetAllCmd()
	for _, typeCmd := range typeCmds {
		cmd := &cobra.Command{
			Use:   typeCmd.GetCurrType(),
			Short: fmt.Sprintf("解析%s格式的文本", typeCmd.GetCurrType()),
			Args:  util.ExactArgsWithPipe(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				args, e := util.InitArgsFromPipe(args)
				if e != nil {
					return e
				}
				text, err := typeCmd.Parse(args[0])
				if err != nil {
					return err.Error()
				}
				fmt.Println(text)
				return nil
			},
		}
		rootCmd.AddCommand(cmd)
		appendChildCmd(cmd, typeCmd)
	}

	var from, to string
	convertCmd := &cobra.Command{
		Use:     "convert",
		Short:   "文本内容格式转换",
		Example: "cat test.yml | zf --from yaml --to json",
		Args:    util.ExactArgsWithPipe(1),
		RunE: func(c *cobra.Command, args []string) error {
			args, err := util.InitArgsFromPipe(args)
			if err != nil {
				return err
			}
			fromCmd, err := cmd.GetCmd(from)
			if err != nil {
				return err
			}
			toCmd, err := cmd.GetCmd(to)
			if err != nil {
				return err
			}

			res, e := fromCmd.GetValues(0, math.MaxUint32, ".", args[0])
			if e != nil {
				return e.Error()
			}

			switch res.(type) {
			case map[string]interface{}:
				res = res.(map[string]interface{})
			case map[interface{}]interface{}:
				res = util.ConvertMap2String(res.(map[interface{}]interface{}))
			case []interface{}:
				res = util.ConvertArray2String(res.([]interface{}))
			}
			text, e := toCmd.Marshal(res)
			if e != nil {
				return e.Error()
			}
			fmt.Println(text)
			return nil
		},
	}
	convertCmd.Flags().StringVarP(&from, "from", "f", "", "源数据格式")
	convertCmd.Flags().StringVarP(&to, "to", "t", "", "目标数据格式")

	rootCmd.AddCommand(convertCmd)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err.Error())
	}

}

func appendChildCmd(cmd *cobra.Command, typeCmd types.TypeCommand) {
	appendParseCmd(cmd, typeCmd)
	appendGetTypeCmd(cmd, typeCmd)
	appendKeysCmd(cmd, typeCmd)
	appendAppendCmd(cmd, typeCmd)
	appendGetValueCmd(cmd, typeCmd)
	appendSetValueCmd(cmd, typeCmd)
}

func appendGetTypeCmd(cmd *cobra.Command, typeCmd types.TypeCommand) {
	var path string

	c := &cobra.Command{
		Use:   "type",
		Short: "获取指定路径值的类别",
		Args:  util.ExactArgsWithPipe(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			args, e := util.InitArgsFromPipe(args)
			if e != nil {
				return e
			}

			typeStr, err := typeCmd.GetType(path, args[0])
			if err != nil {
				return err.Error()
			}
			fmt.Println(typeStr)
			return nil
		},
	}
	c.Flags().StringVarP(&path, "path", "p", ".", "节点路径，jsonpath格式")
	cmd.AddCommand(c)
}

func appendParseCmd(cmd *cobra.Command, typeCmd types.TypeCommand) {
	c := &cobra.Command{
		Use:   "parse",
		Short: "格式化",
		Args:  util.ExactArgsWithPipe(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			args, e := util.InitArgsFromPipe(args)
			if e != nil {
				return e
			}
			text, err := typeCmd.Parse(args[0])
			if err != nil {
				return err.Error()
			}
			fmt.Println(text)
			return nil
		},
	}
	cmd.AddCommand(c)
}

func appendKeysCmd(cmd *cobra.Command, typeCmd types.TypeCommand) {
	var from, to uint
	var path string
	c := &cobra.Command{
		Use:   "keys",
		Short: "获取键列表",
		Args:  util.ExactArgsWithPipe(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			args, e := util.InitArgsFromPipe(args)
			if e != nil {
				return e
			}
			keys, err := typeCmd.Keys(from, to, path, args[0])
			if err != nil {
				return err.Error()
			}
			for _, key := range keys {
				fmt.Println(key)
			}
			return nil
		},
	}

	c.Flags().UintVarP(&from, "from", "f", 0, "范围起始值from")
	c.Flags().UintVarP(&to, "to", "t", math.MaxInt16, "范围终止值to")
	c.Flags().StringVarP(&path, "path", "p", ".", "节点路径，jsonpath格式")
	cmd.AddCommand(c)
}

func appendGetValueCmd(cmd *cobra.Command, typeCmd types.TypeCommand) {
	var from, to uint
	var path string
	c := &cobra.Command{
		Use:   "get",
		Short: "获取值",
		Args:  util.ExactArgsWithPipe(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			args, e := util.InitArgsFromPipe(args)
			if e != nil {
				return e
			}
			res, err := typeCmd.GetValues(from, to, path, args[0])
			if err != nil {
				return err.Error()
			}
			marshal, zfError := typeCmd.Marshal(res)
			fmt.Println(marshal)
			if zfError != nil {
				zfError.Error()
			}
			return nil
		},
	}
	c.Flags().UintVarP(&from, "from", "f", 0, "范围起始值from")
	c.Flags().UintVarP(&to, "to", "t", math.MaxInt16, "范围终止值to")
	c.Flags().StringVarP(&path, "path", "p", ".", "节点路径，jsonpath格式")

	cmd.AddCommand(c)
}

func appendAppendCmd(cmd *cobra.Command, typeCmd types.TypeCommand) {
	var key, value, path string
	var index uint
	c := &cobra.Command{
		Use:   "append",
		Short: "追加值",
		Args:  util.ExactArgsWithPipe(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			args, e := util.InitArgsFromPipe(args)
			if e != nil {
				return e
			}
			text, err := typeCmd.Append(path, key, index, value, args[0])
			if err != nil {
				return err.Error()
			}
			fmt.Println(text)
			return nil
		},
	}
	c.Flags().StringVarP(&path, "path", "p", ".", "节点路径，jsonpath格式")
	c.Flags().UintVarP(&index, "index", "i", math.MaxInt16, "array或string时可以指定，默认插在最后面")
	c.Flags().StringVarP(&key, "key", "k", ".", "当类型为object时需指定key")
	c.Flags().StringVarP(&value, "value", "v", "", "append的值")
	cmd.AddCommand(c)
}

func appendSetValueCmd(cmd *cobra.Command, typeCmd types.TypeCommand) {
	var path, value string

	c := &cobra.Command{
		Use:   "set",
		Short: "修改值，覆盖",
		Args:  util.ExactArgsWithPipe(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			args, e := util.InitArgsFromPipe(args)
			if e != nil {
				return e
			}
			text, err := typeCmd.SetValue(path, value, args[0])
			if err != nil {
				return err.Error()
			}
			fmt.Println(text)
			return nil
		},
	}
	c.Flags().StringVarP(&path, "path", "p", ".", "节点路径，jsonpath格式")
	c.Flags().StringVarP(&value, "value", "v", "", "set的值")

	cmd.AddCommand(c)
}
