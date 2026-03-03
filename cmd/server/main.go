package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/helix-acme-corp-demo/logpipe"

	"github.com/helix-acme-corp-demo/notification-service/internal/handler"
	"github.com/helix-acme-corp-demo/notification-service/internal/sender"
	"github.com/helix-acme-corp-demo/notification-service/internal/store"
)

func main() {
	logger := logpipe.New()

	notifStore := store.New()
	notifSender := sender.NewLog(logger)
	notifHandler := handler.New(notifStore, notifSender, logger)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(logpipe.Middleware(logger))
	r.Use(middleware.Recoverer)

	r.Get("/health", handler.Health())
	r.Post("/notifications", notifHandler.Create())
	r.Get("/notifications", notifHandler.List())
	r.Get("/notifications/{id}", notifHandler.Get())

	port := "8080"
	logger.Info("starting notification-service", logpipe.String("port", port))
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
