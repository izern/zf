package util

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {

}

func ExactArgsWithPipe(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if IsPipe() {
			if (len(args) + 1) != n {
				return fmt.Errorf("accepts %d arg(s), received %d", n, len(args)+1)
			}
		} else {
			if len(args) != n {
				return fmt.Errorf("accepts %d arg(s), received %d", n, len(args))
			}
		}
		return nil
	}
}

func InitArgsFromPipe(args []string) ([]string, error) {
	if IsPipe() {
		if args == nil {
			args = make([]string, 0)
		}
		arg, err := ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		result := make([]string, len(args)+1)
		result[0] = arg
		result = append(result, args...)
		return result, nil
	}
	return args, nil
}
