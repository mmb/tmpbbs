package tmpbbs

import "net/http"

type CSSHandler struct{}

func (rh CSSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/css")
	_, err := w.Write([]byte(css))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var css = `
ul.post {
  list-style-type: none;
}

li.odd {
  background-color: #eee;
}

li.odd, li.even {
  padding: 1em;
}

textarea {
  width: 100%;
  height: 7em;
}
`
