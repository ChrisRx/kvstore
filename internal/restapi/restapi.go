// Package restapi provides a REST API using a gRPC client to integrate with
// the KV service.
package restapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/ChrisRx/kvstore/internal/kvpb"
)

// New constructs a new echo server using the provided KV client connection.
func New(client kvpb.KVClient) (*echo.Echo, error) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/:key", func(c echo.Context) error {
		key := c.Param("key")
		resp, err := client.Get(context.Background(), &kvpb.GetRequest{Key: key})
		if err != nil {
			if errors.Is(err, kvpb.ErrKeyNotFound) {
				return c.JSON(http.StatusNotFound, jsonError(err))
			}
			return c.JSON(http.StatusBadRequest, jsonError(err))
		}
		return c.JSON(http.StatusOK, resp)
	})

	e.POST("/:key", func(c echo.Context) error {
		req := &kvpb.SetRequest{}
		if err := c.Bind(req); err != nil {
			return err
		}
		req.Key = c.Param("key")
		if _, err := client.Set(context.Background(), req); err != nil {
			return c.JSON(http.StatusBadRequest, jsonError(err))
		}
		return c.JSON(http.StatusOK, map[string]string{
			"status": "OK",
		})
	})

	e.DELETE("/:key", func(c echo.Context) error {
		key := c.Param("key")
		if _, err := client.Delete(context.Background(), &kvpb.DeleteRequest{Key: key}); err != nil {
			return c.JSON(http.StatusBadRequest, jsonError(err))
		}
		return c.JSON(http.StatusOK, map[string]string{
			"status": "OK",
		})
	})
	return e, nil
}

// jsonError returns the provided error in a map. This is a helper function to
// ensure the error can be JSON marshaled.
func jsonError(err error) map[string]string {
	return map[string]string{
		"error": err.Error(),
	}
}
