package usecase

import (
	"github.com/go-openapi/strfmt"
	forumRep "github.com/p1xelse/VK_DB_course_project/app/internal/forum/repository"
	postRep "github.com/p1xelse/VK_DB_course_project/app/internal/post/repository"
	threadRep "github.com/p1xelse/VK_DB_course_project/app/internal/thread/repository"
	userRep "github.com/p1xelse/VK_DB_course_project/app/internal/user/repository"
	"github.com/p1xelse/VK_DB_course_project/app/models"
	"strconv"
	"time"
)

type PostUseCaseI interface {
	CreatePosts(posts []*models.Post, slugOrId string) error
	UpdatePost(post *models.Post) error
	GetPost(id uint64, related []string) (*models.PostDetails, error)
	GetThreadPosts(slugOrId string, limit int, since int, desc bool, sort string) ([]*models.Post, error)
}

type postsUsecase struct {
	postRepo   postRep.RepositoryI
	userRepo   userRep.RepositoryI
	threadRepo threadRep.RepositoryI
	forumRepo  forumRep.RepositoryI
}

func (p postsUsecase) CreatePosts(posts []*models.Post, slugOrId string) error {
	var thread *models.Thread
	id, err := strconv.ParseUint(slugOrId, 10, 64)
	if err == nil {
		thread, err = p.threadRepo.GetThreadById(id)
		if err != nil {
			return err
		}
	} else {
		thread, err = p.threadRepo.GetThreadBySlug(slugOrId)
		if err != nil {
			return err
		}
	}

	if len(posts) == 0 {
		return nil
	}

	for idx := range posts {
		posts[idx].Thread = thread.Id
		posts[idx].Forum = thread.Forum
		_, err = p.userRepo.GetUsersByNickname(posts[idx].Author)
		if err != nil {
			return err
		}
		if posts[idx].Parent != 0 {
			selectedPost, err := p.postRepo.GetPostById(posts[idx].Parent)
			if err == models.ErrNotFound {
				return models.ErrConflict
			} else if err != nil {
				return err
			} else if selectedPost.Thread != posts[idx].Thread {
				return models.ErrConflict
			}
		}
	}

	timeNow := time.Now()
	for idx := range posts {
		posts[idx].Created = strfmt.DateTime(timeNow)
	}

	err = p.postRepo.CreatePosts(posts)
	if err != nil {
		return err
	}

	for idx := range posts {
		err = p.forumRepo.CreateForumUser(posts[idx].Forum, posts[idx].Author)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p postsUsecase) UpdatePost(post *models.Post) error {
	selectedPost, err := p.postRepo.GetPostById(post.Id)
	if err != nil {
		return err
	}

	if post.Message == "" || post.Message == selectedPost.Message {
		post.Id = selectedPost.Id
		post.Message = selectedPost.Message
		post.IsEdited = selectedPost.IsEdited
		post.Author = selectedPost.Author
		post.Created = selectedPost.Created
		post.Forum = selectedPost.Forum
		post.Parent = selectedPost.Parent
		post.Thread = selectedPost.Thread
		return nil
	}

	err = p.postRepo.UpdatePost(post)
	if err != nil {
		return err
	}

	post.IsEdited = true

	return nil
}

func (p postsUsecase) GetPost(id uint64, related []string) (*models.PostDetails, error) {
	postDetails := models.PostDetails{}

	post, err := p.postRepo.GetPostById(id)
	if err != nil {
		return nil, err
	}

	postDetails.Post = post

	for _, elem := range related {
		switch elem {
		case "user":
			user, err := p.userRepo.GetUsersByNickname(post.Author)
			if err != nil {
				return nil, err
			}
			postDetails.User = user
		case "thread":
			thread, err := p.threadRepo.GetThreadById(post.Thread)
			if err != nil {
				return nil, err
			}
			postDetails.Thread = thread
		case "forum":
			forum, err := p.forumRepo.GetForumBySlag(post.Forum)
			if err != nil {
				return nil, err
			}
			postDetails.Forum = forum
		}
	}

	return &postDetails, nil
}

func (p postsUsecase) GetThreadPosts(slugOrId string, limit int, since int, desc bool, sort string) ([]*models.Post, error) {
	var selectedThread *models.Thread
	id, err := strconv.ParseUint(slugOrId, 10, 64)
	if err == nil {
		selectedThread, err = p.threadRepo.GetThreadById(id)
		if err != nil {
			return nil, err
		}
	} else {
		selectedThread, err = p.threadRepo.GetThreadBySlug(slugOrId)
		if err != nil {
			return nil, err
		}
	}

	posts, err := p.postRepo.GetThreadPosts(selectedThread.Id, limit, since, desc, sort)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func NewPostUsecase(pr postRep.RepositoryI, ur userRep.RepositoryI, tr threadRep.RepositoryI, fr forumRep.RepositoryI) PostUseCaseI {
	return &postsUsecase{
		postRepo:   pr,
		userRepo:   ur,
		threadRepo: tr,
		forumRepo:  fr,
	}
}
