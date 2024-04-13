package configutils

import (
	"git.hyperchain.cn/dmlab/bitxmesh-ops/internal/chainroll"
	"git.hyperchain.cn/dmlab/bitxmesh-ops/internal/hyperchain"
	"git.hyperchain.cn/dmlab/bitxmesh-ops/internal/mongodb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var CfgFile string

const (
	// recommendedHomeDir 定义放置 bitxmesh-ops 的默认目录.
	recommendedHomeDir = ".bitxmesh-ops"

	// defaultConfigName 指定了 bitxmesh-ops 的默认配置文件名.
	defaultConfigName = "bitxmesh-ops.yaml"
)

// InitConfig 设置需要读取的配置文件名、环境变量，并读取配置文件内容到 viper 中.
func InitConfig() {
	if CfgFile != "" {
		// 从命令行选项指定的配置文件中读取
		viper.SetConfigFile(CfgFile)
	} else {
		// 查找用户主目录
		home, err := os.UserHomeDir()
		// 如果获取用户主目录失败，打印 `'Error: xxx` 错误，并退出程序（退出码为 1）
		cobra.CheckErr(err)

		// 将用 `$HOME/<recommendedHomeDir>` 目录加入到配置文件的搜索路径中
		viper.AddConfigPath(filepath.Join(home, recommendedHomeDir))

		// 把当前目录加入到配置文件的搜索路径中
		viper.AddConfigPath(".")

		// 设置配置文件格式为 YAML (YAML 格式清晰易读，并且支持复杂的配置结构)
		viper.SetConfigType("yaml")

		// 配置文件名称（没有文件扩展名）
		viper.SetConfigName(defaultConfigName)
	}

	// 读取匹配的环境变量
	viper.AutomaticEnv()

	// 读取环境变量的前缀为 BITXMESH，如果是 bitxmesh，将自动转变为大写。
	viper.SetEnvPrefix("BITXMESH")

	// 以下 2 行，将 viper.Get(key) key 字符串中 '.' 和 '-' 替换为 '_'
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// 读取配置文件。如果指定了配置文件名，则使用指定的配置文件，否则在注册的搜索路径中搜索
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading configutils file: %s\n", err)
	}

	// 打印 viper 当前使用的配置文件，方便 Debug.
	log.Println("Using configutils file:", viper.ConfigFileUsed())
}

func Contains(serviceType string, serviceName string) bool {
	value := viper.Get(serviceType + "." + serviceName)
	return value != nil
}

func ParseMongodbs() map[string]*mongodb.Connection {
	result := make(map[string]*mongodb.Connection)
	// 读取配置文件中所有以 `mongodb` 开头的配置项
	mongodbs := viper.GetStringMap("mongodb")
	for key := range mongodbs {
		// 读取每个 MongoDB 配置项
		host := viper.GetString("mongodb." + key + ".host")
		port := viper.GetString("mongodb." + key + ".port")
		user := viper.GetString("mongodb." + key + ".user")
		password := viper.GetString("mongodb." + key + ".password")
		result[key] = mongodb.NewConnection(host, port, user, password)
	}
	return result
}

func ParseHyperchains() map[string][]*hyperchain.Node {
	result := make(map[string][]*hyperchain.Node)
	// 读取配置文件中所有以 `hyperchain` 开头的配置项
	hyperchains := viper.GetStringMap("hyperchain")
	for key, value := range hyperchains {
		for _, node := range value.([]interface{}) {
			nodeInfoMap := node.(map[string]interface{})
			result[key] = append(result[key], hyperchain.NewNode(nodeInfoMap["host"].(string), nodeInfoMap["port"].(int)))
		}
	}
	return result
}

func ParseChainrolls(mongodbs map[string]*mongodb.Connection, hyperchains map[string][]*hyperchain.Node) map[string]*chainroll.Chainroll {
	rootWorkdir := viper.GetString("workdir")
	chainrollPackage := viper.GetString("packages.chainroll")

	result := make(map[string]*chainroll.Chainroll)
	// 读取配置文件中所有以 `chainroll` 开头的配置项
	chainrollConfigMap := viper.GetStringMap("chainroll")
	for key := range chainrollConfigMap {
		chainrollGrpcPort := viper.GetString("chainroll." + key + ".grpc-port")
		chainrollHttpPort := viper.GetString("chainroll." + key + ".http-port")
		chainrollHyperchain := hyperchains[viper.GetString("chainroll."+key+".hyperchain")]
		chainrollMongo := mongodbs[viper.GetString("chainroll."+key+".mongodb")]
		chainrollMongoSuffix := viper.GetString("chainroll." + key + ".mongo-suffix")
		cr := chainroll.NewChainroll(rootWorkdir, chainrollPackage, chainrollGrpcPort, chainrollHttpPort, chainrollMongo, chainrollMongoSuffix, chainrollHyperchain)
		result[key] = cr
	}

	return result
}

//func ParseMeshes(mongodbs map[string]*mongodb.Connection, hyperchains map[string][]*hyperchain.Node) map[string]*mesh.Mesh {
//	rootWorkdir := viper.GetString("workdir")
//	chainrollPackage := viper.GetString("packages.chainroll")
//
//	result := make(map[string]*mesh.Mesh)
//	// 读取配置文件中所有以 `bitxmesh` 开头的配置项
//	meshes := viper.GetStringMap("bitxmesh")
//	for key := range meshes {
//		chainrollGrpcPort := viper.GetString("bitxmesh." + key + ".chainroll.grpc-port")
//		chainrollHttpPort := viper.GetString("bitxmesh." + key + ".chainroll.http-port")
//		chainrollMongo := mongodbs[viper.GetString("bitxmesh."+key+".chainroll.mongodb")]
//		chainrollHyperchain := hyperchains[viper.GetString("bitxmesh."+key+".chainroll.hyperchain")]
//		chainrollMongoSuffix := viper.GetString("bitxmesh." + key + ".chainroll.mongo-suffix")
//		cr := chainroll.NewChainroll(rootWorkdir, chainrollPackage, chainrollGrpcPort, chainrollHttpPort, chainrollMongo, chainrollMongoSuffix, chainrollHyperchain)
//		result[key] = mesh.New(cr)
//	}
//
//	return result
//}
