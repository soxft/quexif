package tool

import "os"

func IsFile(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false // 文件不存在
	}
	return !info.IsDir() // 如果不是目录，则是文件
}
