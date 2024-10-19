package tmpbbs

import "net/http"

func RobotsHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("User-agent: *\nDisallow: /\n"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
