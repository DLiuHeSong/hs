package main

import (
	"hs/internal"
	"hs/internal/middleware"
)

func main() {
	router := internal.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recovery)
	r := router.Group("/a")
	r.Get("/b", func(context *internal.Context) {
		context.JSON(200, "okok")
	})

	router.Start(":8080")
}
