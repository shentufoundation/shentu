package oracle

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/certikfoundation/shentu/toolsets/oracle-operator/types"
)

var (
	Describe = Convey
	It       = Convey
)

const testPort = "54321"

func SimulateEndpointServer() {
	http.HandleFunc("/security/primitive", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, "128\n")
	})
	fmt.Println("Start fake endpoint server")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", testPort), nil); err != nil {
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

			ctx, err := types.NewContextWithDefaultConfigAndLogger()
			So(err, ShouldBeNil)

			config := types.Config{
				Strategy:   make(map[types.Client]types.Strategy),
				Method:     "GET",
				Timeout:    300,
				RetryTimes: 3,
			}
			ctxForTest := ctx.WithConfig(&config)

			payload := types.PrimitivePayload{
				Client:   "eth",
				Address:  "0xf3585fcd969502624c6a8acf73721d1fce214e83",
				Function: "0x00000101",
				Contract: "eth:0xf3585fcd969502624c6a8acf73721d1fce214e83",
			}

			resp, err := handleRequest(ctxForTest, fmt.Sprintf("http://localhost:%s/security/primitive", testPort), payload)
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
			ctx, err := types.NewContextWithDefaultConfigAndLogger()
			So(err, ShouldBeNil)

			config := types.Config{
				Strategy:   make(map[types.Client]types.Strategy),
				Method:     "GET",
				Timeout:    300,
				RetryTimes: 3,
			}
			ctxForTest := ctx.WithConfig(&config)

			payload := types.PrimitivePayload{
				Client:   "eth",
				Address:  "0xf3585fcd969502624c6a8acf73721d1fce214e83",
				Function: "0x00000101",
				Contract: "eth:0xf3585fcd969502624c6a8acf73721d1fce214e83",
			}

			v := url.Values{}
			v.Add("address", payload.Address)
			v.Add("functionSignature", payload.Function)

			u, err := url.Parse(fmt.Sprintf("http://localhost:%s/security/primitive", testPort))
			So(err, ShouldBeNil)
			u.RawQuery = v.Encode()
			endpoint := u.String()

			expect, _ := http.NewRequest(config.Method, endpoint, bytes.NewBuffer([]byte("")))
			actual, err := packRequest(ctxForTest, endpoint, payload)
			So(err, ShouldBeNil)
			So(fmt.Sprint(actual), ShouldEqual, fmt.Sprint(expect))
		})
	})
}
