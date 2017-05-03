package filesystem

import (
	"image"
)

const (
	driverName = "filesystem"
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
