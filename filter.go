package main

import (
	"io"
	"regexp"
	"strings"
)

var (
	// regex for terminal color markers
	color = regexp.MustCompile(`\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[mGK]`)
	// regex for email addresses
	emails = regexp.MustCompile("([a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)")
)

// Filter is a Writer wrapper that optionally filters or masks unwanted characters
type Filter struct {
	// writer where to forward the received input
	delegate io.Writer
	// whether to filter out terminal color markers
	filterColors bool
	// whether to filter out email addresses
	filterEmails bool
}

// clean text according to the filter settings
func (w *Filter) filter(s string) string {
	if w.filterColors {
		s = color.ReplaceAllString(s, "")
	}
	if w.filterEmails {
		s = maskAllEmails(s)
	}
	return s
}

// replaces email address with a masked version of the handle
func maskEmail(mail string) string {
	parts := strings.Split(mail, "@")
	mask := strings.Repeat("*", len(parts[0]))
	return mask + "@" + parts[1]
}

// replaces all email addresses with masked versions
func maskAllEmails(s string) string {
	matches := emails.FindAllString(s, -1)
	for _, m := range matches {
		s = strings.Replace(s, m, maskEmail(m), -1)
	}
	return s
}

// writes the input to the delegate according to the filter settings
func (w *Filter) Write(p []byte) (int, error) {
	s := string(p)
	bs := []byte(w.filter(s))
	i, err := w.delegate.Write(bs)
	if err != nil {
		return i, err
	}
	return len(p), nil
}
