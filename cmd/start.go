package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

var startCmd = &cobra.Command{
	Use: "start",
	Short: "Start microservice",
	Long: "Start microservice",
	Run: func(cmd *cobra.Command, args []string) {
		debugMode := viper.GetBool("DEBUG")
		if debugMode {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
			if viper.GetBool("PRETTY_LOG") {
				log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			}

			log.Info().Str("version", viper.GetString("GIT_COMMIT")).Send()
			e := log.Debug()
			for _, k := range viper.GetViper().AllKeys() {
				e = e.Str(k, fmt.Sprintf("%v", viper.Get(k)))
			}
			e.Msg("settings")
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}

		// Run gRPC-Gateway servers
		ctx := context.Background()
		if err := start.StartServers(ctx); err != nil {
			log.Warn().Err(err).Msg("serving failed")
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Configuring HTTP server
	startCmd.Flags().String("http_endpoint", "127.0.0.1:8080", "Set HTTP server endpoint")
	_ = viper.BindPFlag("HTTP_ENDPOINT", startCmd.Flags().Lookup("http_endpoint"))

	// Configure Tool HTTP server
	startCmd.Flags().String("tool_server_endpoint", "127.0.0.1:7080", "Set Tool HTTP server endpoint")
	_ = viper.BindPFlag("TOOL_SERVER_ENDPOINT", startCmd.Flags().Lookup("tool_server_endpoint"))
}