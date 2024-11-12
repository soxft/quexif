package dir_auto

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"photo_exif_do/fg"
	"photo_exif_do/x_exif"
	"regexp"
	"strings"
	"time"
)

// Run 处理指定文件夹下的照片, 从后向前尝试匹配时间
func Run(path string) {
	// 递归遍历目录
	counter := 1
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("\033[0;33m[SKIP]\033[0m(%d) %s -> %v\n", counter, path, err)
			return nil
		}

		if info.IsDir() {
			// log.Printf("[SKIP](%d) %s -> %v\n", counter, path, "is a directory")
			return nil
		}

		pathLinux := parsePath(path)

		if !x_exif.IsExtValid(pathLinux) {
			log.Printf("\033[0;33m[SKIP]\033[0m(%d) %s -> %v\n", counter, path, x_exif.ErrMediaTypeNotSupport)
			return nil
		}

		// 时间解析失败
		t, err := tryGetDate(pathLinux)
		if err != nil {
			log.Printf("\033[0;33m[SKIP]\033[0m(%d) %s -> %v\n", counter, path, err)
			counter++
			return nil
		}

		// log.Printf("[INFO](%d) %s -> %v\n", counter, path, t.Format("2006-01-02 15.04.05"))

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

func parsePath(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

// tryGetDate 尝试从路径中获取日期
func tryGetDate(path string) (time.Time, error) {
	// 从路径中提取文件名
	sp := strings.Split(path, "/")

	for i := len(sp) - 1; i >= 0; i-- {
		if fg.Regex != "" {
			// 尝试进行正则匹配
			re := regexp.MustCompile(fg.Regex)

			if re.MatchString(sp[i]) {
				timeStr := re.FindStringSubmatch(sp[i])

				for _, v := range timeStr {
					parsedTime, err := time.Parse(fg.DateTpl, v)
					if err == nil {
						return parsedTime, nil
					}
				}

				continue
			}

			continue
		}

		// 对于后缀
		if i == len(sp)-1 && strings.Contains(sp[i], ".") {
			sp := strings.Split(sp[i], ".")
			parsedTime, err := time.Parse(fg.DateTpl, strings.Join(sp[:len(sp)-1], "."))
			if err == nil {
				return parsedTime, nil
			}

			continue
		}

		// 从路径中提取文件名
		parsedTime, err := time.Parse(fg.DateTpl, sp[i])
		if err == nil {
			return parsedTime, nil
		}

	}

	return time.Time{}, errors.New("解析时间失败")
}
