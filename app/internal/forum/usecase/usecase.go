package usecase

import (
	forumRep "github.com/p1xelse/VK_DB_course_project/app/internal/forum/repository"
	userRep "github.com/p1xelse/VK_DB_course_project/app/internal/user/repository"
	"github.com/p1xelse/VK_DB_course_project/app/models"
)

type ForumUseCaseI interface {
	CreateForum(forum *models.Forum) error
	GetForumBySlug(slug string) (*models.Forum, error)
	GetForumUsers(slug string, limit int, since string, desc bool) ([]*models.User, error)
}

type forumUsecase struct {
	forumRepo forumRep.RepositoryI
	userRepo  userRep.RepositoryI
}

func (f forumUsecase) CreateForum(forum *models.Forum) error {
	user, err := f.userRepo.GetUsersByNickname(forum.User)
	if err != nil {
		return err
	}

	existForum, err := f.forumRepo.GetForumBySlag(forum.Slug)
	if err == nil {
		forum.User = existForum.User
		forum.Slug = existForum.Slug
		forum.Title = existForum.Title
		forum.Threads = existForum.Threads
		forum.Posts = existForum.Posts
		return models.ErrConflict
	} else if err != models.ErrNotFound {
		return err
	}

	forum.User = user.Nickname
	err = f.forumRepo.CreateForum(forum)

	if err != nil {
		return err
	}

	return nil
}

func (f forumUsecase) GetForumUsers(slug string, limit int, since string, desc bool) ([]*models.User, error) {
	_, err := f.forumRepo.GetForumBySlag(slug)

	if err != nil {
		return nil, err
	}

	users, err := f.forumRepo.GetForumUsers(slug, limit, since, desc)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (f forumUsecase) GetForumBySlug(slug string) (*models.Forum, error) {
	forum, err := f.forumRepo.GetForumBySlag(slug)

	if err != nil {
		return nil, err
	}

	return forum, nil
}

func NewForumUsecase(fr forumRep.RepositoryI, ur userRep.RepositoryI) ForumUseCaseI {
	return &forumUsecase{
		forumRepo: fr,
		userRepo:  ur,
	}
}
