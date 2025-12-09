package main

import (
	"fmt"
	"log"
	"net/http"

	"proxy/config"
	"proxy/handlers"
)

func main() {
	http.HandleFunc("/", handlers.ProxyHandler)

	addr := fmt.Sprintf("0.0.0.0:%d", config.ProxyPort)
	fmt.Printf("Proxy iniciado en http://localhost:%d\n", config.ProxyPort)
	fmt.Printf("Reenviando todas las solicitudes a %s\n", config.TargetURL)

	log.Fatal(http.ListenAndServe(addr, nil))
}
