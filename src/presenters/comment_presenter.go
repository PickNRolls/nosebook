package presenters

import (
	common_dto "nosebook/src/presenters/dto"
	post_dto "nosebook/src/presenters/post_presenter/dto"
	"nosebook/src/services"

	"github.com/google/uuid"
)

type CommentPresenter struct {
	commentService *services.CommentService
}

func NewCommentPresenter() *CommentPresenter {
	return &CommentPresenter{}
}

func (p *CommentPresenter) WithCommentService(s *services.CommentService) *CommentPresenter {
	p.commentService = s
	return p
}

func (p *CommentPresenter) BatchFindByPostIds(ids []uuid.UUID) *common_dto.BatchQueryResultDTO[*post_dto.CommentDTO] {
	batchResult := p.commentService.BatchFindByPostIds(ids)

	resultDTO := &common_dto.BatchQueryResultDTO[*post_dto.CommentDTO]{
		Err: batchResult.Err,
	}
	if resultDTO.Err != nil {
		return resultDTO
	}

	return resultDTO
}
