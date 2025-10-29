package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/trading-engine/internal/logx"
	"github.com/yourusername/trading-engine/internal/server"
	"github.com/yourusername/trading-engine/internal/storage"
)

// NewServerCommand creates the server command
func NewServerCommand() *cobra.Command {
	var (
		dbPath string
		listen string
	)

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start HTTP API server",
		Long: `Start HTTP API server for Excel VBA integration.

The server provides HTTP endpoints for all trading engine functionality,
allowing Excel VBA to call the engine via HTTP requests instead of CLI.

Both CLI and HTTP return identical JSON, ensuring transport parity.`,
		Example: `  # Start server on default port
  tf-engine server

  # Start server on custom port
  tf-engine server --listen 127.0.0.1:8080

  # Use custom database
  tf-engine server --db /path/to/trading.db`,
		RunE: func(cmd *cobra.Command, args []string) error {
			corrID := cmd.Flag("corr-id").Value.String()
			if corrID == "" {
				corrID = logx.GenerateCorrelationID()
			}
			log := logx.WithCorrelationID(corrID)

			log.WithFields(map[string]interface{}{
				"listen": listen,
				"db":     dbPath,
			}).Info("Starting HTTP server")

			// Open database
			db, err := storage.New(dbPath)
			if err != nil {
				log.WithError(err).Error("Failed to open database")
				return fmt.Errorf("failed to open database: %w", err)
			}
			defer db.Close()

			// Create server
			srv := server.NewServer(db, listen)

			// Start server in goroutine
			errChan := make(chan error, 1)
			go func() {
				if err := srv.Start(); err != nil {
					log.WithError(err).Error("Server error")
					errChan <- err
				}
			}()

			// Wait for interrupt or error
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

			select {
			case <-quit:
				log.Info("Shutdown signal received")
			case err := <-errChan:
				log.WithError(err).Error("Server failed")
				return err
			}

			// Graceful shutdown
			log.Info("Shutting down server gracefully")

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				log.WithError(err).Error("Server shutdown error")
				return fmt.Errorf("server shutdown failed: %w", err)
			}

			fmt.Println("Server stopped gracefully")
			log.Info("Server stopped successfully")
			return nil
		},
	}

	cmd.Flags().StringVar(&dbPath, "db", "./trading.db", "Path to database file")
	cmd.Flags().StringVar(&listen, "listen", "127.0.0.1:18888", "Listen address and port")

	return cmd
}
