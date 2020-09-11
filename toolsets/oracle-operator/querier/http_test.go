package querier

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	querierTypes "github.com/certikfoundation/shentu/toolsets/oracle-operator/querier/types"
	runnerTypes "github.com/certikfoundation/shentu/toolsets/oracle-operator/runner/types"
	oracleTypes "github.com/certikfoundation/shentu/toolsets/oracle-operator/types"
)

var (
	Describe = Convey
	It       = Convey
)

func SimulateEndpointServer() {
	http.HandleFunc("/security/primitive", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, "128\n")
	})
	fmt.Println("Start fake endpoint server")
	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal(err)
	}
}

func TestHandleRequest(t *testing.T) {
	go SimulateEndpointServer()
	Describe("Get Request", t, func() {
		It("should send request with query correctly", func() {
			var wg sync.WaitGroup
			wg.Add(1)

			expectResult := uint8(128)

			ctx, err := oracleTypes.NewContextWithDefaultConfigAndLogger()
			So(err, ShouldBeNil)

			qConf := querierTypes.Config{
				Endpoint:   "http://localhost:1234/security/primitive",
				Method:     "GET",
				Timeout:    300,
				RetryTimes: 3,
			}

			var cConf = oracleTypes.Config{
				Runner:  runnerTypes.DefaultConfig(),
				Querier: qConf,
			}
			ctxForTest := ctx.WithConfig(&cConf)

			msgCreateTask := oracleTypes.MsgCreateTask{
				Contract: "0xf3585fcd969502624c6a8acf73721d1fce214e83",
				Function: "0x00000101",
			}

			resp, err := HandleRequest(ctxForTest, msgCreateTask, qConf.Endpoint)
			So(err, ShouldBeNil)
			wg.Done()
			So(resp, ShouldHaveSameTypeAs, expectResult)
			So(resp, ShouldEqual, expectResult)
		})
	})
}

func TestPackRequest(t *testing.T) {
	Describe("Pack GET Request", t, func() {
		It("should have correct value as data source", func() {
			ctx, err := oracleTypes.NewContextWithDefaultConfigAndLogger()
			So(err, ShouldBeNil)

			qConf := querierTypes.Config{
				Endpoint:   "http://localhost:1234/security/primitive",
				Method:     "GET",
				Timeout:    300,
				RetryTimes: 3,
			}

			var cConf = oracleTypes.Config{
				Runner:  runnerTypes.DefaultConfig(),
				Querier: qConf,
			}
			ctxForTest := ctx.WithConfig(&cConf)
			config := ctxForTest.Config()

			msgCreateTask := oracleTypes.MsgCreateTask{
				Contract: "0xf3585fcd969502624c6a8acf73721d1fce214e83",
				Function: "0x00000101",
			}

			v := url.Values{}
			v.Add("address", msgCreateTask.Contract)
			v.Add("functionSignature", msgCreateTask.Function)

			u, err := url.Parse(qConf.Endpoint)
			So(err, ShouldBeNil)
			u.RawQuery = v.Encode()
			endpoint := u.String()

			expect, _ := http.NewRequest(config.Querier.Method, endpoint, bytes.NewBuffer([]byte("")))
			actual, err := packRequest(ctxForTest, msgCreateTask, endpoint)
			So(err, ShouldBeNil)
			So(fmt.Sprint(actual), ShouldEqual, fmt.Sprint(expect))
		})
	})
}
