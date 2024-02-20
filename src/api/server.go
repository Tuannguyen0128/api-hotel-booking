package api

import (
	"ProjectPractice/src/api/router"
	"ProjectPractice/src/config"
	"fmt"
	"log"
	"net/http"
)

func Run() {
	config.Load()
	//auto.Load()
	fmt.Printf("Running in port %d...\n", config.PORT)
	listen(config.PORT)
}
func listen(port int) {
	r := router.Init()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
