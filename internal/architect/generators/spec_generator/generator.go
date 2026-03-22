package spec_generator

import (
	"project-generator/internal/architect/model"
	"project-generator/internal/architect/modules/templates"
)

func Generate(ctx model.ProjectContext) (string, error) {
	return templates.Render("SPEC", ctx)
}
