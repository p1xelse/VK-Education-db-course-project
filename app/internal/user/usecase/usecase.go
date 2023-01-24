package usecase

import (
	"github.com/p1xelse/VK_DB_course_project/app/internal/user/repository"
	"github.com/p1xelse/VK_DB_course_project/app/models"
)

type UserUseCaseI interface {
	CreateUser(user *models.User) ([]*models.User, error)
	GetUserByNickname(nickname string) (*models.User, error)
	UpdateUser(user *models.User) error
}

type userUsecase struct {
	userRepo repository.RepositoryI
}

func (u userUsecase) UpdateUser(user *models.User) error {
	existedUser, err := u.userRepo.GetUsersByNickname(user.Nickname)
	if err != nil {
		return err
	}

	//nothing change
	if user.Email == "" && user.Fullname == "" && user.About == "" {
		user.About = existedUser.About
		user.Email = existedUser.Email
		user.Fullname = existedUser.Fullname
	}

	userByEmail, err := u.userRepo.GetUsersByEmail(user.Email)

	if err == nil && userByEmail.Nickname != user.Nickname {
		return models.ErrConflict
	} else if err != models.ErrNotFound {
		return err
	}

	err = u.userRepo.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (u userUsecase) GetUserByNickname(nickname string) (*models.User, error) {
	user, err := u.userRepo.GetUsersByNickname(nickname)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u userUsecase) CreateUser(user *models.User) ([]*models.User, error) {
	existUsers, err := u.userRepo.GetUsersByNickNameOrEmail(user.Email, user.Nickname)
	if err != nil {
		return nil, err
	} else if len(existUsers) > 0 {
		return existUsers, models.ErrConflict
	}

	err = u.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func NewServiceUsecase(ps repository.RepositoryI) UserUseCaseI {
	return &userUsecase{
		userRepo: ps,
	}
}
