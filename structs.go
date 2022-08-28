package justrpc

import "encoding/json"

type requestShape struct {
	Version  Version         `json:"justrpc"`
	Method   string          `json:"method"`
	ID       string          `json:"id"`
	ArgsData json.RawMessage `json:"args"`
	MetaData json.RawMessage `json:"meta"`
}

type responseInfo struct {
	Version Version `json:"justrpc"`
	Success bool    `json:"success"`
	ID      string  `json:"id,omitempty"`
}

type okResponse struct {
	responseInfo
	ResultData json.RawMessage `json:"result,omitempty"`
	MetaData   json.RawMessage `json:"meta,omitempty"`
}

type errorResponse struct {
	responseInfo
	Error *Error `json:"error"`
}
