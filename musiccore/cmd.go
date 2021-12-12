package musiccore

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

//执行命令行程序，知道结束后返回
func ExecCommand(commandName string, params []string) (string, bool) {
	cmd := exec.Command(commandName, params...)

	//显示运行的命令
	fmt.Println(cmd.Args)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return "", false
	}

	cmd.Start()

	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	var linesData string
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
		linesData += line
	}

	if err = cmd.Wait(); err != nil {
		log.Printf("exec this command:%v,err:%v", cmd.Args, err)
	}

	return linesData, true
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}
