package response

type RegisterResponse struct {
	Token    string `json:"token,omitempty"`
	Username string `json:"username,omitempty"`
}

type LoginResponse struct {
	Token    string `json:"token,omitempty"`
	Username string `json:"username,omitempty"`
}
