package service

import (
	"bytes"
	"image"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUtilSpec(t *testing.T) {
	Convey("Given you have the cake image", t, func() {
		fImg, err := os.Open("../testimages/cake.jpg")
		defer fImg.Close()
		if err != nil {
			panic(err.Error())
		}

		sourceImage, _, err := image.Decode(fImg)
		buffer := new(bytes.Buffer)

		Convey("When apply the Sepia filter", func() {
			Sepia(sourceImage, buffer)

			Convey("Then there should be at least 10kb in the buffer", func() {
				So(len(buffer.Bytes()), ShouldBeGreaterThan, 10000)
			})
		})
	})
}
