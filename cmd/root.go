package cmd

import (
	"github.com/celio001/prodify/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "productfy-api",
	Short: "teste teste",
	Long:  "teste teste",
	Args:  cobra.MaximumNArgs(1),
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		logger.Log.Fatal("error while executing command", zap.String("error", err.Error()))
	}
}
