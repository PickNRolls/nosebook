package domainlike

import "github.com/google/uuid"

type PostResource struct {
	id    uuid.UUID
	owner *UserOwner
}

func (this *PostResource) Id() uuid.UUID {
	return this.id
}

func (this *PostResource) Type() ResourceType {
	return POST_RESOURCE
}

func (this *PostResource) Owner() Owner {
	return this.owner
}

func NewPostResource(id uuid.UUID, authorId uuid.UUID) *PostResource {
	return &PostResource{
		id:    id,
		owner: NewUserOwner(authorId),
	}
}

type CommentResource struct {
	id    uuid.UUID
	owner *UserOwner
}

func (this *CommentResource) Id() uuid.UUID {
	return this.id
}

func (this *CommentResource) Type() ResourceType {
	return COMMENT_RESOURCE
}

func (this *CommentResource) Owner() Owner {
	return this.owner
}

func NewCommentResource(id uuid.UUID, authorId uuid.UUID) *CommentResource {
	return &CommentResource{
		id:    id,
		owner: NewUserOwner(authorId),
	}
}
