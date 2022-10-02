package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/satori/go.uuid"
)

func main() {
	uid := uuid.Must(uuid.NewV4(), nil)
	fmt.Printf("UUIDv4: %s\n", uid)

	addr := ":8000"
	log.Printf("listen on %s", addr)
	http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Write([]byte(`<!doctype html>
<div id=app></div>
<script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/node-uuid/1.4.8/uuid.min.js"></script>
<script>
// alert(uuidv4());
const $app = document.getElementById('app')
const source = new EventSource('/data')
source.addEventListener('message', (ev) => {
	console.log(ev)
    app.innerHTML += ev.data + "<br>"
})
</script>`))
		}

		if r.URL.Path == "/data" {
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.WriteHeader(http.StatusOK)

			for {
				data := time.Now().Format(time.RFC3339)
				fmt.Fprintf(w, "data: %s\n\n", data)
				log.Printf("data: %s\n", data)
				w.(http.Flusher).Flush()

				select {
				case <-time.After(time.Second * 3):
				case <-r.Context().Done():
					return
				}
			}
		}
	}))
}
