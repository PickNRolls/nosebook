package presenters

import (
	"math"
	"nosebook/src/domain/posts"
	"nosebook/src/errors"
	"nosebook/src/presenters/post_presenter/dto"
	"nosebook/src/presenters/post_presenter/interfaces"
	"nosebook/src/services"
	"nosebook/src/services/auth"
	"nosebook/src/services/posting/commands"
	"nosebook/src/services/posting/structs"
	"slices"

	"github.com/google/uuid"
)

type PostPresenter struct {
	postingService *services.PostingService
	commentService *services.CommentService
	postRepository interfaces.PostRepository
}

func NewPostPresenter() *PostPresenter {
	return &PostPresenter{}
}

func (p *PostPresenter) WithPostingService(s *services.PostingService) *PostPresenter {
	p.postingService = s
	return p
}

func (p *PostPresenter) WithCommentService(s *services.CommentService) *PostPresenter {
	p.commentService = s
	return p
}

func (p *PostPresenter) WithPostRepository(repo interfaces.PostRepository) *PostPresenter {
	p.postRepository = repo
	return p
}

func (p *PostPresenter) mapPosts(posts []*posts.Post, a *auth.Auth) ([]*dto.PostDTO, *errors.Error) {
	if len(posts) == 0 {
		return make([]*dto.PostDTO, 0), nil
	}

	authorMap := make(map[uuid.UUID]bool)
	ownerMap := make(map[uuid.UUID]bool)
	likerMap := make(map[uuid.UUID]bool)

	for _, post := range posts {
		if _, has := authorMap[post.AuthorId]; !has {
			authorMap[post.AuthorId] = true
		}

		if _, has := ownerMap[post.OwnerId]; !has {
			ownerMap[post.OwnerId] = true
		}

		l := math.Min(float64(len(post.LikedBy)), 5)
		for _, likerId := range post.LikedBy[:int(l)] {
			if _, has := likerMap[likerId]; !has {
				likerMap[likerId] = true
			}
		}
	}

	authorIds := make([]uuid.UUID, 0)
	for id := range authorMap {
		authorIds = append(authorIds, id)
	}

	ownerIds := make([]uuid.UUID, 0)
	for id := range ownerMap {
		ownerIds = append(ownerIds, id)
	}

	likerIds := make([]uuid.UUID, 0)
	for id := range likerMap {
		likerIds = append(likerIds, id)
	}

	authors, err := p.postRepository.FindAuthors(authorIds)
	if err != nil {
		return nil, errors.From(err)
	}

	owners, err := p.postRepository.FindOwners(ownerIds)
	if err != nil {
		return nil, errors.From(err)
	}

	likers, err := p.postRepository.FindLikers(likerIds)
	if err != nil {
		return nil, errors.From(err)
	}

	result := make([]*dto.PostDTO, len(posts))

	for i, post := range posts {
		postDTO := &dto.PostDTO{}
		postDTO.Id = post.Id
		postDTO.Message = post.Message
		postDTO.CreatedAt = post.CreatedAt
		postDTO.LikedByUser = slices.Contains(post.LikedBy, a.UserId)
		postDTO.CanBeRemovedByUser = post.CanBeRemovedBy(a.UserId)

		for _, author := range authors {
			if post.AuthorId == author.Id {
				postDTO.Author = author
				break
			}
		}

		for _, owner := range owners {
			if post.OwnerId == owner.Id {
				postDTO.Owner = owner
				break
			}
		}

		postDTO.LikesCount = len(post.LikedBy)
		postDTO.RandomFiveLikers = make([]*dto.UserDTO, 0)
		for _, liker := range likers {
			if slices.Contains(post.LikedBy, liker.Id) {
				postDTO.RandomFiveLikers = append(postDTO.RandomFiveLikers, liker)
			}
		}

		result[i] = postDTO
	}

	return result, nil
}

func (p *PostPresenter) FindByFilter(filter dto.QueryFilterDTO, a *auth.Auth) *dto.QuerySingleResultDTO[*dto.PostDTO] {
	result := p.postingService.FindByFilter(&commands.FindPostsCommand{
		Filter: structs.QueryFilter{
			OwnerId:  filter.OwnerId,
			AuthorId: filter.AuthorId,
			Cursor:   filter.Cursor,
		},
	}, a)

	resultDTO := &dto.QuerySingleResultDTO[*dto.PostDTO]{
		Err:            result.Err,
		Next:           result.Next,
		RemainingCount: result.RemainingCount,
	}
	if resultDTO.Err != nil {
		return resultDTO
	}

	resultDTO.Data, resultDTO.Err = p.mapPosts(result.Data, a)
	return resultDTO
}

func (p *PostPresenter) Publish(c *commands.PublishPostCommand, a *auth.Auth) (*dto.PostDTO, error) {
	post, err := p.postingService.Publish(c, a)
	if err != nil {
		return nil, err
	}

	DTOs, error := p.mapPosts([]*posts.Post{post}, a)
	if error != nil {
		return nil, error
	}

	return DTOs[0], error
}

func (p *PostPresenter) Remove(c *commands.RemovePostCommand, a *auth.Auth) (*dto.PostDTO, error) {
	post, err := p.postingService.Remove(c, a)
	if err != nil {
		return nil, err
	}

	DTOs, error := p.mapPosts([]*posts.Post{post}, a)
	if error != nil {
		return nil, error
	}

	return DTOs[0], error
}

func (p *PostPresenter) Like(c *commands.LikePostCommand, a *auth.Auth) (*dto.PostDTO, error) {
	post, err := p.postingService.Like(c, a)
	if err != nil {
		return nil, err
	}

	DTOs, error := p.mapPosts([]*posts.Post{post}, a)
	if error != nil {
		return nil, error
	}

	return DTOs[0], err
}
