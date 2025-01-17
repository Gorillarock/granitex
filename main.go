package main

import (
	"fmt"
	"granitex/db"
	"granitex/server"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	fmt.Println("Hello, World!")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//e.GET("/", handler)
	e.Static("/v1/src", "client/src")
	e.File("/v1/tx", "client/tx_client.html")
	e.File("/v1/rx", "client/rx_client.html")
	e.POST("/v1/handler/tx", server.PostTxHandler)
	e.GET("/v1/handler/rx", server.GetRxHandler)

	err := db.NewDBClient()
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Logger.Fatal(e.Start(":1323"))
}
