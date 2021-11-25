package middleware

import (
	"github.com/galaxy-future/BridgX/cmd/api/middleware/validation"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func Init() {
	validation.RegisterCustomValidators()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := validation.RegisterValidators(v)
		if err != nil {
			logs.Logger.Fatal(err.Error())
		}
	}
}
