package storagedriver

import (
	"crypto/sha256"
	"encoding/base64"
	"image"
)

type StorageDriver interface {
	// Returns the name of the driver
	Name() string

	// Stores an image
	Store(b []byte) error

	// Fetches an image matching a given RGBA value
	// TODO: enable a user to Get images excluding a list of image
	Get(r, g, b uint8) (image.Image, error)
}

var RawURLEncoding = base64.URLEncoding.WithPadding(base64.NoPadding)

func Encode(b []byte) string {
	h := sha256.New()
	h.Write(b)
	sum := h.Sum(nil)

	return RawURLEncoding.EncodeToString(sum)
}

// TODO: move this into a util package
func GetAverageColor(m image.Image, bounds image.Rectangle) (ra, ga, ba int) {
	// TODO: work with pixel array directly for efficiency
	//bounds := m.Bounds()
	top, bottom := bounds.Min.Y, bounds.Max.Y
	left, right := bounds.Min.X, bounds.Max.X

	rt, gt, bt := 0, 0, 0
	num_pixels := (bottom - top) * (right - left)
	for y := top; y < bottom; y++ {
		for x := left; x < right; x++ {
			r, g, b, _ := m.At(x, y).RGBA()
			rt += int(r)
			gt += int(g)
			bt += int(b)
		}
	}
	rt, gt, bt = rt/257, gt/257, bt/257
	ra, ga, ba = rt/num_pixels, gt/num_pixels, bt/num_pixels
	return ra, ga, ba
}

func bucket(bucketsize, num int) int {
	return num - (num % bucketsize)
}

func ColorBucket(bucketsize, r, g, b int) (int, int, int) {
	return bucket(bucketsize, r), bucket(bucketsize, g), bucket(bucketsize, b)
}
