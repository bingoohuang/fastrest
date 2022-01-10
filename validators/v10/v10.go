package v10

import (
	"github.com/bingoohuang/fastrest"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	fastrest.RegisterPreProcessor(fastrest.PreProcessorFn(func(dtx *fastrest.Context) error {
		return validate.Struct(dtx.Req)
	}))
}
