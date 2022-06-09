package errors

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"reflect"
	"strings"
)

type Validator struct {
	*validator.Validate
}

func InitValidator() *Validator {
	return &Validator{
		validator.New(),
	}
}

func (v *Validator) Struct(s interface{}) error {
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
	return v.StructCtx(context.Background(), s)
}

func (v *Validator) ValidateRequest(ctx *gin.Context, obj any) []*IError {
	err := ctx.ShouldBindJSON(&obj)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, nil)
	}

	err = v.Struct(obj)
	if err != nil {
		return v.buildJsonErrors(err)
	}

	return nil
}

type IError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

func (v *Validator) buildJsonErrors(err error) []*IError {

	var Errors []*IError
	for _, err := range err.(validator.ValidationErrors) {
		var el IError
		el.Field = strings.ToLower(err.Field())
		el.Field = err.Field()
		el.Tag = err.Tag()
		el.Value = err.Param()
		Errors = append(Errors, &el)
	}
	return Errors
}
