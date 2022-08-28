package justrpc

import (
	"encoding/json"
	"errors"
	"reflect"
)

func (s *Service) handleRequest(requestData json.RawMessage) (output json.RawMessage, isSuccess bool) {
	request := new(requestShape)
	respInfo := responseInfo{
		Version: ServiceVersion,
		Success: true,
	}

	if err := json.Unmarshal(requestData, request); err != nil {
		return serializeError(respInfo, Errorf(InvalidRequestFormat, "Invalid request format: %v", err))
	}
	if request.Version.Major != ServiceVersion.Major {
		return serializeError(respInfo,
			Errorf(VersionMismatch, "This server supports only %d.* version of JustRPC", ServiceVersion.Major))
	}
	respInfo.ID = request.ID

	method, ok := s.methods[request.Method]
	if !ok {
		return serializeError(respInfo, Errorf(MethodNotFound, "Method %q not found", request.Method))
	}

	// Calling the method
	var inArgs []reflect.Value
	mType := method.Type()

	if mType.NumIn() > 0 {
		if request.ArgsData == nil {
			return serializeError(respInfo, Errorf(InvalidArgs, "Required arguments not provided"))
		}
		v := reflect.New(mType.In(0))
		if err := json.Unmarshal(request.ArgsData, v.Interface()); err != nil {
			return serializeError(respInfo, Errorf(InvalidArgs, "Invalid args format: %v", err))
		}
		inArgs = append(inArgs, v.Elem())
	}
	if mType.NumIn() > 1 {
		if request.MetaData == nil {
			return serializeError(respInfo, Errorf(InvalidArgs, "Required metadata not provided"))
		}
		v := reflect.New(mType.In(1))
		if err := json.Unmarshal(request.MetaData, v.Interface()); err != nil {
			return serializeError(respInfo, Errorf(InvalidArgs, "Invalid meta format: %v", err))
		}
		inArgs = append(inArgs, v.Elem())
	}

	defer func() {
		if pnc := recover(); pnc != nil {
			output, isSuccess = serializeError(respInfo, Errorf(InternalError, "Panic happened: %v", pnc))
		}
	}()

	results := method.Call(inArgs)
	if len(results) > 0 {
		lastResult := results[len(results)-1]
		if lastResult.Type().Implements(errorInterface) {
			if !lastResult.IsNil() {
				err := lastResult.Interface().(error)
				var e *Error
				if errors.As(err, &e) {
					err = e
				}
				return serializeError(respInfo, err)
			} else {
				results = results[:len(results)-1]
			}
		}
	}

	response := &okResponse{
		responseInfo: respInfo,
	}
	if len(results) > 0 {
		var err error
		if response.ResultData, err = json.Marshal(results[0].Interface()); err != nil {
			return serializeError(respInfo, Errorf(InternalError, "Cannot serialize response: %v", err))
		}
	}
	if len(results) > 1 {
		var err error
		if response.MetaData, err = json.Marshal(results[1].Interface()); err != nil {
			return serializeError(respInfo, Errorf(InternalError, "Cannot serialize response metadata: %v", err))
		}
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		return serializeError(respInfo, Errorf(InternalError, "Cannot serialize response: %v", err))
	}

	return responseData, true
}

func serializeError(respInfo responseInfo, err error) (json.RawMessage, bool) {
	respInfo.Success = false
	je, ok := err.(*Error)
	if !ok {
		je = Errorf(InternalError, err.Error()).(*Error)
	}

	resp := &errorResponse{
		responseInfo: respInfo,
		Error:        je,
	}

	responseData, err := json.Marshal(resp)
	if err != nil {
		resp.Error = Errorf(InternalError, "Cannot serialize error").(*Error)
		responseData, _ = json.Marshal(resp)
	}

	return responseData, false
}
