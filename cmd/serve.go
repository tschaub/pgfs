package cmd

import (
	"database/sql"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tschaub/pgfs/pkg/handlers"
	"github.com/tschaub/pgfs/pkg/models"
)

var (
	servePort int
)

func init() {
	var defaultPort = 5000

	flags := serveCmd.Flags()
	flags.IntVar(&servePort, "port", defaultPort, "listen on this port")

	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve [connection]",
	Short: "Provide WFS",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		connection := args[0]
		db, err := sql.Open("postgres", connection)
		if err != nil {
			return err
		}
		defer db.Close()

		migrateErr := models.Migrate(db)
		if migrateErr != nil {
			return migrateErr
		}

		router := handlers.New(db)

		address := fmt.Sprintf(":%d", servePort)
		fmt.Printf("Listening on http://localhost%s\n", address)
		return router.Start(address)
	},
}
