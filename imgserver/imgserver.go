package imgserver

import (
	"bytes"
	"fmt"
	"image"
	//"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/alvinfeng/mosaic/config"
	"github.com/alvinfeng/mosaic/storage/cache/inmemory"
	"github.com/alvinfeng/mosaic/storage/driver"
	"github.com/alvinfeng/mosaic/storage/driver/filesystem"
	"github.com/alvinfeng/mosaic/storage/driver/simple"
)

type ImgServer struct {
	Router *mux.Router
	store  storagedriver.StorageDriver
}

func New(c *config.Config) (*ImgServer, error) {
	var store storagedriver.StorageDriver
	var err error
	switch storageType := c.StorageType; storageType {
	case "filesystem":
		fmt.Println("Using filesystem")
		store, err = filesystem.New(c.Fs)
	case "simple":
		fmt.Println("Using simple")
		store, err = simple.New()
	default:
		return nil, fmt.Errorf("no storage driver specified")
	}

	switch cacheType := c.CacheType; cacheType {
	case "inmemory":
		fmt.Println("Using inmemory cache")
		store.SetCache(inmemory.New())
	default:
		fmt.Println("No cache specified")
	}

	s := &ImgServer{
		store:  store,
		Router: mux.NewRouter(),
	}

	s.Router.HandleFunc("/image", s.FetchImage)
	s.Router.HandleFunc("/mosaic", s.CreateMosaic)
	s.Router.HandleFunc("/store", s.StoreImage)
	fmt.Println("Image Server created")
	return s, err
}

func (s *ImgServer) CreateMosaic(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	size := ParseTilesize(req)
	scale := 150 / size // default sub image size is 150x150
	start := time.Now()
	fmt.Printf("Creating mosaic with size %v and scale %v\n", size, scale)

	m, _, err := image.Decode(req.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer req.Body.Close()

	bounds := m.Bounds()

	scaledOutput := image.NewRGBA(image.Rect(0, 0, bounds.Max.X*scale, bounds.Max.Y*scale))

	for x := 0; x < bounds.Max.X; x += size {
		for y := 0; y < bounds.Max.Y; y += size {
			maxX := x + size
			if maxX > bounds.Max.X {
				maxX = bounds.Max.X
			}
			maxY := y + size
			if maxY > bounds.Max.Y {
				maxY = bounds.Max.Y
			}
			rect := image.Rect(x, y, maxX, maxY)
			r, g, b := GetAverageColor(m, rect)
			// fmt.Printf("Fetching sub image for (%v, %v, %v, %v) with rgb: (%v, %v, %v)\n", x, y, maxX, maxY, r, g, b)
			subimg, err := s.store.Get(uint8(r), uint8(g), uint8(b))
			if err != nil {
				fmt.Println(err)
			}
			scaledRect := image.Rect(x*scale, y*scale, maxX*scale, maxY*scale)
			//draw.Draw(output, rect, subimg, image.ZP, draw.Src)
			draw.Draw(scaledOutput, scaledRect, subimg, subimg.Bounds().Min, draw.Src)
			//draw.Draw(output, rect, &image.Uniform{color.RGBA{uint8(r), uint8(g), uint8(b), 255}}, image.ZP, draw.Src)
		}
	}

	buff := new(bytes.Buffer)

	var opt jpeg.Options
	opt.Quality = 80
	err = jpeg.Encode(buff, scaledOutput, &opt)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buff.Bytes())))
	if _, err := w.Write(buff.Bytes()); err != nil {
		fmt.Println(err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Started at %v - elapsed %v\n", start, elapsed)
}

func (s *ImgServer) FetchImage(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Fetching image")
	err := req.ParseForm()
	if err != nil {
		fmt.Println(err)
	}

	r, g, b, err := ParseRGB(req)
	if err != nil {
		fmt.Println(err)
	}

	img, err := s.store.Get(r, g, b)
	if err != nil {
		fmt.Println(err)
	}

	buff := new(bytes.Buffer)

	var opt jpeg.Options
	opt.Quality = 80
	err = jpeg.Encode(buff, img, &opt)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buff.Bytes())))
	if _, err := w.Write(buff.Bytes()); err != nil {
		fmt.Println(err)
	}
}

func (s *ImgServer) StoreImage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Storing image")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer r.Body.Close()

	err = s.store.Store(b)
	if err != nil {
		fmt.Println(err)
	}
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

func ParseRGB(r *http.Request) (uint8, uint8, uint8, error) {
	rgb := r.Form.Get("rgb")
	vals := strings.Split(rgb, ",")
	if len(vals) != 3 {
		return 0, 0, 0, fmt.Errorf("failed parsing rgb")
	}

	var nums [3]uint8
	for i := 0; i < 3; i++ {
		num, err := strconv.ParseUint(vals[i], 10, 8)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("could not part %v to uint8: %v", vals[i], err)
		}
		nums[i] = uint8(num)
	}
	return nums[0], nums[1], nums[2], nil
}

func ParseTilesize(r *http.Request) int {
	defaultSize := 30
	size := r.Form.Get("tilesize")
	fmt.Println("got size: ", size)
	num, err := strconv.Atoi(size)
	if err != nil || num <= 0 || num > 150 {
		return defaultSize
	}
	return num
}
