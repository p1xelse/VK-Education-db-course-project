package repository

import (
	"github.com/p1xelse/VK_DB_course_project/app/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryI interface {
	CreateForum(forum *models.Forum) error
	GetForumBySlag(slug string) (*models.Forum, error)
	GetForumUsers(slug string, limit int, since string, desc bool) ([]*models.User, error)
	CreateForumUser(forum string, user string) error
}

type forumRepository struct {
	db *gorm.DB
}

func (f forumRepository) GetForumUsers(slug string, limit int, since string, desc bool) ([]*models.User, error) {
	users := make([]*models.User, 0, 10)

	nicknameCondition := "user_nickname < ?"
	orderCondition := "nickname desc"

	if desc == false {
		nicknameCondition = "user_nickname > ?"
		orderCondition = "nickname"
	}

	if since != "" {
		tx := f.db.Limit(limit).Where("nickname IN (?)", f.db.
			Select("user_nickname").Table("forum_user").Where("forum = ? AND "+nicknameCondition,
			slug, since)).Order(orderCondition).Find(&users)
		if tx.Error != nil {
			return nil, errors.Wrap(tx.Error, "database error (table forum_user)")
		}
	} else {
		tx := f.db.Limit(limit).Where("nickname IN (?)", f.db.
			Select("user_nickname").Table("forum_user").Where("forum = ?", slug)).
			Order(orderCondition).Find(&users)
		if tx.Error != nil {
			return nil, errors.Wrap(tx.Error, "database error (table forum_user)")
		}
	}

	return users, nil
}

func (f forumRepository) CreateForum(forum *models.Forum) error {
	tx := f.db.Create(forum)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table: forums, method: CreateForum)")
	}

	return nil
}

func (f forumRepository) GetForumBySlag(slug string) (*models.Forum, error) {
	forum := models.Forum{}

	tx := f.db.Where("slug = ?", slug).Take(&forum)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table: forums, method: GetForumBySlag)")
	}

	return &forum, nil
}

func (f forumRepository) CreateForumUser(forum string, user string) error {
	fu := models.ForumUserRelation{
		Forum: forum,
		User:  user,
	}
	tx := f.db.Table("forum_user").Clauses(clause.OnConflict{DoNothing: true}).Create(&fu)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table forum_user, method: CreateForumUser)")
	}

	return nil
}

func NewForumRepository(db *gorm.DB) RepositoryI {
	return &forumRepository{
		db: db,
	}
}
