package usecase

import (
	"github.com/labstack/gommon/log"
	forumRep "github.com/p1xelse/VK_DB_course_project/app/internal/forum/repository"
	threadRep "github.com/p1xelse/VK_DB_course_project/app/internal/thread/repository"
	userRep "github.com/p1xelse/VK_DB_course_project/app/internal/user/repository"
	"github.com/p1xelse/VK_DB_course_project/app/models"
	"strconv"
)

type ThreadUseCaseI interface {
	CreateThread(thread *models.Thread) error
	GetForumThreads(slug string, limit int, since string, desc bool) ([]*models.Thread, error)
	GetThread(slugOrId string) (*models.Thread, error)
	UpdateThread(thread *models.Thread, slugOrId string) error
	CreateVote(vote *models.Vote, slugOrId string) (*models.Thread, error)
}

type threadUsecase struct {
	userRepo   userRep.RepositoryI
	threadRepo threadRep.RepositoryI
	forumRepo  forumRep.RepositoryI
}

func (t threadUsecase) CreateThread(thread *models.Thread) error {
	log.Info(thread.Created)
	_, err := t.userRepo.GetUsersByNickname(thread.Author)
	if err != nil {
		return err
	}

	selectedForum, err := t.forumRepo.GetForumBySlag(thread.Forum)
	if err != nil {
		return err
	}

	if thread.Slug != "" {
		existThread, err := t.threadRepo.GetThreadBySlug(thread.Slug)
		if err != models.ErrNotFound && err != nil {
			return err
		} else if err == nil {
			thread.Id = existThread.Id
			thread.Author = existThread.Author
			thread.Created = existThread.Created
			thread.Forum = existThread.Forum
			thread.Message = existThread.Message
			thread.Slug = existThread.Slug
			thread.Title = existThread.Title
			thread.Votes = existThread.Votes
			return models.ErrConflict
		}
	}

	thread.Forum = selectedForum.Slug

	err = t.threadRepo.CreateThread(thread)
	if err != nil {
		return err
	}

	err = t.forumRepo.CreateForumUser(thread.Forum, thread.Author)
	if err != nil {
		return err
	}

	return nil
}

func (t threadUsecase) GetForumThreads(slug string, limit int, since string, desc bool) ([]*models.Thread, error) {
	_, err := t.forumRepo.GetForumBySlag(slug)
	if err != nil {
		return nil, err
	}

	threads, err := t.threadRepo.GetForumThreads(slug, limit, since, desc)
	if err != nil {
		return nil, err
	}

	return threads, nil
}

func (t threadUsecase) GetThread(slugOrId string) (*models.Thread, error) {
	var thread *models.Thread
	id, err := strconv.ParseUint(slugOrId, 10, 64)
	if err == nil {
		thread, err = t.threadRepo.GetThreadById(id)
		if err != nil {
			return nil, err
		}
	} else {
		thread, err = t.threadRepo.GetThreadBySlug(slugOrId)
		if err != nil {
			return nil, err
		}
	}

	return thread, nil
}

func (t threadUsecase) UpdateThread(thread *models.Thread, slugOrId string) error {
	var selectedThread *models.Thread
	id, err := strconv.ParseUint(slugOrId, 10, 64)
	if err == nil {
		selectedThread, err = t.threadRepo.GetThreadById(id)
		if err != nil {
			return err
		}
	} else {
		selectedThread, err = t.threadRepo.GetThreadBySlug(slugOrId)
		if err != nil {
			return err
		}
		id = selectedThread.Id
	}

	if thread.Title == "" && thread.Message == "" {
		thread.Author = selectedThread.Author
		thread.Created = selectedThread.Created
		thread.Forum = selectedThread.Forum
		thread.Id = selectedThread.Id
		thread.Slug = selectedThread.Slug
		thread.Votes = selectedThread.Votes
		thread.Title = selectedThread.Title
		thread.Message = selectedThread.Message
		return nil
	}

	thread.Id = id

	err = t.threadRepo.UpdateThread(thread)
	if err != nil {
		return err
	}

	return nil
}

func (t threadUsecase) CreateVote(vote *models.Vote, slugOrId string) (*models.Thread, error) {
	_, err := t.userRepo.GetUsersByNickname(vote.NickName)
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseUint(slugOrId, 10, 64)
	if err != nil {
		thread, err := t.threadRepo.GetThreadBySlug(slugOrId)
		if err != nil {
			return nil, err
		}
		id = thread.Id
	} else {
		_, err := t.threadRepo.GetThreadById(id)
		if err != nil {
			return nil, err
		}
	}

	vote.ThreadId = id

	err = t.threadRepo.CreateVote(vote)
	if err != nil {
		return nil, err
	}

	thread, err := t.threadRepo.GetThreadById(id)
	if err != nil {
		return nil, err
	}

	return thread, nil
}

func NewThreadUsecase(tr threadRep.RepositoryI, ur userRep.RepositoryI, fr forumRep.RepositoryI) ThreadUseCaseI {
	return &threadUsecase{
		userRepo:   ur,
		threadRepo: tr,
		forumRepo:  fr,
	}
}
