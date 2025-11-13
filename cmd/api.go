package cmd

import (
	"log"

	"github.com/celio001/prodify/internal/fiber"
	"github.com/celio001/prodify/pkg/lifecycle"
	"github.com/celio001/prodify/pkg/postgress"
	"github.com/celio001/prodify/product"
	"github.com/spf13/cobra"
)

var (
	httpCommand = &cobra.Command{
		Use: "http",
		Short: "Initializes the codebase running as http server",
		Long: "Initializes the codebase running as http server",
		RunE: ApiExecute,
	}
)

func init() {
	rootCmd.AddCommand(httpCommand)
}

func ApiExecute(cmd *cobra.Command, args []string) error {

	db, err := postgress.NewInstance()

	if err != nil {
		log.Printf("%v", err)
	}

	productRepository := product.NewRepository(db)

	s := fiber.CreateServer(productRepository)

	lifecycle.New(cmd.Context(), "http", s.Start, s.Stop)

	return nil
}
