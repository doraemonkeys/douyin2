package response

import "github.com/Doraemonkeys/douyin2/internal/app/models"

const CommentResFormat = "2006-01-02 15:04:05"

type PostCommentResponse struct {
	CommonResponse
	Comment Comment `json:"comment"`
}

type Comment struct {
	ID         int    `json:"id"`
	User       User   `json:"user"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
}

func (c *Comment) SetValue(comment models.CommentModel, user models.UserModel, followedMyself bool) {
	c.ID = int(comment.ID)
	c.Content = comment.Content
	c.CreateDate = comment.CreatedAt.Format(CommentResFormat)
	c.User.SetValue(user, followedMyself)
}
