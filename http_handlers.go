package neoutil

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/Financial-Times/neoism"
	"net/http"
)

type httpHandlers struct {
	db *neoism.Database
	cw CypherRunner
	nc NeoEngine
}

func (hh *httpHandlers) putHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	dec := json.NewDecoder(req.Body)
	inst, docID, err := hh.nc.DecodeJSON(dec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if docID != id {
		http.Error(w, fmt.Sprintf("id does not match: '%v' '%v'", docID, id), http.StatusBadRequest)
		return
	}

	err = hh.nc.CreateOrUpdate(hh.cw, inst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (hh *httpHandlers) deleteHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	deleted, err := hh.nc.Delete(hh.cw, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if deleted {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (hh *httpHandlers) getHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	obj, found, err := hh.nc.Read(hh.cw, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
