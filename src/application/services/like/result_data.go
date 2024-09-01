package like

import "github.com/google/uuid"

type resultData struct {
	PostId    *uuid.UUID `json:"postId,omitempty"`
	CommentId *uuid.UUID `json:"commentId,omitempty"`
	Liked     bool       `json:"liked"`
}
