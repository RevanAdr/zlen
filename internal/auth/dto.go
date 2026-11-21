package auth

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
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

type UserCreate struct {
	Id    int
	Name  string
	Email string
	Role  string
}

type UserCreateResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type WebResponse[T any] struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   T      `json:"data"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Name     string
	Email    string
	Password string
}

type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserDisplay struct {
	Name  string `json:"name"`
	Email string `json:"email"`
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

func ToCreateRepsonse(u UserCreate) UserCreateResponse {
	return UserCreateResponse{
		Id:    u.Id,
		Name:  u.Name,
		Email: u.Email,
		Role:  u.Role,
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

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "could not hash password", err
	}
	return string(hash), nil
}

func MatchPassword(password string, hashed string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		return PasswordNotMatched
	}
	return nil
}
