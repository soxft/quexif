package main

import (
	"fmt"
	"log"
	"photo_exif_do/dir"
	"photo_exif_do/dir_date"
	"photo_exif_do/fg"
	"photo_exif_do/qumagie"
	"photo_exif_do/x_exif"
	"strings"
)

func main() {
	fg.Parse()

	// 安全 QA
	if !fg.SkipSafeQA {
		log.Println("请确保已经设置了快照，程序将会直接修改文件的 exif 元数据, 是否继续? (y/N)")
		var confirm string

		if _, err := fmt.Scanln(&confirm); err != nil {
			return
		} else if strings.ToLower(confirm) != "y" {
			log.Fatal("已取消")
		}
	}

	// 读取目录
	switch fg.Mode {
	case "dir": // 指定文件夹批量修改
		dir.Run(fg.Path)
		break
	case "dir_date": // 按照上级文件夹名称修改
		dir_date.Run(fg.Path)
		break
	case "test":
		// log.Println(x_exif.SetDate("pics/qumagie/2006-01-02 15.04.05.jpg", time.Now(), true))
		log.Println(x_exif.ReadExif("pics/qumagie/2022-06-15 10.13.50.png"))
	default:
		qumagie.Run(fg.Path)
	}

}
