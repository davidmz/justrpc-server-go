package justrpc

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	srv *Service
}

func (s *ServiceTestSuite) SetupSuite() {
	s.srv = NewService().
		Register("sum", func(in []int) int {
			sum := 0
			for _, a := range in {
				sum += a
			}
			return sum
		}).
		Register("panic", func() {
			panic("Just panic!")
		}).
		Register("untypedError", func() error {
			return errors.New("hi, I am error")
		}).
		Register("typedError", func() error {
			return Errorf("myError", "hi, I am error")
		}).
		Register("wrappedError", func() error {
			return fmt.Errorf("wrapped: %w", Errorf("myError", "hi, I am error"))
		})
}

func (s *ServiceTestSuite) TestIntArray() {
	out, success := s.srv.handleRequest([]byte(`{
		"justrpc": "1.0",
		"method": "sum",
		"args": [2,3,5]
	}`))

	s.True(success)
	s.JSONEq(`{
		"justrpc": "1.0",
		"success": true,
		"result": 10
	}`, string(out))
}

func (s *ServiceTestSuite) TestEmptyArray() {
	out, success := s.srv.handleRequest([]byte(`{
		"justrpc": "1.0",
		"method": "sum",
		"args": []
	}`))

	s.True(success)
	s.JSONEq(`{
		"justrpc": "1.0",
		"success": true,
		"result": 0
	}`, string(out))
}

func (s *ServiceTestSuite) TestNullArg() {
	out, success := s.srv.handleRequest([]byte(`{
		"justrpc": "1.0",
		"method": "sum",
		"args": null
	}`))

	s.True(success)
	s.JSONEq(`{
		"justrpc": "1.0",
		"success": true,
		"result": 0
	}`, string(out))
}

func (s *ServiceTestSuite) TestInvalidArg() {
	out, success := s.srv.handleRequest([]byte(`{
		"justrpc": "1.0",
		"method": "sum",
		"args": "null"
	}`))

	s.False(success)
	s.JSONEq(`{
		"justrpc": "1.0",
		"success": false,
		"error": {
			"code": "justrpc.invalidArgs",
			"message": "Invalid args format: json: cannot unmarshal string into Go value of type []int"
		}
	}`, string(out))
}

func (s *ServiceTestSuite) TestNoArg() {
	out, success := s.srv.handleRequest([]byte(`{
		"justrpc": "1.0",
		"method": "sum"
	}`))

	s.False(success)
	s.JSONEq(`{
		"justrpc": "1.0",
		"success": false,
		"error": {
			"code": "justrpc.invalidArgs",
			"message": "Required arguments not provided"
		}
	}`, string(out))
}

func (s *ServiceTestSuite) TestInvalidJSON() {
	out, success := s.srv.handleRequest([]byte(`{hi,!`))

	s.False(success)
	s.JSONEq(`{
		"justrpc": "1.0",
		"success": false,
		"error": {
			"code": "justrpc.invalidRequestFormat",
			"message": "Invalid request format: invalid character 'h' looking for beginning of object key string"
		}
	}`, string(out))
}

func (s *ServiceTestSuite) TestPanicInHandler() {
	out, success := s.srv.handleRequest([]byte(`{
		"justrpc": "1.0",
		"method": "panic"
	}`))

	s.False(success)
	s.JSONEq(`{
		"justrpc": "1.0",
		"success": false,
		"error": {
			"code": "justrpc.internalError",
			"message": "Panic happened: Just panic!"
		}
	}`, string(out))
}

func (s *ServiceTestSuite) TestUntypedError() {
	out, success := s.srv.handleRequest([]byte(`{
		"justrpc": "1.0",
		"method": "untypedError"
	}`))

	s.False(success)
	s.JSONEq(`{
		"justrpc": "1.0",
		"success": false,
		"error": {
			"code": "justrpc.internalError",
			"message": "hi, I am error"
		}
	}`, string(out))
}

func (s *ServiceTestSuite) TestTypedError() {
	out, success := s.srv.handleRequest([]byte(`{
		"justrpc": "1.0",
		"method": "typedError"
	}`))

	s.False(success)
	s.JSONEq(`{
		"justrpc": "1.0",
		"success": false,
		"error": {
			"code": "myError",
			"message": "hi, I am error"
		}
	}`, string(out))
}

func (s *ServiceTestSuite) TestWrappedError() {
	out, success := s.srv.handleRequest([]byte(`{
		"justrpc": "1.0",
		"method": "wrappedError"
	}`))

	s.False(success)
	s.JSONEq(`{
		"justrpc": "1.0",
		"success": false,
		"error": {
			"code": "myError",
			"message": "hi, I am error"
		}
	}`, string(out))
}

///

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
