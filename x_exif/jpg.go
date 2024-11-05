package x_exif

import (
	"bytes"
	"fmt"
	"github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	jpeg "github.com/dsoprea/go-jpeg-image-structure/v2"
	"os"
	"time"
)

func setJpgExif(filePath string, t time.Time) error {

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
