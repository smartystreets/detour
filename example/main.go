package main

import (
	"log"
	"net/http"

	"github.com/smartystreets/detour/v3"
	"github.com/smartystreets/detour/v3/example/app"
)

func main() {
	router := http.NewServeMux()
	router.Handle("/", detour.New(NewProcessPaymentDetour, app.NewHandler()))
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
