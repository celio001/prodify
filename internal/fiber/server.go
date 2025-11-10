package fiber

import (
	"context"

	v1 "github.com/celio001/prodify/internal/fiber/v1"
)

func (h HttpServer) Start(ctx context.Context) {

	router := h.app.Group("/api/")

	v1Router := router.Group(v1.HandlerPath)
	v1.RegisterRouter(v1Router, h.productRepository)
}
