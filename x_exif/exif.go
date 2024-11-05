package x_exif

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	jpeg "github.com/dsoprea/go-jpeg-image-structure/v2"
	"log"
	"os"
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

// SetDate 为文件设置日期 如果已经存在则跳过
func SetDate(filePath string, t time.Time, skip bool) error {
	if !strings.HasSuffix(strings.ToLower(filePath), ".jpg") && !strings.HasSuffix(strings.ToLower(filePath), ".jpeg") {
		return ErrMediaTypeNotSupport
	}

	// 检测是否已经有日期
	if skip {
		aDate, err := ReadExif(filePath)
		if aDate != "" || err != nil {
			return ErrAlreadyHasDate
		}
	}

	parser := jpeg.NewJpegMediaParser()
	intfc, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse JPEG file: %v", err)
	}

	sl := intfc.(*jpeg.SegmentList)

	rootIb, err := sl.ConstructExifBuilder()
	if err != nil {
		im, err := exifcommon.NewIfdMappingWithStandard()
		if err != nil {
			return fmt.Errorf("failed to create new IFD mapping with standard tags: %v", err)
		}
		ti := exif.NewTagIndex()
		if err := exif.LoadStandardTags(ti); err != nil {
			return fmt.Errorf("failed to load standard tags: %v", err)
		}

		rootIb = exif.NewIfdBuilder(im, ti, exifcommon.IfdStandardIfdIdentity,
			exifcommon.EncodeDefaultByteOrder)
	}

	ts := exifcommon.ExifFullTimestampString(t)

	// Set DateTime
	if err := setExifTag(rootIb, "IFD0", "DateTime", ts); err != nil {
		return fmt.Errorf("failed to set tag %v: %v", "DateTime", err)
	}

	// Set DateTimeOriginal
	if err := setExifTag(rootIb, "IFD/Exif", "DateTimeOriginal", ts); err != nil {
		return fmt.Errorf("failed to set tag %v: %v", "DateTimeOriginal", err)
	}

	// Update the exif segment.
	if err := sl.SetExif(rootIb); err != nil {
		return fmt.Errorf("failed to set EXIF to jpeg: %v", err)
	}

	// Write the modified file
	b := new(bytes.Buffer)
	if err := sl.Write(b); err != nil {
		return fmt.Errorf("failed to create JPEG data: %v", err)
	}

	// Save the file
	if err := os.WriteFile(filePath, b.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write JPEG file: %v", err)
	}

	return nil
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
