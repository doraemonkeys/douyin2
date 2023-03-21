package response

// errors
const (
	ErrInvalidUsername = "用户名不合法"
	ErrInvalidPassword = "密码不合法"
)

type RegisterResponse struct {
	CommonResponse
	UserID int    `json:"user_id"`
	Token  string `json:"token"`
}
