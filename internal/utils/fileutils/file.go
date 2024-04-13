package fileutils

import (
	"log"
	"os"
	"os/exec"
)

func ExtractTar(tarFilePath string, targetDir string) {
	// TODO 不依赖系统的 tar 命令，使用 Go 语言的 tar 包解压 tar 文件
	cmd := exec.Command("tar", "-zxvf", tarFilePath, "-C", targetDir)
	//cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalln("Error extracting tar file: ", err)
	}
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
