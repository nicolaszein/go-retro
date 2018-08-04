package handlers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code   int               `json:"-"`
	Errors map[string]string `json:"errors"`
}

type DefaultResponse struct {
	Response
	Data interface{} `json:"data"`
}

func (r DefaultResponse) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Code)

	json.NewEncoder(w).Encode(r)
}

func (r *Response) AddError(key string, value string) {
	if r.Errors == nil {
		r.Errors = make(map[string]string)
	}
	r.Errors[key] = value
}
