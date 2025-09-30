package request

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type RegisterRequest struct {
	UserRequest
}

type LoginRequest struct {
	UserRequest
}
