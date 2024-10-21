package tmpbbs

import (
	"html/template"
	"net/http"
	"strconv"
	"time"
)

const html = `
{{- define "post_title" -}}
<a href="{{ .URL }}">{{ .Title }}</a>
{{- if .Author }} by {{ .Author }}{{ end -}}
{{ if .Replies }} ({{ len .Replies }} replies){{ end }} {{ .TimeAgo -}}
{{- end -}}

{{ define "post" -}}
<div class="post">
<p>
{{ template "post_title" . }}
</p>
{{- if .Body }}
{{ .BodyHTML }}
{{- end }}
</div>
{{- end -}}

<!DOCTYPE html>
<html lang="en">
<head>
<title>{{ .title }}</title>
<meta charset="UTF-8">
<link rel="stylesheet" href="{{ .cssURL }}">
</head>
<body>
<ul class="post">
{{- if .post.Parent }}
<li>
<div class="post">
<p>
{{ template "post_title" .post.Parent }}
</p>
</div>
</li>
<li>
<ul class="post">
{{- end }}
<li>
{{ template "post" .post }}
</li>
<li>
<ul class="post">
{{- range .post.Replies }}
<li>
{{ template "post" . }}
</li>
{{- end }}
<li>
<div class="post">
<form action="" method="post">
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
</div>
</li>
</ul>
{{- if .post.Parent }}
</li>
</ul>
</li>
{{- end }}
</ul>
</body>
</html>
`

func CreatePostPostHandler(postStore *postStore, tripCoder *tripCoder) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		parentID, err := castID(r.PathValue("parentID"))
		if err != nil {
			http.NotFound(w, r)
			return
		}

		p := newPost(r.FormValue("title"), r.FormValue("author"), r.FormValue("body"), tripCoder)
		postStore.put(p, parentID)

		http.Redirect(w, r, r.URL.String(), http.StatusSeeOther)
	}
}

func CreateGetPostHandler(postStore *postStore, cssURL *string, title *string) func(http.ResponseWriter, *http.Request) {
	t := template.Must(template.New("index").Parse(html))

	return func(w http.ResponseWriter, r *http.Request) {
		id, err := castID(r.PathValue("id"))
		if err != nil {
			http.NotFound(w, r)
			return
		}

		found := postStore.get(id, func(post *post) {
			err = t.Execute(w, map[string]interface{}{
				"cssURL":     cssURL,
				"post":       post,
				"title":      title,
				"timeFormat": time.DateTime,
			})
		})

		if !found {
			http.NotFound(w, r)
		}
	}
}

func castID(id string) (int, error) {
	if id == "" {
		return 0, nil
	}

	return strconv.Atoi(id)
}
