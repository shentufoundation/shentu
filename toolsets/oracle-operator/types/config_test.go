package types

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadConfig(t *testing.T) {
	Describe("Load Config", t, func() {
		It("should load config content correctly", func() {
			err := initConfig()
			So(err, ShouldBeNil)
		})
	})
}
