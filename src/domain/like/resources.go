package domainlike

import "github.com/google/uuid"

type PostResource struct {
	id uuid.UUID
}

func (this *PostResource) Id() uuid.UUID {
	return this.id
}

func (this *PostResource) Type() ResourceType {
	return POST_RESOURCE
}

func NewPostResource(id uuid.UUID) *PostResource {
	return &PostResource{
		id: id,
	}
}

type CommentResource struct {
	id uuid.UUID
}

func (this *CommentResource) Id() uuid.UUID {
	return this.id
}

func (this *CommentResource) Type() ResourceType {
	return COMMENT_RESOURCE
}

func NewCommentResource(id uuid.UUID) *CommentResource {
	return &CommentResource{
		id: id,
	}
}
