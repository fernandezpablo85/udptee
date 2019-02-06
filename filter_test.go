package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestFilterNothing(t *testing.T) {
	w := Filter{}
	original := "my email is pablo@mail.com and my favourite color is \x1b[34mblue\x1b[0m"
	filtered := w.filter(original)
	expected := original
	if expected != filtered {
		t.Errorf("expected '%s' but was '%s'", expected, filtered)
	}
}

func TestFilterColor(t *testing.T) {
	w := Filter{filterColors: true}
	original := "roses are red, violets are \x1b[34mblue\x1b[0m"
	filtered := w.filter(original)
	expected := "roses are red, violets are blue"
	if expected != filtered {
		t.Errorf("expected '%s' but was '%s'", expected, filtered)
	}
}

func TestFilterEmail(t *testing.T) {
	w := Filter{filterEmails: true}
	original := "my email is pablo@mail.com and i also have pablo.fernandez@mail.com"
	filtered := w.filter(original)
	expected := "my email is *****@mail.com and i also have ***************@mail.com"
	if expected != filtered {
		t.Errorf("expected '%s' but was '%s'", expected, filtered)
	}
}

func TestFilterEmailRepeated(t *testing.T) {
	w := Filter{filterEmails: true}
	original := "my email is pablo@mail.com and i also have pablo@mail.com again"
	filtered := w.filter(original)
	expected := "my email is *****@mail.com and i also have *****@mail.com again"
	if expected != filtered {
		t.Errorf("expected '%s' but was '%s'", expected, filtered)
	}
}

func TestFilterBoth(t *testing.T) {
	w := Filter{filterEmails: true, filterColors: true}
	original := "my email is pablo@mail.com my favourite color is \x1b[34mblue\x1b[0m"
	filtered := w.filter(original)
	expected := "my email is *****@mail.com my favourite color is blue"
	if expected != filtered {
		t.Errorf("expected '%s' but was '%s'", expected, filtered)
	}
}

func TestAsWriter(t *testing.T) {
	b := bytes.Buffer{}
	w := Filter{delegate: &b, filterEmails: true, filterColors: true}
	fmt.Fprintf(&w, "my email is pablo@mail.com my favourite color is \x1b[34mblue\x1b[0m")
	filtered, err := b.ReadString('\n')
	if err != nil && err != io.EOF {
		t.Error(err)
	}

	expected := "my email is *****@mail.com my favourite color is blue"
	if expected != filtered {
		t.Errorf("expected '%s' but was '%s'", expected, filtered)
	}
}
