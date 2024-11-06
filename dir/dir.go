package dir

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"photo_exif_do/fg"
	"photo_exif_do/x_exif"
	"strings"
	"time"
)

// Run 处理指定文件夹下的照片
func Run(path string) {
	// 尝试解析时间
	t, err := time.Parse(fg.DateTpl, fg.DateTime)
	if err != nil {
		log.Fatalf("解析时间失败: %v", err)
	}

	if !fg.SkipSafeQA {
		log.Printf("即将对 %s 下的所有文件设置时间为 %s, 继续? (y/N)\n", path, t.Format(fg.DateTpl))
		var confirm string

		if _, err := fmt.Scanln(&confirm); err != nil {
			return
		} else if strings.ToLower(confirm) != "y" {
			log.Fatal("已取消")
		}
	}

	// 递归遍历目录
	counter := 1
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("\033[0;33m[SKIP]\033[0m(%d) %s -> %v\n", counter, path, err)
			return nil
		}

		if info.IsDir() {
			// log.Printf("[SKIP](%d) %s -> %v\n", counter, path, "is a directory")
			return nil
		}

		if err := x_exif.SetDate(path, t); err == nil {
			log.Printf("\033[0;32m[SUCC]\033[0m(%d) %s -> %s\n", counter, path, t.Format("2006-01-02 15.04.05"))
		} else if errors.Is(err, x_exif.ErrAlreadyHasDate) || errors.Is(err, x_exif.ErrMediaTypeNotSupport) {
			log.Printf("\033[0;33m[SKIP]\033[0m(%d) %s -> %v\n", counter, path, err)
		} else {
			log.Printf("\033[0;31m[FAIL]\033[0m(%d) %s -> %v\n", counter, path, err)
		}
		counter++

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
