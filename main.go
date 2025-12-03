package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/neilyinliang/k620/global"
	"github.com/neilyinliang/k620/server"
)

var config global.Config

func init() {
	config = global.Config{
		AllowUsers:              "a420aa94-5e8a-415d-9537-484be3774daa",
		IntervalSecond:          "7200",
		EnableDataUsageMetering: "true",
		BufferSize:              "8192",
		AppPort:                 "8226",
	}
}

func main() {
	flag.Parse()

	// Parse subcommands
	args := flag.Args()
	if len(args) == 0 {
		args = append(args, "run") // default to "run" if no subcommand is provided
	}

	subcommand := strings.ToLower(args[0])
	switch subcommand {
	case "run":
		runServer()
	default:
		fmt.Printf("Unknown subcommand: %s\n\n", subcommand)
		os.Exit(1)
	}
}

func runServer() {
	fd := global.SetupLogger("", "DEBUG")
	defer fd.Close()

	// Channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	app := server.NewApp(&config, stop)

	go app.Run()
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	app.Shutdown(ctx)
}
