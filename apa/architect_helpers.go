package apa

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"project-generator/internal/architect/model"
	architectgenerators "project-generator/internal/architect/modules/generators"
	"project-generator/internal/i18n"
)

// resolveIdea returns the idea string from either the direct flag value or a file.
// When ideaFile is "-", reads from stdin. When ideaFile is non-empty, reads from
// that file path. Otherwise returns idea as-is. ideaFile takes precedence over idea.
func resolveIdea(idea, ideaFile string) (string, error) {
	if strings.TrimSpace(ideaFile) == "" {
		return idea, nil
	}
	var data []byte
	var err error
	if ideaFile == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(ideaFile)
	}
	if err != nil {
		return "", fmt.Errorf("read idea file %q: %w", ideaFile, err)
	}
	return strings.TrimSpace(string(data)), nil
}

// isInteractive returns true when stdin is an interactive terminal.
// Non-TTY contexts (pipes, tests, CI) should not trigger interactive prompts.
func isInteractive() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}

func promptText(reader *bufio.Reader, label, defaultValue string) (string, error) {
	if strings.TrimSpace(defaultValue) != "" {
		fmt.Printf("%s [%s]: ", label, defaultValue)
	} else {
		fmt.Printf("%s: ", label)
	}
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	value := strings.TrimSpace(line)
	if value == "" {
		value = strings.TrimSpace(defaultValue)
	}
	return value, nil
}

// promptMultilineText reads multiple lines until two consecutive empty lines
// or EOF (Ctrl+D). A single empty line is preserved as a paragraph break.
// This allows users to paste multi-paragraph prompts safely.
func promptMultilineText(reader *bufio.Reader, label string) (string, error) {
	fmt.Printf("%s (%s):\n", label, i18n.T("prompt.multiline.hint"))
	var lines []string
	emptyStreak := 0
	for {
		line, err := reader.ReadString('\n')
		trimmed := strings.TrimRight(line, "\r\n")
		if trimmed == "" {
			emptyStreak++
			if err != nil {
				break
			}
			if emptyStreak >= 2 {
				break
			}
			// Keep single empty lines as paragraph separators.
			lines = append(lines, "")
		} else {
			emptyStreak = 0
			lines = append(lines, trimmed)
		}
		if err != nil {
			break
		}
	}
	// If user finishes with "double Enter", remove the terminator empty line.
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	return strings.Join(lines, "\n"), nil
}

func writeArtifacts(projectRoot string, ctx model.ProjectContext, artifacts []architectgenerators.Artifact) ([]string, error) {
	written := make([]string, 0, len(artifacts))
	for _, artifact := range artifacts {
		content, err := architectgenerators.Generate(ctx, artifact)
		if err != nil {
			return written, err
		}
		fullPath := filepath.Join(projectRoot, artifact.FilePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			return written, err
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
			return written, err
		}
		written = append(written, fullPath)
	}
	return written, nil
}

