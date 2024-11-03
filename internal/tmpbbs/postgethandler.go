package tmpbbs

import (
	"html/template"
	"io"
	"net/http"
	"strconv"

	"golang.org/x/text/message"
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
{{- if .TripCode }} <span class="trip-code">!{{ .TripCode }}</span>{{ end }}{{ end }} <span class="time">{{ .TimeAgo }}</span>
{{- if .Replies }} ({{ .NumRepliesLocalized }}){{ end -}}
{{ end -}}

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
<li>
{{- if .post.Parent }}
{{ template "post_title" .post.ParentDisplayPost }}
{{- else }}
&nbsp;
{{- end }}
<ul class="post">
<li>
{{ template "post" .post }}
<ul id="replies" class="post">
{{- $class := "even" }}
{{- range .post.RepliesPage .repliesPage .repliesPerPage }}
<li class="{{ $class }}">
<details open>
<summary>
{{ template "post_title" . }}
</summary>
{{- if .Body }}
{{ .BodyHTML }}
{{- end }}
</details>
</li>
{{- if eq $class "even" }}{{ $class = "odd" }}{{ else }}{{ $class = "even" }}{{ end -}}
{{ end -}}
{{ .post.RepliesNav .repliesPage .repliesPerPage $class -}}
{{ if .repliesEnabled }}
<li class="{{ $class }}">
<details open>
<summary>
Reply
</summary>
<form id="reply" action="{{ .post.URL }}" method="post">
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
</details>
</li>
{{- end }}
</ul>
</ul>
</ul>
</body>
</html>
`

func NewPostGetHandler(title string, repliesPerPage int, cssURLs []string, repliesEnabled bool, postStore *postStore) *postGetHandler {
	template := template.Must(template.New("index").Parse(html))

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
		printer := message.NewPrinter(message.MatchLanguage(r.Header.Get("Accept-Language"), "en"))

		displayPost := newDisplayPost(post, printer)
		if !displayPost.HasRepliesPage(repliesPage, pgh.repliesPerPage) {
			http.NotFound(w, r)

			return
		}

		err = pgh.renderPost(displayPost, w, repliesPage)
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

func (pgh postGetHandler) renderPost(displayPost *displayPost, w io.Writer, repliesPage int) error {
	return pgh.template.Execute(w, map[string]interface{}{
		"title":          pgh.title,
		"repliesPerPage": pgh.repliesPerPage,
		"cssURLs":        pgh.cssURLs,
		"repliesEnabled": pgh.repliesEnabled,
		"post":           displayPost,
		"repliesPage":    repliesPage,
	})
}
