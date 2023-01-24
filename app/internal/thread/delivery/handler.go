package delivery

import (
	"errors"
	"github.com/labstack/echo/v4"
	threadUsecase "github.com/p1xelse/VK_DB_course_project/app/internal/thread/usecase"
	"github.com/p1xelse/VK_DB_course_project/app/models"
	"net/http"
	"strconv"
)

type Delivery struct {
	ThreadUC threadUsecase.ThreadUseCaseI
}

func (delivery *Delivery) CreateThread(c echo.Context) error {
	var thread models.Thread
	err := c.Bind(&thread)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	thread.Forum = c.Param("slug")

	err = delivery.ThreadUC.CreateThread(&thread)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		case errors.Is(err, models.ErrConflict):
			c.Logger().Error(err)
			return c.JSON(http.StatusConflict, thread)
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, thread)
}

func (delivery *Delivery) GetForumThreads(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 100
	}

	since := c.QueryParam("since")

	desc, err := strconv.ParseBool(c.QueryParam("desc"))
	if err != nil {
		desc = false
	}

	threads, err := delivery.ThreadUC.GetForumThreads(c.Param("slug"), limit, since, desc)
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

	return c.JSON(http.StatusOK, threads)
}

func (delivery *Delivery) GetThread(c echo.Context) error {
	thread, err := delivery.ThreadUC.GetThread(c.Param("slug_or_id"))

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

	return c.JSON(http.StatusOK, thread)
}

func (delivery *Delivery) UpdateThread(c echo.Context) error {
	var thread models.Thread
	err := c.Bind(&thread)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	err = delivery.ThreadUC.UpdateThread(&thread, c.Param("slug_or_id"))
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

	return c.JSON(http.StatusOK, thread)
}

func (delivery *Delivery) CreateVote(c echo.Context) error {
	var vote models.Vote
	err := c.Bind(&vote)
	if err != nil || (vote.Voice != -1 && vote.Voice != 1) {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	thread, err := delivery.ThreadUC.CreateVote(&vote, c.Param("slug_or_id"))

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

	return c.JSON(http.StatusOK, thread)
}

func NewDelivery(e *echo.Echo, tu threadUsecase.ThreadUseCaseI) {
	handler := &Delivery{
		ThreadUC: tu,
	}

	e.POST("/api/forum/:slug/create", handler.CreateThread)
	e.GET("/api/forum/:slug/threads", handler.GetForumThreads)
	e.GET("/api/thread/:slug_or_id/details", handler.GetThread)
	e.POST("/api/thread/:slug_or_id/details", handler.UpdateThread)
	e.POST("/api/thread/:slug_or_id/vote", handler.CreateVote)
}
