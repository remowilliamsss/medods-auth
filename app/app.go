package app

import (
	"context"
	"github.com/gocraft/web"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Run(ctx context.Context) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Print(err)
		}
	}()
	dao := NewSessionDAO(client)
	jwt := NewJWT([]byte(os.Getenv("secret")))
	service := NewService(dao, jwt)
	handler := NewHandler(service)

	return http.ListenAndServe("localhost:8080", initEndpoints(handler))
}

func initEndpoints(h *Handler) *web.Router {
	router := web.New(*h)
	router.Get("/auth", WrapEndpoint(h.Auth))
	router.Get("/auth/refresh", WrapEndpoint(h.Refresh))

	return router
}
