package generator

func backendDirectories(backendType string) []string {
	switch backendType {
	case "go":
		return []string{"backend/cmd/server", "backend/internal", "backend/pkg"}
	case "python":
		return []string{"backend"}
	case "node":
		return []string{"backend/src"}
	case "none":
		return []string{"backend"}
	default:
		return []string{"backend"}
	}
}

func backendFiles(backendType string) []plannedFile {
	switch backendType {
	case "go":
		return []plannedFile{
			{RelPath: "backend/go.mod", TemplateName: "backend_go_mod.tmpl", Mode: 0o644},
			{RelPath: "backend/cmd/server/main.go", TemplateName: "backend_go_main.tmpl", Mode: 0o644},
			{RelPath: "backend/cmd/server/main_test.go", TemplateName: "backend_go_main_test.tmpl", Mode: 0o644},
			{RelPath: "backend/internal/.gitkeep", RawContent: "", Mode: 0o644},
			{RelPath: "backend/pkg/.gitkeep", RawContent: "", Mode: 0o644},
		}
	case "python":
		return []plannedFile{
			{RelPath: "backend/main.py", TemplateName: "backend_python_main.tmpl", Mode: 0o644},
			{RelPath: "backend/test_main.py", TemplateName: "backend_python_test_main.tmpl", Mode: 0o644},
			{RelPath: "backend/requirements.txt", TemplateName: "backend_python_requirements.tmpl", Mode: 0o644},
			{RelPath: "backend/README.md", TemplateName: "backend_python_readme.tmpl", Mode: 0o644},
		}
	case "node":
		return []plannedFile{
			{RelPath: "backend/package.json", TemplateName: "backend_node_package.tmpl", Mode: 0o644},
			{RelPath: "backend/src/index.js", TemplateName: "backend_node_index.tmpl", Mode: 0o644},
			{RelPath: "backend/src/index.test.js", TemplateName: "backend_node_index_test.tmpl", Mode: 0o644},
			{RelPath: "backend/README.md", TemplateName: "backend_node_readme.tmpl", Mode: 0o644},
		}
	case "none":
		return []plannedFile{
			{RelPath: "backend/README.md", TemplateName: "backend_none_readme.tmpl", Mode: 0o644},
		}
	default:
		return []plannedFile{
			{RelPath: "backend/README.md", TemplateName: "backend_none_readme.tmpl", Mode: 0o644},
		}
	}
}
