package main

import (
	"circadian/internal/config"
	"circadian/internal/db"
	"circadian/internal/loop"
	"circadian/internal/routes"
	"context"
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/sync/errgroup"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	configPath := "./sample/config.yaml"

	cfg, err := config.Read(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read config: %v\n", err)
		os.Exit(-1)
	}

	dbType := os.Getenv("DB_TYPE")
	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")

	conn, err := db.New(dbType, dbConnectionString)
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", conn)
			c.Set("config", cfg)
			return next(c)
		}
	})

	e.GET("/api/v1/healthz", routes.Healthz)
	routes.RegisterHealthCheckResultRoutes(e, conn)
	e.GET("/api/v1/uptime", routes.Uptime)

	e.Static("/", "ui/dist")

	// if len(os.Args) < 2 {
	// 	fmt.Fprintf(os.Stderr, "Path must be passed\n")
	// 	os.Exit(-1)
	// }
	// configPath := os.Args[1]
	// if len(configPath) == 0 {
	// 	fmt.Fprintf(os.Stderr, "Valid path not found\n")
	// 	os.Exit(-1)
	// }

	eg, egCtx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		return loop.Run(egCtx, cfg, conn)
	})

	eg.Go(func() error {
		return e.Start(":9081")
	})

	err = eg.Wait()
	if err != nil {
		e.Logger.Fatal(err)
	}
}
