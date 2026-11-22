package httpx

import (
	"encoding/json"
	"net/http"
)

type WebResponse[T any] struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   T      `json:"data"`
}

func WriteToJSON(w http.ResponseWriter, status int, str string, args any) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := WebResponse[any]{
		Code:   status,
		Status: str,
		Data:   args,
	}
	encoder.Encode(response)
}

func ReadBody(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(dst)
	if err != nil {
		return err
	}
	return nil
}
