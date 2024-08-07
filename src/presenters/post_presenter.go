package presenters

import (
	"math"
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
	postRepository interfaces.PostRepository
}

func NewPostPresenter() *PostPresenter {
	return &PostPresenter{}
}

func (p *PostPresenter) WithPostingService(s *services.PostingService) *PostPresenter {
	p.postingService = s
	return p
}

func (p *PostPresenter) WithPostRepository(repo interfaces.PostRepository) *PostPresenter {
	p.postRepository = repo
	return p
}

func (p *PostPresenter) FindByFilter(filter dto.QueryFilterDTO, a *auth.Auth) *dto.QueryResultDTO {
	result := p.postingService.FindByFilter(&commands.FindPostsCommand{
		Filter: structs.QueryFilter{
			OwnerId:  filter.OwnerId,
			AuthorId: filter.AuthorId,
			Cursor:   filter.Cursor,
		},
	}, a)

	resultDTO := &dto.QueryResultDTO{
		Err:            result.Err,
		Next:           result.Next,
		RemainingCount: result.RemainingCount,
	}

	if len(result.Data) > 0 {
		authorMap := make(map[uuid.UUID]bool)
		ownerMap := make(map[uuid.UUID]bool)
		likerMap := make(map[uuid.UUID]bool)

		for _, post := range result.Data {
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
			resultDTO.Err = err
			return resultDTO
		}

		owners, err := p.postRepository.FindOwners(ownerIds)
		if err != nil {
			resultDTO.Err = err
			return resultDTO
		}

		likers, err := p.postRepository.FindLikers(likerIds)
		if err != nil {
			resultDTO.Err = err
			return resultDTO
		}

		resultDTO.Data = make([]*dto.PostDTO, len(result.Data))

		for i, post := range result.Data {
			postDTO := &dto.PostDTO{}
			postDTO.Id = post.Id
			postDTO.Message = post.Message
			postDTO.CreatedAt = post.CreatedAt

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

			resultDTO.Data[i] = postDTO
		}
	}

	return resultDTO
}
