package tmpbbs

import "golang.org/x/text/message"

// Strings in templates that need translation are duplicated here because
// gotext cannot find usage in templates and it will remove any translations
// that it doesn't find used.
// See https://github.com/golang/go/issues/51144
func neverCalled() { //nolint:unused // never called, see comment
	printer := message.NewPrinter(message.MatchLanguage("en"))

	printer.Sprintf("Title")
	printer.Sprintf("Author#tripcode-secret")
	printer.Sprintf("Markdown is supported.")
	printer.Sprintf("Insert emoji using shortcode between colons (:mushroom: becomes 🍄).")
	printer.Sprintf("Reply")
	printer.Sprintf("URL QR Code")
	printer.Sprintf("Close")
}
