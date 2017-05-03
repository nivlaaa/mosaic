package imgserver

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ImgServerTestSuite struct {
	suite.Suite
}

func (suite *ImgServerTestSuite) TestParseRGB() {
	fmt.Println("Testing ParseRGB")
	req, err := http.NewRequest("GET", "http://example.com", nil)
	suite.Nil(err)

	req.Form = url.Values{}
	req.Form.Add("rgb", "127,127,127")
	r, g, b, err := ParseRGB(req)
	suite.Nil(err)
	suite.Equal(r, uint8(127))
	suite.Equal(g, uint8(127))
	suite.Equal(b, uint8(127))

	req.Form = url.Values{}
	req.Form.Add("rgb", "999,999,999")
	r, g, b, err = ParseRGB(req)
	suite.NotNil(err)

	req.Form = url.Values{}
	req.Form.Add("rgb", "abcdefg")
	r, g, b, err = ParseRGB(req)
	suite.NotNil(err)

	req.Form = url.Values{}
	r, g, b, err = ParseRGB(req)
	suite.NotNil(err)
}

func TestImgServerTestSuite(t *testing.T) {
	suite.Run(t, new(ImgServerTestSuite))
}
