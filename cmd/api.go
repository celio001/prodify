package cmd

import (
	"os"

	"github.com/celio001/prodify/internal/fiber"
	"github.com/celio001/prodify/pkg/lifecycle"
	"github.com/celio001/prodify/pkg/logger"
	"github.com/celio001/prodify/pkg/postgress"
	"github.com/celio001/prodify/product"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	httpCommand = &cobra.Command{
		Use:   "http",
		Short: "Initializes the codebase running as http server",
		Long:  "Initializes the codebase running as http server",
		RunE:  ApiExecute,
	}
)

func init() {
	rootCmd.AddCommand(httpCommand)
}

func ApiExecute(cmd *cobra.Command, args []string) error {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	logger.Init(env)
	defer logger.Log.Sync()

	connPostgres, err := postgress.NewInstance()

	if err != nil {
		logger.Log.Fatal("failed to connect to Postgres", zap.String("error", err.Error()))
	}

	productRepository := product.NewRepository(connPostgres)

	s := fiber.CreateServer(productRepository)

	lifecycle.New(cmd.Context(), "product-api", s.Start, s.Stop)

	return nil
}
