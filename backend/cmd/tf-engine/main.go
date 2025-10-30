package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tf-engine <command> [options]")
		fmt.Println("Commands:")
		fmt.Println("  server     Start HTTP API server")
		fmt.Println("  init       Initialize database")
		fmt.Println("  settings   Manage settings")
		fmt.Println("  size       Calculate position size")
		fmt.Println("  checklist  Evaluate checklist")
		fmt.Println("  heat       Check heat levels")
		fmt.Println("  gates      Check gates")
		os.Exit(1)
	}

	command := os.Args[1]
	os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	switch command {
	case "server":
		ServerCommand()
	case "init":
		fmt.Println("TODO: Implement init command")
		os.Exit(1)
	case "settings":
		fmt.Println("TODO: Implement settings command")
		os.Exit(1)
	case "size":
		fmt.Println("TODO: Implement size command")
		os.Exit(1)
	case "checklist":
		fmt.Println("TODO: Implement checklist command")
		os.Exit(1)
	case "heat":
		fmt.Println("TODO: Implement heat command")
		os.Exit(1)
	case "gates":
		fmt.Println("TODO: Implement gates command")
		os.Exit(1)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
