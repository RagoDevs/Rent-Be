package main

import "github.com/labstack/echo/v4"

func (app *application) test(c echo.Context) error {
	return c.JSON(200, "test")
}
