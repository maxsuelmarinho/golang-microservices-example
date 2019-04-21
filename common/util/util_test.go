package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestResolveIP(t *testing.T) {
	Convey("Given a call request", t, func() {
		Convey("When", func() {
			ipAddress, err := ResolveIPFromHostFile()

			Convey("Then", func() {
				So(err, ShouldBeNil)
				So(ipAddress, ShouldNotBeNil)
				So(string(ipAddress), ShouldContainSubstring, ".")
			})
		})
	})
}
