package delivery

import (
	"github.com/labstack/echo/v4"
	serviceUsecase "github.com/p1xelse/VK_DB_course_project/app/internal/service/usecase"
	"net/http"
)

type Delivery struct {
	ServiceUC serviceUsecase.ServiceUseCaseI
}

func (delivery *Delivery) ClearData(c echo.Context) error {
	err := delivery.ServiceUC.ClearData()
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (delivery *Delivery) GetStatus(c echo.Context) error {
	status, err := delivery.ServiceUC.GetStatus()
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, status)
}

func NewDelivery(e *echo.Echo, serviceUC serviceUsecase.ServiceUseCaseI) {
	handler := &Delivery{
		ServiceUC: serviceUC,
	}

	e.POST("/api/service/clear", handler.ClearData)
	e.GET("/api/service/status", handler.GetStatus)
}
