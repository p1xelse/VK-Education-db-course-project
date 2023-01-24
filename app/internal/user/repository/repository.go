package repository

import (
	"github.com/labstack/gommon/log"
	"github.com/p1xelse/VK_DB_course_project/app/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryI interface {
	CreateUser(user *models.User) error
	GetUsersByNickNameOrEmail(email string, nickname string) ([]*models.User, error)
	GetUsersByNickname(nickname string) (*models.User, error)
	GetUsersByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func (u userRepository) GetUsersByNickname(nickname string) (*models.User, error) {
	user := models.User{}

	tx := u.db.Where("nickname = ?", nickname).Take(&user)
	log.Info("param: ", nickname, "value:", user.Nickname)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table: users, method: GetUsersByNickname)")
	}

	return &user, nil
}

func (u userRepository) GetUsersByEmail(email string) (*models.User, error) {
	user := models.User{}

	tx := u.db.Where("email = ?", email).Take(&user)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table: users, method: GetUsersByNickname)")
	}

	return &user, nil
}

func (u userRepository) CreateUser(user *models.User) error {
	tx := u.db.Create(user)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table: users, method: CreateUser)")
	}

	return nil
}

func (u userRepository) GetUsersByNickNameOrEmail(email string, nickname string) ([]*models.User, error) {
	users := make([]*models.User, 0, 10)

	tx := u.db.Where("email = ? OR nickname = ?", email, nickname).Find(&users)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table: users, method: GetUsersByNickNameOrEmail)")
	}

	return users, nil
}

func (u userRepository) UpdateUser(user *models.User) error {
	tx := u.db.Model(user).Clauses(clause.Returning{}).Updates(user)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table: users, method: UpdateUser)")
	}

	return nil
}

func NewUserRepository(db *gorm.DB) RepositoryI {
	return &userRepository{
		db: db,
	}
}
