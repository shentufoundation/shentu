package types

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	Describe = Convey
	It       = Convey
)

func TestContext(t *testing.T) {
	Describe("Context initialization", t, func() {
		It("should be initialized successfully", func() {
			ctx, err := NewContextWithDefaultConfigAndLogger()
			So(err, ShouldBeNil)
			So(ctx.Config(), ShouldHaveSameTypeAs, Config{})
			So(ctx.Logger(), ShouldNotBeNil)
		})
	})
}
