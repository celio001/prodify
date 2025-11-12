package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "productfy-api",
	Short: "teste teste",
	Long: "teste teste",
	Args: cobra.MaximumNArgs(1),
}

func Execute(){
	err := rootCmd.Execute()
	if err != nil {
		log.Printf("error while executing command")
	}
}