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

type WebResponse struct {
	Code   int
	Status string
	Data   any
}

func ToResponses(u []User) []UserResponse {
	results := []UserResponse{}
	for _, v := range u {
		result := UserResponse{Id: v.Id, Name: v.Name, Role: v.Role}
		results = append(results, result)
	}
	return results
}

func WriteToJSON(w http.ResponseWriter, status int, str string, args any) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := WebResponse{
		Code:   status,
		Status: str,
		Data:   args,
	}
	encoder.Encode(response)
}
