package tmpbbs

import (
	"html/template"
	"io"
	"net/http"
	"strconv"
)

type postGetHandler struct {
	title          string
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
{{- range .post.Replies }}
<li>
<details open>
<summary>{{ template "post_title" . }}</summary>
{{- if .Body }}
{{ .BodyHTML }}
{{- end }}
</details>
</li>
{{- end -}}
{{ if .repliesEnabled }}
<li>
<details open>
<summary>Reply</summary>
<form action="{{ .post.URL }}" method="post">
<p>
<label for="title">Title</label>
<input type="text" id="title" name="title">
</p>
<p>
<label for="author">Author</label>
<input type="text" id="author" name="author">
</p>
<p>
<label for="body"></label>
<textarea id="body" name="body"></textarea>
</p>
<p>
<label for="body"></label>
<input type="submit" value="Reply">
</p>
</form>
</details>
</li>
{{- end }}
</ul>
{{- if .post.Parent }}
</ul>
{{- end }}
</ul>
</body>
</html>
`

func NewPostGetHandler(title string, cssURLs []string, repliesEnabled bool, postStore *postStore) *postGetHandler {
	return &postGetHandler{
		title:          title,
		cssURLs:        cssURLs,
		repliesEnabled: repliesEnabled,
		template:       template.Must(template.New("index").Parse(html)),
		postStore:      postStore,
	}
}

func (pgh postGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := castID(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if !pgh.postStore.get(id, func(post *post) {
		err = pgh.renderPost(post, w)
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

func (pgh postGetHandler) renderPost(post *post, w io.Writer) error {
	return pgh.template.Execute(w, map[string]interface{}{
		"title":          pgh.title,
		"cssURLs":        pgh.cssURLs,
		"repliesEnabled": pgh.repliesEnabled,
		"post":           post,
	})
}
