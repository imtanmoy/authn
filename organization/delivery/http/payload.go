package http

import (
	"gopkg.in/thedevsaddam/govalidator.v1"
	"net/http"
	"net/url"
)

type OrganizationPayload struct {
	Name string `json:"name"`
}

func (o *OrganizationPayload) Bind(r *http.Request) error {
	return nil
}

func (o *OrganizationPayload) validate() url.Values {
	rules := govalidator.MapData{
		"name": []string{"required", "min:4", "max:20"},
	}
	opts := govalidator.Options{
		Data:  o,
		Rules: rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateStruct()
	return e
}
