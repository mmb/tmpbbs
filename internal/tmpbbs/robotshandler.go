package tmpbbs

import "net/http"

type RobotsHandler struct{}

func (rh RobotsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("User-agent: *\nDisallow: /\n"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
