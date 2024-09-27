package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) CreateExpenseHandler(c echo.Context) error {
	exp := Expense{}
	err := c.Bind(&exp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	// Append to the Books table
	h.DB.Create(&exp)

	return c.JSON(http.StatusCreated, exp)
}
