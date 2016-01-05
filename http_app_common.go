package neoutil

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmcvetta/neoism"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func EnsureAllIndexes(db *neoism.Database, engs map[string]NeoEngine) {
	for _, eng := range engs {
		EnsureIndexes(db, eng.SuggestedIndexes())
	}
}

func RunServer(engs map[string]NeoEngine, port int) {

	m := mux.NewRouter()
	http.Handle("/", m)

	for path, eng := range engs {
		handlers := httpHandlers{eng}
		m.HandleFunc(fmt.Sprintf("/%s/{id}", path), handlers.getHandler).Methods("GET")
		m.HandleFunc(fmt.Sprintf("/%s/{id}", path), handlers.putHandler).Methods("PUT")
		m.HandleFunc(fmt.Sprintf("/%s/{id}", path), handlers.deleteHandler).Methods("DELETE")
	}

	go func() {
		log.Printf("listening on %d", port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			log.Printf("web stuff failed: %v\n", err)
		}
	}()

	// wait for ctrl-c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Println("exiting")
}
