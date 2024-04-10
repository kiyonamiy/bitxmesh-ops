package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a complete BitXMesh environment",
	Run: func(cmd *cobra.Command, args []string) {
		// 打印所有的配置项及其值
		settings, _ := json.Marshal(viper.AllSettings())
		fmt.Println(string(settings))
		// 打印 db -> username 配置项的值
		fmt.Println(viper.GetString("packages.chainroll"))
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
