/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/util/modbus"
	"github.com/nexeck/modbus-proxy/proxy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run modbus proxy service",
	Run:   runServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	serveCmd.Flags().StringP("device", "d", "", "modbus device")
	serveCmd.Flags().StringP("uri", "u", "", "modbus uri")
	serveCmd.Flags().IntP("port", "p", 11502, "listening port")

	if err := viper.BindPFlags(serveCmd.Flags()); err != nil {
		panic(fmt.Errorf("could not bind viper flags"))
	}
}

func runServe(cmd *cobra.Command, args []string) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	util.LogLevel("trace", nil)

	modbusSettings := modbus.Settings{
		Device:   viper.GetString("device"),
		URI:      viper.GetString("uri"),
		Baudrate: 9600,
		Comset:   "8N1",
	}

	if err := proxy.StartProxy(viper.GetInt("port"), modbusSettings, false); err != nil {
		panic(err)
	}

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	const shutdownTimeout = 5

	_, cancel := context.WithTimeout(context.Background(), shutdownTimeout*time.Second)
	defer cancel()
}
