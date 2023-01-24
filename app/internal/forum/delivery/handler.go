package delivery

import (
	"github.com/labstack/echo/v4"
	forumUsecase "github.com/p1xelse/VK_DB_course_project/app/internal/forum/usecase"
	"github.com/p1xelse/VK_DB_course_project/app/models"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type Delivery struct {
	ForumUC forumUsecase.ForumUseCaseI
}

func (delivery *Delivery) CreateForum(c echo.Context) error {
	var forum models.Forum

	err := c.Bind(&forum)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	err = delivery.ForumUC.CreateForum(&forum)

	if err != nil {
		c.Logger().Error(err)
		switch {
		case errors.Is(err, models.ErrConflict):
			return c.JSON(http.StatusConflict, forum)
		case errors.Is(err, models.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError)
		}
	}

	return c.JSON(http.StatusCreated, forum)
}

func (delivery *Delivery) GetForum(c echo.Context) error {
	slug := c.Param("slug")
	forum, err := delivery.ForumUC.GetForumBySlug(slug)

	if err != nil {
		c.Logger().Error(err)
		switch {
		case errors.Is(err, models.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError)
		}
	}

	return c.JSON(http.StatusOK, forum)
}

func (delivery *Delivery) GetForumUsers(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 100
	}

	since := c.QueryParam("since")

	desc, err := strconv.ParseBool(c.QueryParam("desc"))
	if err != nil {
		desc = false
	}

	users, err := delivery.ForumUC.GetForumUsers(c.Param("slug"), limit, since, desc)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, users)
}

func NewDelivery(e *echo.Echo, fu forumUsecase.ForumUseCaseI) {
	handler := &Delivery{
		ForumUC: fu,
	}

	e.POST("/api/forum/create", handler.CreateForum)
	e.GET("/api/forum/:slug/details", handler.GetForum)
	e.GET("/api/forum/:slug/users", handler.GetForumUsers)
}
