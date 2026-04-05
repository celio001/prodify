package fiber

import (
	"context"
	"fmt"
	"net/http"
	"os"

	v1 "github.com/celio001/prodify/internal/fiber/v1"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/celio001/prodify/pkg/logger"
)

// @title prodify API
// @version 0.1
// @description API documentation for Hercules
// @host localhost:8080
// @BasePath /api/
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func (h HttpServer) Start(ctx context.Context) error {
	
	router := h.app.Group("/api/")

	h.app.Get("/health", healthCheck)

	swagcontent, _ := os.ReadFile("docs/swagger.json")
	h.app.Get("/docs/swagger.json", adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(swagcontent); err != nil {
			logger.Log.Error("Failed to write swagger content: " + err.Error())
		}
	}))

	h.app.Get("/", adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/index.html", http.StatusMovedPermanently)
	}))

	h.app.Get("/docs/*", adaptor.HTTPHandlerFunc(httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"),
	)))

	h.app.Get("/api/health", healthCheck)

	v1Router := router.Group(v1.HandlerPath)
	v1.RegisterRouter(v1Router, h.productRepository, h.auth_service, h.userService)

	addr := fmt.Sprint(":8080")
	logger.Log.Info("Starting server on " + addr)
	return h.app.Listen(addr)
}

func (h HttpServer) Stop(ctx context.Context) error {
	return h.app.Shutdown()
}

// HealthCheck godoc
// @Summary Health check
// @Description Returns API status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func healthCheck(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"status": "ok",
	})
}
