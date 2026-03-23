package generator

// docsFiles returns the planned doc files for the given docsType.
//   basic – core design docs only (PRD, SPEC, ARCHITECTURE)
//   full  – all docs including API, DEPLOYMENT, TESTING
func docsFiles(docsType string) []plannedFile {
	basic := []plannedFile{
		{RelPath: "docs/PRD.md", TemplateName: "docs_prd_md.tmpl", Mode: 0o644},
		{RelPath: "docs/SPEC.md", TemplateName: "docs_spec_md.tmpl", Mode: 0o644},
		{RelPath: "docs/ARCHITECTURE.md", TemplateName: "docs_architecture_md.tmpl", Mode: 0o644},
		{RelPath: "docs/IMPLEMENTATION_STATUS.md", TemplateName: "implementation_status_md.tmpl", Mode: 0o644},
	}
	if docsType == "full" {
		return append(basic,
			plannedFile{RelPath: "docs/API.md", TemplateName: "docs_api_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: "docs/DEPLOYMENT.md", TemplateName: "docs_deployment_md.tmpl", Mode: 0o644},
			plannedFile{RelPath: "docs/TESTING.md", TemplateName: "docs_testing_md.tmpl", Mode: 0o644},
		)
	}
	return basic
}
