package repository

import (
	"github.com/p1xelse/VK_DB_course_project/app/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryI interface {
	CreatePosts(posts []*models.Post) error
	UpdatePost(post *models.Post) error
	GetPostById(id uint64) (*models.Post, error)
	GetThreadPosts(id uint64, limit int, since int, desc bool, sort string) ([]*models.Post, error)
}

type postRepository struct {
	db *gorm.DB
}

func (p postRepository) CreatePosts(posts []*models.Post) error {
	tx := p.db.Create(&posts)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table posts)")
	}

	return nil
}

func (p postRepository) UpdatePost(post *models.Post) error {
	tx := p.db.Model(post).Clauses(clause.Returning{}).Updates(models.Post{Message: post.Message, IsEdited: true})
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table posts)")
	}

	return nil
}

func (p postRepository) GetPostById(id uint64) (*models.Post, error) {
	post := models.Post{}

	tx := p.db.Where("id = ?", id).Take(&post)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table posts)")
	}

	return &post, nil
}

func (p postRepository) GetThreadPosts(id uint64, limit int, since int, desc bool, sort string) ([]*models.Post, error) {
	posts := make([]*models.Post, 0, 10)

	if sort == "flat" {
		if desc {
			if since != 0 {
				tx := p.db.Limit(limit).Where("thread = ? AND id < ?", id, since).
					Order("id desc").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			} else {
				tx := p.db.Limit(limit).Where("thread = ?", id).Order("id desc").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			}
		} else {
			if since != 0 {
				tx := p.db.Limit(limit).Where("thread = ? AND id > ?", id, since).
					Order("id").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			} else {
				tx := p.db.Limit(limit).Where("thread = ?", id).Order("id").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			}
		}
	} else if sort == "tree" {
		if desc {
			if since != 0 {
				tx := p.db.Limit(limit).Where("thread = ? AND post_tree < (?)", id,
					p.db.Table("posts").Select("post_tree").Where("id = ?", since)).
					Order("post_tree desc").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			} else {
				tx := p.db.Limit(limit).Where("thread = ?", id).Order("post_tree desc").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			}
		} else {
			if since != 0 {
				tx := p.db.Limit(limit).Where("thread = ? AND post_tree > (?)", id,
					p.db.Table("posts").Select("post_tree").Where("id = ?", since)).
					Order("post_tree").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			} else {
				tx := p.db.Limit(limit).Where("thread = ?", id).Order("post_tree").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			}
		}
	} else if sort == "parent_tree" {
		if desc {
			if since != 0 {
				tx := p.db.Where("post_tree[1] IN (?)", p.db.
					Table("posts").Limit(limit).Where("parent = 0 AND thread = ? AND id < (?)", id,
					p.db.Table("posts").Select("post_tree[1]").Where("id = ?", since)).
					Order("id desc").Select("id")).Order("post_tree[1] desc, post_tree").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			} else {
				tx := p.db.Where("post_tree[1] IN (?)", p.db.
					Table("posts").Limit(limit).Where("parent = 0 AND thread = ?", id).
					Order("id desc").Select("id")).Order("post_tree[1] desc, post_tree").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			}
		} else {
			if since != 0 {
				tx := p.db.Where("post_tree[1] IN (?)", p.db.
					Table("posts").Limit(limit).Where("parent = 0 AND thread = ? AND id > (?)", id,
					p.db.Table("posts").Select("post_tree[1]").Where("id = ?", since)).
					Order("id").Select("id")).Order("post_tree").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			} else {
				tx := p.db.Where("post_tree[1] IN (?)", p.db.
					Table("posts").Limit(limit).Where("parent = 0 AND thread = ?", id).
					Order("id").Select("id")).Order("post_tree").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			}
		}
	}

	return posts, nil
}

func NewPostRepository(db *gorm.DB) RepositoryI {
	return &postRepository{
		db: db,
	}
}
