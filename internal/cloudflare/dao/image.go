package dao

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/arwoosa/media/internal/pb/image"
)

type Image struct {
	ID       string
	Filename string
	Uploaded time.Time
	Meta     map[string]string
	Variants []string
}

func (i *Image) getFloat64Ptr(key string) *float64 {
	val := i.Meta[key]
	if val == "" {
		return nil
	}
	valFloat, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return nil
	}
	return &valFloat
}

func (i *Image) getUint64(key string) uint64 {
	val, err := strconv.ParseInt(i.Meta[key], 10, 64)
	if err != nil {
		return 0
	}
	return uint64(val)
}

func (i *Image) getUint32(key string) uint32 {
	val, err := strconv.Atoi(i.Meta[key])
	if err != nil {
		return 0
	}
	return uint32(val)
}

func (i *Image) GetLatitude() *float64 {
	return i.getFloat64Ptr("latitude")
}

func (i *Image) GetLongitude() *float64 {
	return i.getFloat64Ptr("longitude")
}

func (i *Image) GetFormat() string {
	return i.Meta["format"]
}

func (i *Image) GetImageFormat() image.ImageFormat {
	return image.ImageFormat(image.ImageFormat_value[i.GetFormat()])
}

func (i *Image) GetWidth() uint32 {
	return i.getUint32("width")
}

func (i *Image) GetHeight() uint32 {
	return i.getUint32("height")
}

func (i *Image) GetSize() uint64 {
	return i.getUint64("size")
}

type ImageMetadata struct {
	Width     uint32
	Height    uint32
	Format    string
	Size      uint64
	Latitude  *float64
	Longitude *float64
}

func (i *ImageMetadata) ToCoudflareFieldMetadata() any {
	data := map[string]string{
		"width":  strconv.FormatUint(uint64(i.Width), 10),
		"height": strconv.FormatUint(uint64(i.Height), 10),
		"format": i.Format,
		"size":   strconv.FormatUint(i.Size, 10),
	}
	if i.Latitude != nil {
		data["latitude"] = strconv.FormatFloat(*i.Latitude, 'f', -1, 64)
	}
	if i.Longitude != nil {
		data["longitude"] = strconv.FormatFloat(*i.Longitude, 'f', -1, 64)
	}
	metadataJSON, _ := json.Marshal(data)
	return string(metadataJSON)
}
