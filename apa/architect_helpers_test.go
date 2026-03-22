package apa

import (
	"bufio"
	"strings"
	"testing"
)

func TestPromptMultilineTextAllowsParagraphBreaks(t *testing.T) {
	input := "第一段第一行\n第一段第二行\n\n第二段第一行\n\n\n"
	reader := bufio.NewReader(strings.NewReader(input))

	got, err := promptMultilineText(reader, "Project idea")
	if err != nil {
		t.Fatalf("promptMultilineText returned error: %v", err)
	}

	want := "第一段第一行\n第一段第二行\n\n第二段第一行"
	if got != want {
		t.Fatalf("unexpected multiline text\nwant:\n%q\ngot:\n%q", want, got)
	}
}

func TestPromptMultilineTextStopsOnEOF(t *testing.T) {
	input := "line1\nline2"
	reader := bufio.NewReader(strings.NewReader(input))

	got, err := promptMultilineText(reader, "Project idea")
	if err != nil {
		t.Fatalf("promptMultilineText returned error: %v", err)
	}

	want := "line1\nline2"
	if got != want {
		t.Fatalf("unexpected multiline text on EOF\nwant: %q\ngot: %q", want, got)
	}
}
