package repository

import (
	"github.com/p1xelse/VK_DB_course_project/app/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type RepositoryI interface {
	ClearData() error
	GetStatus() (*models.ServiceStatus, error)
}

type serviceRepository struct {
	db *gorm.DB
}

func (s serviceRepository) ClearData() error {
	tx := s.db.Exec("TRUNCATE posts, threads, forums, users, forum_user cascade;")
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error")
	}

	return nil
}

func (s serviceRepository) GetStatus() (*models.ServiceStatus, error) {
	status := models.ServiceStatus{}

	var count int64
	tx := s.db.Model(&models.User{}).Count(&count)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table users)")
	}
	status.UserCount = count
	tx = s.db.Model(&models.Forum{}).Count(&count)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table forums)")
	}
	status.ForumCount = count
	tx = s.db.Model(&models.Thread{}).Count(&count)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table posts)")
	}
	status.ThreadCount = count
	tx = s.db.Model(&models.Post{}).Count(&count)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table threads)")
	}
	status.PostCount = count

	return &status, nil
}

func NewServiceRepository(db *gorm.DB) RepositoryI {
	return &serviceRepository{
		db: db,
	}
}
