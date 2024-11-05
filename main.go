package main

import (
	"bytes"
	"fmt"
	"github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	jpeg "github.com/dsoprea/go-jpeg-image-structure/v2"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// getFileName 从文件路径中解析出文件名，并将其转换为时间对象
// 时间格式为 "2006-01-02 15.04.05"
func getFileName(filePath string) (time.Time, error) {
	// 从路径中提取文件名
	fileName := filepath.Base(filePath)

	fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	fileNameWithoutExt = fileNameWithoutExt[:19]

	// 将文件名解析为时间
	layout := "2006-01-02 15.04.05"
	parsedTime, err := time.Parse(layout, fileNameWithoutExt)
	if err != nil {
		return time.Time{}, fmt.Errorf("解析时间失败: %v", err)
	}

	return parsedTime, nil
}

// ReadExif 获取exif 中的 DateTimeOriginal
func ReadExif(path string) string {
	opt := exif.ScanOptions{}
	dt, err := exif.SearchFileAndExtractExif(path)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	ets, _, err := exif.GetFlatExifData(dt, &opt)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	for _, et := range ets {
		if et.TagName == "DateTimeOriginal" {
			return et.Value.(string)
		}
	}

	return ""
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

// setDateIfNone 为文件设置日期 如果已经存在则跳过
func setDateIfNone(filePath string, counter int) error {
	parser := jpeg.NewJpegMediaParser()
	intfc, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse JPEG file: %v", err)
	}

	sl := intfc.(*jpeg.SegmentList)

	rootIb, err := sl.ConstructExifBuilder()
	if err != nil {
		//fmt.Println("No EXIF; creating it from scratch")

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
		err = rootIb.AddStandardWithName("ProcessingSoftware", "photos-uploader")
		if err != nil {
			return err
		}
	}

	// 检测是否已经有日期
	aDate := ReadExif(filePath)
	if aDate != "" {
		log.Printf("[SKIP](%d) %s already has date: %s\n", counter, filepath.Base(filePath), aDate)
		return nil
	}

	// Form our timestamp string
	// 从名称解析日期
	t, err := getFileName(filePath)
	if err != nil {
		return fmt.Errorf("从文件名解析日期失败: %v", err)
	}
	ts := exifcommon.ExifFullTimestampString(t)

	// Set DateTime
	ifdPath := "IFD0"
	if err := setExifTag(rootIb, ifdPath, "DateTime", ts); err != nil {
		return fmt.Errorf("failed to set tag %v: %v", "DateTime", err)
	}

	// Set DateTimeOriginal
	ifdPath = "IFD/Exif"
	if err := setExifTag(rootIb, ifdPath, "DateTimeOriginal", ts); err != nil {
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

	log.Printf("[SUCC](%d) %s -> %s\n", counter, filepath.Base(filePath), ts)

	return nil
}

// removeEditStr 如果文件名中包含 (edited) 则自动重命名该文件
// 主要用于处理 QFiling 归档文件
func removeEditStr(path string) {
	if strings.Contains(path, "(edited)") {
		newPath := strings.ReplaceAll(path, "(edited)", "")
		if err := os.Rename(path, newPath); err != nil {
			log.Printf("重命名失败: %v\n", err)
		} else {
			log.Printf("重命名成功: %s -> %s\n", path, newPath)
		}
	}
}

func main() {
	// 读取命令行参数 - 目录
	if len(os.Args) < 2 {
		log.Fatal("Usage: exif_pass <directory>")
	}

	// 安全 QA
	log.Println("请确保已经设置了快照，程序将会直接修改文件的 exif 元数据, 建议在使用前选择少量照片进行测试后再使用")
	log.Println("您确认要继续吗? (y/n)")
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" {
		log.Fatal("已取消")
	}

	// 读取目录
	dir := os.Args[1]
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	counter := 1

	// 遍历目录下的文件
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := filepath.Join(dir, file.Name())
		if strings.HasSuffix(strings.ToLower(filePath), ".jpg") || strings.HasSuffix(strings.ToLower(filePath), ".jpeg") {
			if err := setDateIfNone(filePath, counter); err != nil {
				log.Printf("[FAIL](%d) %s: %v\n", counter, filePath, err)
			} else {
				removeEditStr(filePath)
			}
		} else {
			log.Printf("[SKIP](%d) %s\n", counter, filePath)
		}

		counter++
	}
}
