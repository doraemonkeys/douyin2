package response

type LoginResponse struct {
	CommonResponse
	UserID int    `json:"user_id"`
	Token  string `json:"token"`
}
