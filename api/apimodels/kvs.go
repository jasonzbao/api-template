package apimodels

import (
	"github.com/jasonzbao/api-template/api/response"
)

type V1KvsResponse struct {
	response.V1ResponseBase
	Value string `json:"value"`
}

type V1KvsGetRequest struct {
	Key string `form:"key"`
}

type V1KVSPostRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
