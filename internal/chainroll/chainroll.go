package chainroll

import (
	"fmt"
	"git.hyperchain.cn/dmlab/bitxmesh-ops/internal/hyperchain"
	"git.hyperchain.cn/dmlab/bitxmesh-ops/internal/mongodb"
	"git.hyperchain.cn/dmlab/bitxmesh-ops/internal/utils/fileutils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const envFilename = "chainroll-auto.env"

type Chainroll struct {
	workdir           string
	packagePath       string
	grpcPort          string
	httpPort          string
	mongodbConnection *mongodb.Connection
	mongoSuffix       string
	hyperchain        []*hyperchain.Node
}

// NewChainroll 创建 Chainroll 实例
func NewChainroll(rootWorkdir string, packagePath string, grpcPort string, httpPort string, mongodbConnection *mongodb.Connection, mongoSuffix string, hyperchain []*hyperchain.Node) *Chainroll {
	workdir := filepath.Join(rootWorkdir, fmt.Sprintf(".%s-grpc%s-http%s", "chainroll", grpcPort, httpPort))
	return &Chainroll{
		workdir,
		packagePath,
		grpcPort,
		httpPort,
		mongodbConnection,
		mongoSuffix,
		hyperchain,
	}
}

// Start 部署 chainroll
func (c *Chainroll) Start() {
	// TODO 检查 hyperchain 和 mongodb 状态（并创建 mongodb 相关数据库）
	c.makeWorkspace()
	//c.generateEnvFile()
	c.replaceHyperchain()
	c.execStartScript()
	isRunnging := c.status()
	if isRunnging {
		log.Printf("🎉🎉🎉 Chainroll started at %s\n", c.workdir)
	} else {
		log.Fatalln("❌ Error starting chainroll")
	}

}

// makeWorkspace 创建工作目录
func (c *Chainroll) makeWorkspace() {
	if fileutils.FileExists(c.workdir) {
		log.Printf("Workspace %s already exists, continue", c.workdir)
		return
	} else {
		err := os.MkdirAll(c.workdir, os.ModePerm)
		if err != nil {
			log.Fatalln("Error creating workspace: ", err)
		}
	}
	if !fileutils.FileExists(c.packagePath) {
		log.Fatalf("Package %s does not exist\n", c.packagePath)
	}
	fileutils.ExtractTar(c.packagePath, c.workdir)
}

// TODO 可以删去该步骤
// generateEnvFile 生成环境变量文件，返回生成的文件路径
func (c *Chainroll) generateEnvFile() string {
	if c.mongodbConnection == nil {
		log.Fatalln("MongoDB connection is nil")
	}
	// 写入文件
	env := map[string]string{
		"CHAINROLL_GRPC_PORT":          c.grpcPort,
		"CHAINROLL_HTTP_PORT":          c.httpPort,
		"CHAINROLL_DB_MASTER_IP":       c.mongodbConnection.Host(),
		"CHAINROLL_DB_MASTER_PORT":     c.mongodbConnection.Port(),
		"CHAINROLL_DB_MASTER_USERNAME": c.mongodbConnection.User(),
		"CHAINROLL_DB_MASTER_PASSWORD": c.mongodbConnection.Password(),
		"CHAINROLL_DB_MASTER_SUFFIX":   c.mongoSuffix,
	}

	envFilePath := fmt.Sprintf("%s/%s", c.workdir, envFilename)
	file, err := os.Create(envFilePath)
	if err != nil {
		log.Fatalln("Error creating env file: ", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalln("Error closing env file: ", err)
		}
	}(file)

	for key, value := range env {
		_, err := file.WriteString(fmt.Sprintf("export %s=%s\n", key, value))
		if err != nil {
			log.Fatalln("Error writing to env file: ", err)
		}
	}
	err = file.Sync()
	if err != nil {
		log.Fatalln("Error syncing env file: ", err)
	}
	return envFilePath
}

// replaceHyperchain 替换 hyperchain
func (c *Chainroll) replaceHyperchain() {
	hpcTomlFile := filepath.Join(c.workdir, "config", "blockchain", "hyperchain", "hpc.toml")
	content, err := os.ReadFile(hpcTomlFile)
	if err != nil {
		log.Fatalln("Error reading hpc.toml: ", err)
	}

	nodesString := func(nodes []string) string {
		// 构造节点字符串，格式为 `["node1", "node2", ...]`
		nodesStr := `["` + nodes[0] + `"`
		for _, node := range nodes[1:] {
			nodesStr += `, "` + node + `"`
		}
		nodesStr += `]`
		return nodesStr
	}

	// [jsonRPC]
	//	 nodes = ["172.22.67.127", "172.22.67.127", "172.22.67.127", "172.22.67.127"]
	//	 ports = ["8081", "8082", "8083", "8084"]

	// 使用正则表达式替换节点和端口
	nodeRegex := regexp.MustCompile(`nodes\s*=\s*\[.*]`)
	portRegex := regexp.MustCompile(`ports\s*=\s*\[.*]`)

	nodes := make([]string, 0)
	ports := make([]string, 0)
	for _, node := range c.hyperchain {
		nodes = append(nodes, node.Host())
		ports = append(ports, fmt.Sprintf("%d", node.Port()))
	}
	// 替换节点
	content = nodeRegex.ReplaceAll(content, []byte(fmt.Sprintf("nodes = %s", nodesString(nodes))))
	// 替换端口
	content = portRegex.ReplaceAll(content, []byte(fmt.Sprintf("ports = %s", nodesString(ports))))

	// 将替换后的内容写回文件
	err = os.WriteFile(hpcTomlFile, content, 0644)
	if err != nil {
		log.Fatalln("Error writing to hpc.toml: ", err)
	}
}

// execStartScript 执行启动脚本
func (c *Chainroll) execStartScript() {
	//cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cd %s && source ./%s && ./start.sh", c.workdir, envFilename))
	//cmd := exec.Command(fmt.Sprintf("source ./%s && ./start.sh", envFilename))
	envMap := map[string]string{
		"CHAINROLL_GRPC_PORT":          c.grpcPort,
		"CHAINROLL_HTTP_PORT":          c.httpPort,
		"CHAINROLL_DB_MASTER_IP":       c.mongodbConnection.Host(),
		"CHAINROLL_DB_MASTER_PORT":     c.mongodbConnection.Port(),
		"CHAINROLL_DB_MASTER_USERNAME": c.mongodbConnection.User(),
		"CHAINROLL_DB_MASTER_PASSWORD": c.mongodbConnection.Password(),
		"CHAINROLL_DB_MASTER_SUFFIX":   c.mongoSuffix,
	}
	for key, value := range envMap {
		err := os.Setenv(key, value)
		if err != nil {
			log.Fatalln("Error setting env: ", err)
		}
	}
	cmd := exec.Command("bash", "./start.sh")
	//cmd := exec.Command("source", fmt.Sprintf("./%s", envFilename))
	cmd.Dir = c.workdir
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalln("Error starting chainroll: ", err)
	}
	fmt.Printf("%s", out)
}

func (c *Chainroll) status() bool {
	count := 1
	for count < 10 {
		cmd := exec.Command("bash", "./status.sh")
		cmd.Dir = c.workdir
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalln("Error getting chainroll status: ", err)
		}
		if strings.HasPrefix(string(out), "chainroll is running") {
			return true
		}
		count++
		time.Sleep(1 * time.Second)
	}
	return false
}
