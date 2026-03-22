package apa

import (
	"fmt"
	"os"

	"project-generator/internal/config"
)

func defaultLocalSkillsPath() (string, error) {
	return config.ExpandPath("skills")
}

func resolveExistingSkillsPath(path string) (string, bool, error) {
	absPath, err := config.ExpandPath(path)
	if err != nil {
		return "", false, err
	}

	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return absPath, false, nil
		}
		return "", false, err
	}
	if !info.IsDir() {
		return "", false, fmt.Errorf("skills path is not a directory: %s", absPath)
	}
	return absPath, true, nil
}
