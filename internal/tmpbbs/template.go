package tmpbbs

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
)

//go:embed template
var templateFS embed.FS

var templates *template.Template

func init() {
	templateDir, err := fs.Sub(templateFS, "template")
	if err != nil {
		log.Fatal(err)
	}
	templates = template.Must(template.New("templates").ParseFS(templateDir, "*.gohtml"))
}
