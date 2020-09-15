package types

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBasenameAndExtension(t *testing.T) {
	Describe("getBasenameAndExtension", t, func() {
		testcases := []struct {
			File     string
			Expected []string
		}{
			{File: "a.toml", Expected: []string{"a", "toml"}},
			{File: "c", Expected: []string{"c", ""}},
			{File: ".toml", Expected: []string{"", "toml"}},
			{File: "d/e.toml", Expected: []string{"e", "toml"}},
			{File: "", Expected: []string{"", ""}},
		}
		for _, tt := range testcases {
			tt := tt
			It(fmt.Sprintf("should obtain correct base name and extension for %s", tt.File), func() {
				base, ext := getBasenameAndExtension(tt.File)
				So(base, ShouldEqual, tt.Expected[0])
				So(ext, ShouldEqual, tt.Expected[1])
			})
		}
	})
}

func TestLoadConfig(t *testing.T) {
	Describe("Load Config", t, func() {
		It("should load config content correctly", func() {
			err := initConfig()
			So(err, ShouldBeNil)
			t.Log("Config", config)
		})
	})
}
