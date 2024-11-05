package x_exif

import (
	"bytes"
	"fmt"
	"github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	log "github.com/dsoprea/go-logging"
	png "github.com/dsoprea/go-png-image-structure/v2"
	riimage "github.com/dsoprea/go-utility/v2/image"
	"os"
	"time"
)

func setPngDate(filePath string, t time.Time) error {
	praser := png.NewPngMediaParser()
	intfc, err := praser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse PNG file: %v", err)
	}

	rootIb, err := pngConstructExifBuilder(intfc)
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

	cs := intfc.(*png.ChunkSlice)

	err = cs.SetExif(rootIb)
	log.PanicIf(err)

	b := new(bytes.Buffer)

	// Write to a `bytes.Buffer`.
	err = cs.WriteTo(b)

	// Save the file
	if err := os.WriteFile(filePath, b.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write JPEG file: %v", err)
	}

	return nil
}

func pngConstructExifBuilder(intfc riimage.MediaContext) (rootIb *exif.IfdBuilder, err error) {

	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	rootIfd, _, err := intfc.Exif()
	if log.Is(err, exif.ErrNoExif) == true {
		// No EXIF. Just create a boilerplate builder.

		im := exifcommon.NewIfdMapping()

		err := exifcommon.LoadStandardIfds(im)
		log.PanicIf(err)

		ti := exif.NewTagIndex()

		rootIb :=
			exif.NewIfdBuilder(
				im,
				ti,
				exifcommon.IfdStandardIfdIdentity,
				exifcommon.EncodeDefaultByteOrder)

		return rootIb, nil
	} else if err != nil {
		log.Panic(err)
	}

	rootIb = exif.NewIfdBuilderFromExistingChain(rootIfd)

	return rootIb, nil
}
