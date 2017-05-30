package simple

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/alvinfeng/mosaic/storage/cache"
	"github.com/alvinfeng/mosaic/storage/driver"
)

const (
	driverName = "simple"
	baseDir    = "./buckets"
)

type driver struct {
	buckets cache.Cache
}

func New() (*driver, error) {
	return &driver{}, nil
}

func (d *driver) Name() string {
	return driverName
}

func (d *driver) SetCache(c cache.Cache) {
	d.buckets = c
}

func (d *driver) Store(data []byte) error {
	uuid := storagedriver.Encode(data)
	m, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return err
	}

	r, g, b := storagedriver.GetAverageColor(m, m.Bounds())
	rbucket, gbucket, bbucket := storagedriver.ColorBucket(15, r, g, b)
	fmt.Printf("uuid: %v, format: %v, rgb: (%v, %v, %v), bucket: (%v, %v, %v)\n", uuid, format, r, g, b, rbucket, gbucket, bbucket)

	bucketPath := filepath.Join(baseDir, fmt.Sprintf("r%v_g%v_b%v", rbucket, gbucket, bbucket))
	err = os.MkdirAll(bucketPath, os.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(bucketPath, fmt.Sprintf("%v.%v", uuid, format)), data, 0644)
	return err
}

func (d *driver) Get(r, g, b uint8) (image.Image, error) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	c := color.RGBA{r, g, b, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{c}, image.ZP, draw.Src)
	return img, nil
}
