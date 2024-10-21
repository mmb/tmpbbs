package tmpbbs

import "net/http"

func CSSHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/css")
	_, err := w.Write([]byte(`
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
`))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
