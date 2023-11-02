package response

import "github.com/labstack/echo/v4"

type (
	MsgReponse struct {
		Message string `json:"message"`
	}
)

func ErrorResponse(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, &MsgReponse{Message: message})
}

func SuccessResponse(c echo.Context, statusCode int, data any) error {
	return c.JSON(statusCode, data)
}
