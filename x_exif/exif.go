package x_exif

import (
	"errors"
	"fmt"
	"github.com/dsoprea/go-exif/v3"
	"log"
	"os"
	"path/filepath"
	"photo_exif_do/fg"
	"strings"
	"time"
)

// setExifTag 设置exif标签
func setExifTag(rootIB *exif.IfdBuilder, ifdPath, tagName, tagValue string) error {
	// fmt.Printf("setTag(): ifdPath: %v, tagName: %v, tagValue: %v",
	//	ifdPath, tagName, tagValue)

	ifdIb, err := exif.GetOrCreateIbFromRootIb(rootIB, ifdPath)
	if err != nil {
		return fmt.Errorf("failed to get or create IB: %v", err)
	}

	if err := ifdIb.SetStandardWithName(tagName, tagValue); err != nil {
		return fmt.Errorf("failed to set DateTime tag: %v", err)
	}

	return nil
}

// ReadExif 获取exif 中的 DateTimeOriginal
func ReadExif(path string) (string, error) {
	opt := exif.ScanOptions{}
	dt, err := exif.SearchFileAndExtractExif(path)
	if err != nil {
		return "", fmt.Errorf("failed with Err: %v", err)
	}
	ets, _, err := exif.GetFlatExifData(dt, &opt)
	if err != nil {
		return "", fmt.Errorf("failed with Err: %v", err)
	}
	for _, et := range ets {
		if et.TagName == "DateTimeOriginal" {
			return et.Value.(string), nil
		}
	}

	return "", errors.New("mo DateTime from exif")
}

// SetDate 为文件设置日期 如果已经存在则跳过
func SetDate(filePath string, t time.Time) error {
	// 前置后缀检查
	if !IsExtValid(filePath) {
		return ErrMediaTypeNotSupport
	}

	// 检测是否已经有日期
	if !fg.Force {
		aDate, _ := ReadExif(filePath)
		if aDate != "" {
			return ErrAlreadyHasDate
		}
	}

	switch strings.ToLower(filepath.Ext(filePath)) {
	case ".jpg", ".jpeg":
		return setJpgExif(filePath, t)
	case ".png":
		return setPngExif(filePath, t)
	}

	return ErrMediaTypeNotSupport
}

// RemoveEditStr 如果文件名中包含 (edited) 则自动重命名该文件
// 主要用于处理 QFiling 归档文件
func RemoveEditStr(path string) {
	if strings.Contains(path, "(edited)") {
		newPath := strings.ReplaceAll(path, "(edited)", "")
		if err := os.Rename(path, newPath); err != nil {
			log.Printf("重命名失败: %v\n", err)
		} else {
			log.Printf("重命名成功: %s -> %s\n", path, newPath)
		}
	}
}

func IsExtValid(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg", ".png":
		return true
	}
	return false
}

var ErrAlreadyHasDate = fmt.Errorf("已经存在 Exif 日期")
var ErrMediaTypeNotSupport = fmt.Errorf("不支持的媒体类型, 仅支持 jpg 和 png")
