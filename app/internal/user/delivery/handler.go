package delivery

import (
	"github.com/labstack/echo/v4"
	userUsecase "github.com/p1xelse/VK_DB_course_project/app/internal/user/usecase"
	"github.com/p1xelse/VK_DB_course_project/app/models"
	"github.com/pkg/errors"
	"net/http"
)

type Delivery struct {
	UserUC userUsecase.UserUseCaseI
}

func (delivery *Delivery) CreateUser(c echo.Context) error {
	var user models.User

	err := c.Bind(&user)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	user.Nickname = c.Param("nickname")
	conflictUsers, err := delivery.UserUC.CreateUser(&user)

	if err != nil {
		c.Logger().Error(err)
		switch {
		case errors.Is(err, models.ErrConflict):
			return c.JSON(http.StatusConflict, conflictUsers)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError)
		}
	}

	return c.JSON(http.StatusCreated, user)
}

func (delivery *Delivery) UpdateUser(c echo.Context) error {
	var user models.User

	err := c.Bind(&user)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	user.Nickname = c.Param("nickname")
	err = delivery.UserUC.UpdateUser(&user)

	if err != nil {
		c.Logger().Error(err)
		switch {
		case errors.Is(err, models.ErrConflict):
			return c.JSON(http.StatusConflict, models.ErrConflict.Error())
		case errors.Is(err, models.ErrNotFound):
			return c.JSON(http.StatusNotFound, models.ErrNotFound.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError)
		}
	}

	return c.JSON(http.StatusOK, user)
}

func (delivery *Delivery) GetUser(c echo.Context) error {
	nickname := c.Param("nickname")
	user, err := delivery.UserUC.GetUserByNickname(nickname)

	if err != nil {
		c.Logger().Error(err)
		switch {
		case errors.Is(err, models.ErrNotFound):
			return c.JSON(http.StatusNotFound, models.ErrNotFound.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError)
		}
	}

	return c.JSON(http.StatusOK, user)
}

func NewDelivery(e *echo.Echo, uu userUsecase.UserUseCaseI) {
	handler := &Delivery{
		UserUC: uu,
	}
	e.POST("/api/user/:nickname/create", handler.CreateUser)
	e.GET("/api/user/:nickname/profile", handler.GetUser)
	e.POST("/api/user/:nickname/profile", handler.UpdateUser)
}
