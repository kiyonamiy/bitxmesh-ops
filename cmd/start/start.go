package start

import (
	"fmt"
	"git.hyperchain.cn/dmlab/bitxmesh-ops/internal/utils/configutils"
	"git.hyperchain.cn/dmlab/bitxmesh-ops/internal/utils/stringutils"
	"github.com/spf13/cobra"
	"log"
)

var allServiceTypes = []string{"chainroll", "datasys", "dataflow"}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start [chainroll|datasys|dataflow] [name]",
	Short: "Start a type of service of BitXMesh",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("requires 2 arg(s), only received %d", len(args))
		}
		serviceType := args[0]
		if !stringutils.ContainsString(allServiceTypes, serviceType) {
			return fmt.Errorf("service type must be one of %v", allServiceTypes)
		}
		serviceName := args[1]
		if !configutils.Contains(serviceType, serviceName) {
			return fmt.Errorf("service name %s not found in config file", serviceName)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		serviceType := args[0]
		serviceName := args[1]
		log.Println(serviceType, serviceName)
		if serviceType == "chainroll" {
			mongodbs := configutils.ParseMongodbs()
			hyperchains := configutils.ParseHyperchains()
			chainrolls := configutils.ParseChainrolls(mongodbs, hyperchains)
			for name, chainroll := range chainrolls {
				if name == serviceName {
					chainroll.Start()
					return nil
				}
			}
		}
		if serviceType == "datasys" {
			//mongodbs := configutils.ParseMongodbs()
			//datasyss := configutils.ParseDatasyss(mongodbs)
			//for name, datasys := range datasyss {
			//	if name == serviceName {
			//		datasys.Deploy()
			//		return nil
			//	}
			//}
		}

		// TODO 应该还要增加一个 flag 为 chainroll 或者 datasys 或者 dataflow 配置参数
		//// 解析出所有可用的 mongodb
		//mongodbs := configutils.ParseMongodbs()
		//// 解析出所有可用的 hyperchain
		//hyperchains := configutils.ParseHyperchains()
		//// 解析出所有的 mesh
		//meshes := configutils.ParseMeshes(mongodbs, hyperchains)
		//// 部署所有的 mesh
		//for _, mesh := range meshes {
		//	//
		//	// TODO 应该改为 start，如果目录不存在，则创建，否则直接启动（所以不能用 log.Fatal）
		//	mesh.Deploy()
		//}
		return nil
	},
}

func Command() *cobra.Command {
	return startCmd
}
