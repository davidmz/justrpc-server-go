package justrpc_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidmz/justrpc-server"
	"github.com/stretchr/testify/suite"
)

type HTTPTestSuite struct {
	suite.Suite
	srv *httptest.Server
}

func (s *HTTPTestSuite) SetupSuite() {
	s.srv = httptest.NewServer(justrpc.NewService().
		Register("sum", func(in []int) int {
			sum := 0
			for _, a := range in {
				sum += a
			}
			return sum
		}))
}

func (s *HTTPTestSuite) TearDownSuite() {
	s.srv.Close()
}

func (s *HTTPTestSuite) TestOKResponse() {
	resp, err := http.Post(
		s.srv.URL,
		"application/json",
		bytes.NewBufferString(`{
		"justrpc": "1.0",
		"method": "sum",
		"args": [2,3,5]
		}`))

	s.Nil(err)
	s.Equal(http.StatusOK, resp.StatusCode)
	respBody, _ := io.ReadAll(resp.Body)
	s.JSONEq(`{
		"justrpc": "1.0",
		"success": true,
		"result": 10
	}`, string(respBody))
}

func (s *HTTPTestSuite) TestInvalidMethod() {
	resp, err := http.Get(s.srv.URL)
	s.Nil(err)
	s.Equal(http.StatusMethodNotAllowed, resp.StatusCode)
	s.Equal("POST", resp.Header.Get("Allow"))
	respBody, _ := io.ReadAll(resp.Body)
	s.JSONEq(`{
		"justrpc": "1.0",
		"success": false,
		"error": {
			"code": "justrpc.invalidRequestFormat",
			"message": "Only POST requests are allowed"
		}
	}`, string(respBody))
}

///

func TestHTTPTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPTestSuite))
}
