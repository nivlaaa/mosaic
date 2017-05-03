package s3

import (
	"image"
)

const (
	driverName = "s3"
)

type Driver struct{}

func (d *Driver) Name() string {
	return driverName
}

func (d *Driver) Store(m image.Image) error {
	return nil
}

func (d *Driver) Get(r, g, b, a int) (image.Image, error) {
	return nil, nil
}
