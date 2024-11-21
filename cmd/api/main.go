package main

import (
	"fmt"
	// "os"
	// "os/signal"
	// "syscall"

	"github.com/snipep/iot/internal/handler"
	"github.com/snipep/iot/internal/server"
)

func main()  {
	server := server.NewServer()
	// Initialize the controller
	handler.InitializeController()

	// // Wait for system interrupt to gracefully shut down
	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	// <-sigChan

	fmt.Println("Server running on port ", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}