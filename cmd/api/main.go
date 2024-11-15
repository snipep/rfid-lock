package main

import (
	"fmt"

	"github.com/snipep/iot/internal/server"
)

func main()  {
	server := server.NewServer()

	fmt.Println("Server running on port ", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}