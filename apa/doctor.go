package apa

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"project-generator/internal/config"
	"project-generator/internal/i18n"
)

type doctorCheck struct {
	Name   string
	Passed bool
	Detail string
}

func newDoctorCmd() *cobra.Command {
	var skillsPath string
	var checkWrite bool

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: i18n.T("doctor.short"),
		Long:  i18n.T("doctor.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			userProvidedSkillsPath := strings.TrimSpace(skillsPath) != ""
			if skillsPath == "" {
				defaultPath, err := defaultLocalSkillsPath()
				if err != nil {
					return err
				}
				if isInteractive() {
					reader := bufio.NewReader(os.Stdin)
					val, err := promptText(reader, "Skills path (press Enter to use ./skills or skip)", defaultPath)
					if err != nil {
						return err
					}
					skillsPath = strings.TrimSpace(val)
				} else {
					skillsPath = defaultPath
				}
			}

			checks := make([]doctorCheck, 0, 4)

			goPath, err := exec.LookPath("go")
			if err != nil {
				checks = append(checks, doctorCheck{
					Name:   "Go executable",
					Passed: false,
					Detail: "go not found in PATH",
				})
			} else {
				checks = append(checks, doctorCheck{
					Name:   "Go executable",
					Passed: true,
					Detail: fmt.Sprintf("found at %s", goPath),
				})

				out, versionErr := exec.Command(goPath, "version").CombinedOutput()
				if versionErr != nil {
					checks = append(checks, doctorCheck{
						Name:   "Go version",
						Passed: false,
						Detail: strings.TrimSpace(string(out)),
					})
				} else {
					checks = append(checks, doctorCheck{
						Name:   "Go version",
						Passed: true,
						Detail: strings.TrimSpace(string(out)),
					})
				}
			}

			if checkWrite {
				wd, wdErr := os.Getwd()
				if wdErr != nil {
					checks = append(checks, doctorCheck{
						Name:   "Filesystem write",
						Passed: false,
						Detail: wdErr.Error(),
					})
				} else {
					testDir := filepath.Join(wd, fmt.Sprintf(".project-generator_doctor_%d", time.Now().UnixNano()))
					mkErr := os.Mkdir(testDir, 0o755)
					if mkErr != nil {
						checks = append(checks, doctorCheck{
							Name:   "Filesystem write",
							Passed: false,
							Detail: mkErr.Error(),
						})
					} else {
						_ = os.RemoveAll(testDir)
						checks = append(checks, doctorCheck{
							Name:   "Filesystem write",
							Passed: true,
							Detail: fmt.Sprintf("write test passed under %s", wd),
						})
					}
				}
			} else {
				checks = append(checks, doctorCheck{
					Name:   "Filesystem write",
					Passed: true,
					Detail: "skipped",
				})
			}

			if skillsPath != "" {
				absPath, pErr := config.ExpandPath(skillsPath)
				if pErr != nil {
					checks = append(checks, doctorCheck{
						Name:   "Skills path",
						Passed: false,
						Detail: pErr.Error(),
					})
				} else {
					info, statErr := os.Stat(absPath)
					if statErr != nil {
						passed := false
						detail := statErr.Error()
						if !userProvidedSkillsPath && os.IsNotExist(statErr) {
							passed = true
							detail = fmt.Sprintf("not found under %s (optional)", absPath)
						}
						checks = append(checks, doctorCheck{
							Name:   "Skills path",
							Passed: passed,
							Detail: detail,
						})
					} else if !info.IsDir() {
						checks = append(checks, doctorCheck{
							Name:   "Skills path",
							Passed: false,
							Detail: "path exists but is not a directory",
						})
					} else {
						checks = append(checks, doctorCheck{
							Name:   "Skills path",
							Passed: true,
							Detail: absPath,
						})
					}
				}
			} else {
				checks = append(checks, doctorCheck{
					Name:   "Skills path",
					Passed: true,
					Detail: "skipped",
				})
			}

			hasFailure := false
			for _, check := range checks {
				status := "PASS"
				if !check.Passed {
					status = "FAIL"
					hasFailure = true
				}
				fmt.Printf("[%s] %s: %s\n", status, check.Name, check.Detail)
			}

			if hasFailure {
				return fmt.Errorf("doctor found issues")
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&skillsPath, "skills-path", "", i18n.T("doctor.flag.skills-path"))
	cmd.Flags().BoolVar(&checkWrite, "check-write", true, i18n.T("doctor.flag.check-write"))
	return cmd
}
