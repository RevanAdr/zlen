package auth

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Id   int
	Name string
	Role string
}

type UserResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type WebResponse[T any] struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   T      `json:"data"`
}

func ToResponses(u []User) []UserResponse {
	results := []UserResponse{}
	for _, v := range u {
		result := UserResponse{Id: v.Id, Name: v.Name, Role: v.Role}
		results = append(results, result)
	}
	return results
}

func ToResponse(u User) UserResponse {
	return UserResponse{
		Id:   u.Id,
		Name: u.Name,
		Role: u.Role,
	}
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
