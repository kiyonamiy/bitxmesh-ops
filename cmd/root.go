package cmd

import (
	"git.hyperchain.cn/dmlab/bitxmesh-ops/cmd/start"
	"git.hyperchain.cn/dmlab/bitxmesh-ops/internal/utils/configutils"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bitxmesh-ops",
	Short: "The bitxmesh-ops is a tool for managing BitXMesh",
}

func init() {
	rootCmd.AddCommand(start.Command())
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// 以下设置，使得 InitConfig 函数在每个命令运行时都会被调用以读取配置
	cobra.OnInitialize(configutils.InitConfig)
	// Cobra 支持持久性标志(PersistentFlag)，该标志可用于它所分配的命令以及该命令下的每个子命令
	rootCmd.PersistentFlags().StringVarP(&configutils.CfgFile, "config", "c", "", "The path to the bitxmesh-ops configuration file. Empty string for no configuration file.")

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
