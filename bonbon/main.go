package bonbon

import (
	"net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.ListenAndServe(":8080", nil)
}
