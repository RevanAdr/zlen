package auth

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
