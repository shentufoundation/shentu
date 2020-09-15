package oracle

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseMsgCreateTaskContract(t *testing.T) {
	Describe("parse MsgCreateTask target contract", t, func() {
		testcases := []struct {
			Input    string
			Expected []string
		}{
			{Input: "eth:0xa", Expected: []string{"eth", "0xa"}},
			{Input: ":0xa", Expected: []string{"", "0xa"}},
			{Input: "eth:", Expected: []string{"eth", ""}},
			{Input: "eth:ropsten:0xa", Expected: []string{"eth:ropsten", "0xa"}},
			{Input: "", Expected: []string{"eth", ""}},
			{Input: "0xa", Expected: []string{"eth", "0xa"}},
		}
		for _, tt := range testcases {
			tt := tt
			It(fmt.Sprintf("should parse task contract %s correctly", tt.Input), func() {
				prefix, address, err := parseMsgCreateTaskContract(tt.Input)
				So(err, ShouldBeNil)
				So(prefix, ShouldEqual, tt.Expected[0])
				So(address, ShouldEqual, tt.Expected[1])
			})
		}
	})
}
