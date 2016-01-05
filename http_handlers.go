package neoutil

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type httpHandlers struct {
	ne NeoEngine
}

func (hh *httpHandlers) putHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	dec := json.NewDecoder(req.Body)
	inst, docID, err := hh.ne.DecodeJSON(dec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if docID != id {
		http.Error(w, fmt.Sprintf("id does not match: '%v' '%v'", docID, id), http.StatusBadRequest)
		return
	}

	err = hh.ne.CreateOrUpdate(inst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (hh *httpHandlers) deleteHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	deleted, err := hh.ne.Delete(id)

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

	obj, found, err := hh.ne.Read(id)

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
