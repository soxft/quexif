package read

import (
	"log"
	"os"
	"path/filepath"
	"photo_exif_do/fg"
	"photo_exif_do/tool"
	"photo_exif_do/x_exif"
	"strings"
)

func isPathHasPrefix(path string, prefix string) bool {
	x := strings.Split(path, "/")
	for _, v := range x {
		if strings.HasPrefix(v, prefix) {
			return true
		}
	}

	return false
}

func Run(path string) {
	// check is file
	if tool.IsFile(path) {
		filename := filepath.Base(path)
		t, e := x_exif.ReadExif(path)
		if e != nil {
			log.Printf("\033[0;31m[%s]\033[0m -> %v\n", filename, e.Error())
			return
		}

		log.Printf("\033[0;32m[%s]\033[0m -> %v\n", filename, t)
		return
	}

	// 递归遍历目录
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		filename := strings.Replace(path, fg.Path, "", -1)

		if err != nil {
			log.Printf("\033[0;31m[%s]\033[0m -> %v\n", filename, err)
			return nil
		}

		if info.IsDir() || !x_exif.IsExtValid(path) || isPathHasPrefix(filename, ".") {
			// log.Printf("[SKIP](%d) %s -> %v\n", counter, path, "is a directory")
			return nil
		}

		t, e := x_exif.ReadExif(path)
		if e != nil {
			log.Printf("\033[0;31m[%s]\033[0m -> %v\n", filename, e.Error())
			return nil
		}

		log.Printf("\033[0;32m[%s]\033[0m -> %v\n", filename, t)

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
