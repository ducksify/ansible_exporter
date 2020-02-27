package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	h1 := func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "Hello from a HandleFunc #1!\n")
	}
	h2 := func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "Hello from a HandleFunc #2!\n")
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/endpoint", h2)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

// package main
//
// import (
//   "net/http"
//   "io"
//   log "github.com/Sirupsen/logrus"
//   "github.com/prometheus/client_golang/prometheus/promhttp"
// )
//
// func main() {
//
//   h1 := func(w http.ResponseWriter, _ *http.Request) {
//   	io.WriteString(w, "Hello from a HandleFunc #1!\n")
//   }
//   h2 := func(w http.ResponseWriter, _ *http.Request) {
//     io.WriteString(w, "Hello from a HandleFunc #2!\n")
//   }
//
//   http.HandleFunc("/", h1)
//   http.HandleFunc("/endpoint", h2)
//   //This section will start the HTTP server and expose
//   //any metrics on the /metrics endpoint.
//   //http.Handle("/metrics", promhttp.Handler())
//
//   log.Info("Beginning to serve on port :8080")
//   log.Fatal(http.ListenAndServe(":8080", nil))
// }
//
//
// package main
//
// import (
// 	"io"
// 	"log"
// 	"net/http"
// )
//
// func main() {
// 	h1 := func(w http.ResponseWriter, _ *http.Request) {
// 		io.WriteString(w, "Hello from a HandleFunc #1!\n")
// 	}
// 	h2 := func(w http.ResponseWriter, _ *http.Request) {
// 		io.WriteString(w, "Hello from a HandleFunc #2!\n")
// 	}
//
// 	http.HandleFunc("/", h1)
// 	http.HandleFunc("/endpoint", h2)
//
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }