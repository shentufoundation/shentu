package types

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadLogger(t *testing.T) {
	Describe("Initialize Logger", t, func() {
		It("should initialize logger successfully", func() {
			err := initLogger()
			if err != nil {
				logger.Error(err.Error())
			}
			So(func() { logger.Debug("logger test") }, ShouldNotPanic)
		})
	})
}
