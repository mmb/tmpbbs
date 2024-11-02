package tmpbbs

import (
	"html/template"
	"io"
	"net/http"
	"strconv"
)

type postGetHandler struct {
	title          string
	repliesPerPage int
	cssURLs        []string
	repliesEnabled bool
	template       *template.Template
	postStore      *postStore
}

const html = `
{{- define "post_title" -}}
<a href="{{ .URL }}">{{ .DisplayTitle }}</a>
{{- if .Author }} by <span class="author">{{ .Author }}</span>
{{- if .TripCode }} <span class="trip-code">!{{ .TripCode }}</span>{{ end }}{{ end -}}
{{ if .Replies }} ({{ len .Replies }}{{ if eq (len .Replies) 1 }} reply{{ else }} replies{{ end }}){{ end }} <span class="time">{{ .TimeAgo }}</span>
{{- end -}}

{{ define "post" -}}
<p>
{{ template "post_title" . }}
</p>
{{- if .Body }}
{{ .BodyHTML }}
{{- end }}
{{- end -}}

<!DOCTYPE html>
<html lang="en">
<head>
<title>{{ .title }}</title>
<meta charset="UTF-8">
{{- range .cssURLs }}
<link rel="stylesheet" href="{{ . }}">
{{- end }}
</head>
<body>
<ul class="post">
{{- if .post.Parent }}
<li>
<p>
{{ template "post_title" .post.Parent }}
</p>
<ul class="post">
{{- end }}
<li>
{{ template "post" .post }}
<ul class="post">
{{ $class := "even" }}
{{- range .post.RepliesPage .repliesPage .repliesPerPage }}
<li class="{{ $class }}">
<details open>
<summary>{{ template "post_title" . }}</summary>
{{- if .Body }}
{{ .BodyHTML }}
{{- end }}
</details>
</li>
{{ if eq $class "even" }}{{ $class = "odd" }}{{ else }}{{ $class = "even" }}{{ end }}
{{- end -}}
<li class="{{ $class }}">
<details open>
<summary>
Replies
<a href="{{ .post.BeginRepliesPageURL }}">begin</a>
<a href="{{ .post.PrevRepliesPageURL .repliesPage }}">prev</a>
<a href="{{ .post.NextRepliesPageURL .repliesPage .repliesPerPage }}">next</a>
<a href="{{ .post.EndRepliesPageURL .repliesPerPage }}">end</a>
</summary>
{{ if .repliesEnabled }}
<form action="{{ .post.URL }}" method="post">
<p>
<input type="text" id="title" name="title" placeholder="Title">
<input type="text" id="author" name="author" placeholder="Author">
</p>
<p>
<textarea id="body" name="body"></textarea>
</p>
<p>
<input type="submit" value="Reply">
</p>
</form>
{{- end }}
</details>
</li>
</ul>
{{- if .post.Parent }}
</ul>
{{- end }}
</ul>
</body>
</html>
`

func NewPostGetHandler(title string, repliesPerPage int, cssURLs []string, repliesEnabled bool, postStore *postStore) *postGetHandler {
	template := template.Must(template.New("index").Funcs(template.FuncMap{
		"even": func(i int) bool {
			return i%2 == 0
		}}).Parse(html))

	return &postGetHandler{
		title:          title,
		repliesPerPage: repliesPerPage,
		cssURLs:        cssURLs,
		repliesEnabled: repliesEnabled,
		template:       template,
		postStore:      postStore,
	}
}

func (pgh postGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := castID(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	repliesPage, err := strconv.Atoi(r.URL.Query().Get("p"))
	if err != nil {
		repliesPage = 1
	}

	if !pgh.postStore.get(id, func(post *post) {
		err = pgh.renderPost(post, w, repliesPage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}) {
		http.NotFound(w, r)
	}
}

func castID(id string) (int, error) {
	if id == "" {
		return 0, nil
	}

	return strconv.Atoi(id)
}

func (pgh postGetHandler) renderPost(post *post, w io.Writer, repliesPage int) error {
	return pgh.template.Execute(w, map[string]interface{}{
		"title":          pgh.title,
		"repliesPerPage": pgh.repliesPerPage,
		"cssURLs":        pgh.cssURLs,
		"repliesEnabled": pgh.repliesEnabled,
		"post":           post,
		"repliesPage":    repliesPage,
	})
}
