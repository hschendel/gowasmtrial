package main

import (
	"net/http"
	"log"
	"github.com/hschendel/wasmtrial/shared"
	"fmt"
	"encoding/json"
	"bytes"
)

type Repository struct {
	entities []shared.SomeEntity
}

func (r *Repository) Get(index int, entity *shared.SomeEntity) error {
	if r.entities == nil || index >= len(r.entities) || index < 0 {
		return fmt.Errorf("entity %d not found", index)
	}
	*entity = r.entities[index]
	return nil
}

func main() {
	repo := &Repository{
		entities: []shared.SomeEntity{{A: "one", B: 1}},
	}
	http.HandleFunc("/entity", func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(repo.entities[0]); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(buf.Bytes()); err != nil {
			log.Println(err)
		}
	})
	log.Print("Listening at :8080")
	http.Handle("/", http.FileServer(http.Dir("_web")))
	http.ListenAndServe(":8080", nil)
}
