package qumagie

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"photo_exif_do/x_exif"
	"regexp"
	"time"
)

// getTime 从文件路径中解析出文件名，并将其转换为时间对象
// 时间格式为 "2006-01-02 15.04.05"
// NOTICE 仅适用于 QuMagie 备份的文件名格式
func getTime(filePath string) (time.Time, error) {
	// 从路径中提取文件名
	fileName := filepath.Base(filePath)

	pattern := `\d{4}-\d{2}-\d{2} \d{2}\.\d{2}\.\d{2}`
	re := regexp.MustCompile(pattern)
	timeStr := re.FindString(fileName)

	// 将文件名解析为时间
	layout := "2006-01-02 15.04.05"
	parsedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("解析时间失败: %v", err)
	}

	return parsedTime, nil
}

// Run 处理 QuMagie 备份的照片
func Run(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	counter := 1

	// 遍历目录下的文件
	for _, file := range files {
		// QuMagie 不会有子文件夹, 所以不需要递归
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(path, file.Name())
		// 从路径中提取文件名
		filename := filepath.Base(filePath)

		t, err := getTime(filePath)
		if err != nil {
			log.Printf("\033[0;33m[SKIP]\033[0m(%d) %s -> %v\n", counter, filename, err)
			counter++
			continue
		}

		if err := x_exif.SetDate(filePath, t); err == nil {
			log.Printf("\033[0;32m[SUCC]\033[0m(%d) %s -> %s\n", counter, filename, t.Format("2006-01-02 15.04.05"))
			x_exif.RemoveEditStr(filePath)
		} else if errors.Is(err, x_exif.ErrAlreadyHasDate) || errors.Is(err, x_exif.ErrMediaTypeNotSupport) {
			log.Printf("\033[0;33m[SKIP]\033[0m(%d) %s -> %v\n", counter, filename, err)
		} else {
			log.Printf("\033[0;31m[FAIL]\033[0m(%d) %s -> %v\n", counter, filename, err)
		}

		counter++
	}
}
