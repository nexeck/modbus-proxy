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
	serveCmd.Flags().String("device", "", "modbus device")
	serveCmd.Flags().String("uri", "", "modbus uri")
	serveCmd.Flags().Int("baudrate", 9600, "baudrate")
	serveCmd.Flags().String("comset", "8N1", "Comset")
	serveCmd.Flags().Int("port", 11502, "listening port")

	serveCmd.Flags().Duration("connect-delay", 0, "initial delay after connecting before starting communication")
	serveCmd.Flags().Duration("delay", 0, "delay so use between subsequent modbus operations")
	serveCmd.Flags().Duration("timeout", 1*time.Second, "request timeout")

	if err := viper.BindPFlags(serveCmd.Flags()); err != nil {
		panic(fmt.Errorf("could not bind viper flags"))
	}
}

func runServe(cmd *cobra.Command, args []string) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	util.LogLevel("trace", nil)

	modbusSettings := proxy.Settings{
		Device:       viper.GetString("device"),
		URI:          viper.GetString("uri"),
		Baudrate:     viper.GetInt("baudrate"),
		Comset:       viper.GetString("comset"),
		ConnectDelay: viper.GetDuration("connect-delay"),
		Delay:        viper.GetDuration("delay"),
		Timeout:      viper.GetDuration("timeout"),
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
