package apimodels

import (
	"github.com/jasonzbao/api-template/api/response"
)

type V1PingResponse struct {
	response.V1ResponseBase
	Version string `json:"version"`
}
