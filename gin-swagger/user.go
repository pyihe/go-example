package gin_swagger

type LoginResponse struct {
	Username string `json:"username,omitempty"`
	Token    string `json:"token,omitempty"`
}
