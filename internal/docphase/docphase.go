package docphase

import (
	"os"
	"path/filepath"
	"regexp"
)

var phaseHeadingPattern = regexp.MustCompile(`(?mi)^#{1,6}\s+Phase\s+[0-9]+\b`)

type Status struct {
	RelPath    string
	Exists     bool
	PhaseBased bool
}

func IsPhaseBased(content []byte) bool {
	return phaseHeadingPattern.Match(content)
}

func Check(root, relPath string) (Status, error) {
	status := Status{RelPath: relPath}
	content, err := os.ReadFile(filepath.Join(root, relPath))
	if err != nil {
		if os.IsNotExist(err) {
			return status, nil
		}
		return status, err
	}
	status.Exists = true
	status.PhaseBased = IsPhaseBased(content)
	return status, nil
}
