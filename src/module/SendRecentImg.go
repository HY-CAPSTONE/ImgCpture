package module

import (
	"io/ioutil"
	"log"
	"net/http"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {

	buf, err := ioutil.ReadFile("imgStuck/output.jpg")

	if err != nil {

		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(buf)
}
