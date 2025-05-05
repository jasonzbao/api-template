package response

import (
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	apiErrors "github.com/jasonzbao/api-template/api/errors"
	"github.com/jasonzbao/api-template/api/logger"
)

const (
	KeyRequestResponse = "request-apistructs-key"
	KeyHttpStatusCode  = "apistructs-http-status-code-key"
)

type V1ResponseBase struct {
	HttpStatusCode int    `json:"status_code,omitempty"`
	Error          error  `json:"-"`
	ErrorCode      int    `json:"error_code,omitempty"`
	ErrorString    string `json:"error_str,omitempty"`
}

func formatReturn(c *gin.Context, res interface{}, rb *V1ResponseBase) {
	if rb.Error != nil {
		rb.ErrorString = rb.Error.Error()
		var apiError *apiErrors.APIError
		if errors.As(rb.Error, &apiError) {
			rb.ErrorCode = apiError.InternalCode
			c.Set(KeyHttpStatusCode, apiError.StatusCode)
		} else {
			c.Set(KeyHttpStatusCode, http.StatusInternalServerError)
		}
	} else {
		c.Set(KeyHttpStatusCode, rb.HttpStatusCode)
	}
	// done so that there are no 0 status codes
	rb.HttpStatusCode = c.GetInt(KeyHttpStatusCode)

	if rb.HttpStatusCode >= 500 && rb.HttpStatusCode != http.StatusNotFound {
		logger.WithContext(c).Error(
			fmt.Sprintf("[ERROR] | %d | %s", rb.HttpStatusCode, c.Request.URL.Path),
			zap.Any("error", rb.Error),
		)
	}

	c.Set(KeyRequestResponse, res)
}

func (rb *V1ResponseBase) FormatReturn(c *gin.Context, res interface{}) {
	// if panic, just re-panic but add this so we don't send a 200 prematurely
	if e := recover(); e != nil {
		logger.WithContext(c).Error(
			fmt.Sprintf("[ERROR] | %d | %s", rb.HttpStatusCode, c.Request.URL.Path),
			zap.Any("error", e),
		)
		panic(e)
	}
	formatReturn(c, res, rb)
}

func middleware(c *gin.Context, formatFunc func(*gin.Context, interface{})) {
	defer func() {
		if c.IsAborted() {
			return
		}

		res, requestResponseKeyExists := c.Get(KeyRequestResponse)
		if !requestResponseKeyExists {
			res = &V1ResponseBase{
				HttpStatusCode: http.StatusInternalServerError,
				Error:          apiErrors.ErrorNoResponse,
			}
			c.Set(KeyHttpStatusCode, http.StatusInternalServerError)
		}
		formatFunc(c, res)
	}()
	c.Next()
}

func Middleware(c *gin.Context) {
	middleware(c, func(c *gin.Context, res interface{}) {
		c.JSON(c.GetInt(KeyHttpStatusCode), res)
	})
}
