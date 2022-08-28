package justrpc

import (
	"io"
	"net/http"
	"reflect"
)

var ServiceVersion = Version{1, 0}

type Service struct {
	methods map[string]reflect.Value
}

func NewService() *Service {
	return &Service{
		methods: make(map[string]reflect.Value),
	}
}

func (s *Service) Register(method string, handler any) *Service {
	handlerValue := reflect.ValueOf(handler)
	checkHandlerType(handlerValue.Type())
	s.methods[method] = handlerValue
	return s
}

func (s *Service) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Body != nil {
		defer req.Body.Close()
	}
	respInfo := responseInfo{Version: ServiceVersion}
	if req.Method != "POST" {
		respBody, _ := serializeError(respInfo, Errorf(InvalidRequestFormat, "Only POST requests are allowed"))
		rw.Header().Add("Allow", "POST")
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write(respBody)
		return
	}
	reqData, err := io.ReadAll(req.Body)
	if err != nil {
		respData, _ := serializeError(respInfo, Errorf(InvalidRequestFormat, "Cannot read request body: %v", err))
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(respData)
		return
	}
	respData, success := s.handleRequest(reqData)
	if !success {
		rw.WriteHeader(http.StatusBadRequest)
	}
	rw.Write(respData)
}
