package middleware

import "github.com/galaxy-future/BridgX/cmd/api/middleware/validation"

func Init() {
	validation.RegisterCustomValidators()
}
