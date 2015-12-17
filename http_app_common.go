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

func RunServer(engs map[string]NeoEngine, neoURL string, port int) {

	db, err := neoism.Connect(neoURL)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("connected to %s\n", neoURL)

	for _, eng := range engs {
		EnsureIndexes(db, eng.SuggestedIndexes())
	}

	cypherRunner := NewBatchWriter(db, 1024)

	m := mux.NewRouter()
	http.Handle("/", m)

	for path, eng := range engs {
		handlers := httpHandlers{db, cypherRunner, eng}
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
