package repository

import (
	"github.com/p1xelse/VK_DB_course_project/app/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryI interface {
	CreateThread(thread *models.Thread) error
	GetThreadBySlug(slug string) (*models.Thread, error)
	GetForumThreads(slug string, limit int, since string, desc bool) ([]*models.Thread, error)
	GetThreadById(id uint64) (*models.Thread, error)
	UpdateThread(thread *models.Thread) error
	CreateVote(vote *models.Vote) error
}

type threadRepository struct {
	db *gorm.DB
}

func (t threadRepository) CreateThread(thread *models.Thread) error {
	if thread.Slug == "" {
		tx := t.db.Omit("votes", "slug").Create(thread)
		if tx.Error != nil {
			return errors.Wrap(tx.Error, "database error (table threads, method: CreateThread)")
		}
	} else {
		tx := t.db.Omit("votes").Create(thread)
		if tx.Error != nil {
			return errors.Wrap(tx.Error, "database error (table threads, method: CreateThread)")
		}
	}

	return nil
}

func (t threadRepository) GetThreadBySlug(slug string) (*models.Thread, error) {
	thread := models.Thread{}

	tx := t.db.Where("slug = ?", slug).Take(&thread)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table threads, method: GetThreadBySlug)")
	}

	return &thread, nil
}

func (t threadRepository) GetForumThreads(slug string, limit int, since string, desc bool) ([]*models.Thread, error) {
	threads := make([]*models.Thread, 0, 10)

	createdCondition := "created <= ?"
	orderCondition := "created desc"

	if desc == false {
		createdCondition = "created >= ?"
		orderCondition = "created"
	}

	if since != "" {
		tx := t.db.Limit(limit).Where("forum = ? AND "+createdCondition, slug, since).
			Order(orderCondition).Find(&threads)
		if tx.Error != nil {
			return nil, errors.Wrap(tx.Error, "database error (table threads)")
		}
	} else {
		tx := t.db.Limit(limit).Where("forum = ?", slug).Order("created desc").
			Find(&threads)
		if tx.Error != nil {
			return nil, errors.Wrap(tx.Error, "database error (table threads)")
		}
	}

	return threads, nil
}

func (t threadRepository) GetThreadById(id uint64) (*models.Thread, error) {
	thread := models.Thread{}

	tx := t.db.Where("id = ?", id).Take(&thread)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table threads)")
	}

	return &thread, nil
}

func (t threadRepository) UpdateThread(thread *models.Thread) error {
	tx := t.db.Model(thread).Clauses(clause.Returning{}).Updates(models.Thread{Message: thread.Message, Title: thread.Title})
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table threads, method: UpdateThread)")
	}

	return nil
}

func (t threadRepository) CreateVote(vote *models.Vote) error {
	tx := t.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "thread_id"}, {Name: "nickname"}},
		DoUpdates: clause.AssignmentColumns([]string{"voice"}),
	}).Create(vote)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table votes)")
	}

	return nil
}

func NewThreadRepository(db *gorm.DB) RepositoryI {
	return &threadRepository{
		db: db,
	}
}
