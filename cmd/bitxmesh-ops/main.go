package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	command := newCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

func newCommand() *cobra.Command {
	cmd := &cobra.Command{
		// 作为主命令，Use 只影响 -h 的提示输出，并不影响命令本身 ./_output/bitxmesh-ops；而如果是子命令，则会作为子命令的命令名
		Use:   "bitxmesh-ops",
		Short: "The bitxmesh-ops is a tool for managing BitXMesh",
		// 命令出错时，不打印帮助信息。不需要打印帮助信息，设置为 true 可以保持命令出错时一眼就能看到错误信息
		// SilenceUsage: true,
		// 指定调用 cmd.Execute() 时，执行的 Run 函数，函数执行失败会返回错误信息（由 cobra 打印）
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Hello World")
			return nil
		},
		// 这里设置命令运行时，不需要指定命令行参数
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		},
	}

	return cmd
}
