package server

import (
	"api-hotel-booking/internal/router"
	"fmt"
	"log"
	"net/http"
)

func Run() {
	listen(3000)
	fmt.Printf("Running in port %d...\n", 3000)
}
func listen(port int) {
	r := router.Init()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
