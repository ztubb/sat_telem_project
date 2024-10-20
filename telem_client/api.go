package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo"
)

func StartApi() {
	fmt.Println("STARTING UP API")

	db := pg.Connect(&pg.Options{
		Addr:     "db:5432",
		User:     "postgres",
		Password: os.Getenv("POSTGRES_PASSWORD"),
	})
	defer db.Close()

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/sat/:id", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))

		var data []TelemEntity

		err := db.Model(&data).Where("satellite_id = ?", id).Column("*").Select()
		if err != nil {
			fmt.Println("Unable to select")
			fmt.Println(err)
			panic(err)
		}

		fmt.Println(data)

		return c.JSONPretty(http.StatusOK, data, "  ")
	})

	e.GET("/range", func(c echo.Context) error {
		var data []TelemEntity

		start := c.QueryParam("start")
		end := c.QueryParam("end")

		if start == "" {
			return c.JSONPretty(http.StatusInternalServerError, "must provide start query parameter", "  ")
		}

		if end == "" {
			return c.JSONPretty(http.StatusInternalServerError, "must provide end query parameter", "  ")
		}

		err := db.Model(&data).Column("*").
			Where("created_at >= ?", start).
			Where("created_at <= ?", end).
			Select()

		if err != nil {
			fmt.Println("Unable to select")
			fmt.Println(err)
			panic(err)
		}

		fmt.Println(data)

		return c.JSONPretty(http.StatusOK, data, "  ")
	})

	port := os.Getenv("API_PORT")
	e.Logger.Fatal(e.Start(":" + port))
}
