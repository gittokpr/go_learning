package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

// implement the ServeHTTP method to make Hello a handler (similar to implementing methods of an interface)
func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log.Println("Hello World")
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Oops", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(rw, "Hello %s", d)
}
