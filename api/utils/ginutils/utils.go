package ginutils

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"

	apiErrors "github.com/jasonzbao/api-template/api/errors"
)

func ShouldBindWith(c *gin.Context, obj interface{}, binding binding.Binding) error {
	err := apiErrors.NewJSONDecodeError(apiErrors.NewValidationError(c.ShouldBindWith(obj, binding)))
	if err != nil {
		return errors.Wrap(apiErrors.ErrorBadBody, err.Error())
	}
	return nil
}
