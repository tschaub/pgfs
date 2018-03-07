package cmd

import (
	"fmt"
	"os"

	_ "github.com/lib/pq" // only works with postgres
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pgfs",
	Short: "Postgres backed WFS 3",
}

// Execute runs the command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
