package fiber

import (
	"context"
	"fmt"
	"log"

	v1 "github.com/celio001/prodify/internal/fiber/v1"
)

func (h HttpServer) Start(ctx context.Context) error {
	router := h.app.Group("/api/")

	v1Router := router.Group(v1.HandlerPath)
	v1.RegisterRouter(v1Router, h.productRepository)

	addr := fmt.Sprint(":8080")
	log.Printf("Starting server on %v", addr)
	return h.app.Listen(addr)
}

func (h HttpServer) Stop(ctx context.Context) error {
	return h.app.Shutdown()
}
