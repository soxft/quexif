package x_exif

import (
	"errors"
	"fmt"
	"github.com/dsoprea/go-exif/v3"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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
func SetDate(filePath string, t time.Time, skip bool) error {
	// 检测是否已经有日期
	if skip {
		aDate, err := ReadExif(filePath)
		if aDate != "" || err != nil {
			return ErrAlreadyHasDate
		}
	}

	switch strings.ToLower(filepath.Ext(filePath)) {
	case ".jpg", ".jpeg":
		return setJpgExif(filePath, t)
	case ".png":
		return setPngDate(filePath, t)
	default:
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

var ErrAlreadyHasDate = fmt.Errorf("already has date")
var ErrMediaTypeNotSupport = fmt.Errorf("media type not support")