package apa

import (
	"fmt"

	"github.com/spf13/cobra"

	"project-generator/internal/buildinfo"
	"project-generator/internal/i18n"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: i18n.T("version.short"),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(buildinfo.String())
		},
	}
}
