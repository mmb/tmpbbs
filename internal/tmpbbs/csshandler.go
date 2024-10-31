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
label {
  display: inline-block;
  width: 50px;
  text-align: right;
}

textarea {
  height: 200px;
}

textarea, input {
  width: 400px;
}

ul.post {
  list-style-type: none;
}

li.odd {
  background-color: #eee;
}

li.odd, li.even {
  padding: 1em;
}
`
