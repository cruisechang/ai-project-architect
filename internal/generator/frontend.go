package generator

func frontendDirectories(frontendType string) []string {
	switch frontendType {
	case "next":
		return []string{"frontend/pages"}
	case "nuxt":
		return []string{"frontend/pages"}
	case "react":
		return []string{"frontend/src", "frontend/public"}
	case "vue":
		return []string{"frontend/src"}
	case "pure-typescript":
		return []string{"frontend/src"}
	case "none":
		return []string{"frontend"}
	default:
		return []string{"frontend"}
	}
}

func frontendFiles(frontendType string) []plannedFile {
	switch frontendType {
	case "next":
		return []plannedFile{
			{RelPath: "frontend/package.json", TemplateName: "frontend_next_package.tmpl", Mode: 0o644},
			{RelPath: "frontend/pages/index.tsx", TemplateName: "frontend_next_index.tmpl", Mode: 0o644},
		}
	case "nuxt":
		return []plannedFile{
			{RelPath: "frontend/package.json", TemplateName: "frontend_nuxt_package.tmpl", Mode: 0o644},
			{RelPath: "frontend/pages/index.vue", TemplateName: "frontend_nuxt_index.tmpl", Mode: 0o644},
		}
	case "react":
		return []plannedFile{
			{RelPath: "frontend/package.json", TemplateName: "frontend_react_package.tmpl", Mode: 0o644},
			{RelPath: "frontend/src/main.jsx", TemplateName: "frontend_react_main.tmpl", Mode: 0o644},
			{RelPath: "frontend/public/index.html", TemplateName: "frontend_react_index_html.tmpl", Mode: 0o644},
		}
	case "vue":
		return []plannedFile{
			{RelPath: "frontend/package.json", TemplateName: "frontend_vue_package.tmpl", Mode: 0o644},
			{RelPath: "frontend/src/main.js", TemplateName: "frontend_vue_main.tmpl", Mode: 0o644},
			{RelPath: "frontend/src/App.vue", TemplateName: "frontend_vue_app.tmpl", Mode: 0o644},
		}
	case "pure-typescript":
		return []plannedFile{
			{RelPath: "frontend/package.json", TemplateName: "frontend_ts_package.tmpl", Mode: 0o644},
			{RelPath: "frontend/tsconfig.json", TemplateName: "frontend_ts_tsconfig.tmpl", Mode: 0o644},
			{RelPath: "frontend/index.html", TemplateName: "frontend_ts_index_html.tmpl", Mode: 0o644},
			{RelPath: "frontend/src/main.ts", TemplateName: "frontend_ts_main.tmpl", Mode: 0o644},
		}
	case "none":
		return []plannedFile{
			{RelPath: "frontend/README.md", TemplateName: "frontend_none_readme.tmpl", Mode: 0o644},
		}
	default:
		return []plannedFile{
			{RelPath: "frontend/README.md", TemplateName: "frontend_none_readme.tmpl", Mode: 0o644},
		}
	}
}
