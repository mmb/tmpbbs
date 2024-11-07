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
img {
  max-width: 100%;
}

input[type=text] {
  width: 11em;
}

li#replies-start, li#replies-end {
  height: 0em;
}

li.even {
  background-color: #eee;
}

li.odd, li.even {
  padding: 1em;
}

textarea {
  width: 100%;
  height: 7em;
}

ul.post {
  list-style-type: none;
  padding: 0em 0.5em 0em 1em;
}

ul.replies {
   padding: 0em;
}
`
