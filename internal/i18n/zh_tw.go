package i18n

func init() {
	register("zh-TW", map[string]string{
		// root
		"root.short": "從產品 idea 到可執行專案的全流程 CLI",
		"root.long": `apa（AI Project Architect）

` + iterSep + `
推薦使用流程
` + iterSep + `

  【首次建立】（在專案目錄外執行）
  apa init --idea "線上訂餐平台" --name food-platform --path ~/projects

  【進入 repo 後，由 skill 指導 agent 持續工作】
  cd ~/projects/food-platform
  apa list-skills              # 查看可用 skills
  apa iterate                  # 產出「持續迭代直到完成」AI 提示詞
                              # Codex: 明確要求 agent 使用 apa-loop + apa-implement
                              # Claude Code: /apa-loop --max-iterations 30
                              # Claude Code: /cancel-apa-loop
  make test                    # 執行測試（repo 原生 Makefile）

` + iterSep + `
apa init — 首次 bootstrap（核心指令）
` + iterSep + `

  # 互動模式（依序詢問 idea / 名稱 / 路徑 / 技術選項）
  apa init

  # 非互動模式（--idea 自動推論技術棧）
  apa init --idea "線上訂餐平台" --name food-platform --path ~/projects

  # 進階：覆蓋推論結果
  apa init --idea "線上訂餐平台" --name food-platform --backend python --agent claude-code

  # 預覽將建立的檔案（不實際寫入）
  apa init --dry-run

  產物：
    .architect/context.json   推論出的技術棧
    docs/                     PRD、SPEC、ARCHITECTURE、API、DB Schema、實作計畫
    backend/ frontend/ tests/ 可執行程式碼起始檔案
    Makefile                  make test / make lint / make build
    CLAUDE.md 或 PROMPT.md    agent 設定
    agents/ skills/           agent 與 skill 模板

` + iterSep + `
其他指令
` + iterSep + `

  apa iterate      輸出「持續迭代直到完成」AI 提示詞（任何時段皆可執行）
  apa list-skills  列出可用的 skills
  apa doctor       環境檢查
  apa version      顯示版本

執行 "apa <指令> --help" 查看各指令詳細說明。`,

		// init
		"init.short": "bootstrap 新專案（context + docs + 可執行骨架一次完成）",
		"init.long": `從 idea 一鍵建立完整的新專案。

執行步驟：
  [1/4] 從 idea 推論技術棧（可用 flag 覆蓋任一欄位）
  [2/4] 建立可執行程式碼骨架（含 Makefile / 測試設施 / agent 設定）
  [3/4] 產生設計文件（PRD / SPEC / ARCHITECTURE / API / DB Schema / 實作計畫）
  [4/4] 完成

使用場景：
  專案外首次建立  →  apa init
  進入 repo 後    →  apa iterate → 明確要求 agent 使用 apa-loop + apa-implement → make test（循環）

產物：
  .architect/context.json          推論出的技術棧
  docs/                            PRD / SPEC / ARCHITECTURE / API / DB Schema / 實作計畫
  backend/ frontend/ tests/        可執行程式碼起始檔案
  Makefile                         make test / make lint / make build
  CLAUDE.md 或 PROMPT.md           agent 設定（依 --agent）
  agents/ skills/                  agent 與 repo-local skill 模板`,

		"init.flag.idea":             "產品 idea（觸發自動推論技術棧）",
		"init.flag.idea-file":        "從檔案讀取 idea（使用 - 從 stdin 讀取）",
		"init.flag.name":             "專案名稱（留空則從 idea 自動推論）",
		"init.flag.path":             "專案建立的父目錄路徑",
		"init.flag.type":             "專案類型：web-app | ai-app | devops-tool | internal-tool | platform-service",
		"init.flag.ai-feature":       "AI 功能：none | prompt-workflow | rag | agent-system",
		"init.flag.agent":            "AI Agent 類型：codex | claude-code",
		"init.flag.backend":          "後端語言：go | python | node | none",
		"init.flag.frontend":         "前端框架：react | next | nuxt | vue | pure-typescript | none",
		"init.flag.architecture":     "架構類型：cli-tool | backend-service | frontend-app | fullstack-web-app | frontend-backend",
		"init.flag.stack":            "技術棧描述",
		"init.flag.docs":             "文件類型：basic | full",
		"init.flag.unit-test":        "加入單元測試：yes | no",
		"init.flag.integration-test": "加入整合測試：yes | no",
		"init.flag.e2e-test":         "加入 E2E 測試：yes | no",
		"init.flag.docker-compose":   "加入 Docker Compose：yes | no",
		"init.flag.skills":           "要複製的 skill 名稱，逗號分隔",
		"init.flag.skills-path":      "skills 來源目錄路徑（預設：當前 repo 的 ./skills）",
		"init.flag.description":      "專案描述",
		"init.flag.force":            "備份既有目錄後覆蓋重建",
		"init.flag.dry-run":          "預覽將建立的檔案與目錄，不實際寫入",

		// iterate
		"iterate.short": "輸出「持續迭代直到完成」AI 提示詞（任何時段皆可執行）",
		"iterate.long": `讀取當前 repo 狀態，輸出給 AI 的「持續迭代直到完成」指令。

將輸出複製並貼到 AI 後，建議搭配 ` + "`apa-loop`" + ` 使用，以強制進入交付循環。
若是 Codex 專案，請明確要求 agent 使用 ` + "`apa-loop`" + ` skill。
若是 Claude Code 專案，則可直接使用產生出的 slash command。
AI 會自動：
  1. 盤點現況（docs、tasks、測試、CI 狀態）
  2. 循環實作 → 測試 → 修復 → 更新文件
  3. 直到所有核心需求完成、測試通過

使用方式：
  apa iterate                        # 在當前 repo 目錄執行
  apa iterate --root ~/projects/foo  # 指定專案目錄
  apa iterate | pbcopy               # 複製到剪貼簿（macOS）
  apa iterate > prompt.md            # 輸出到檔案`,

		"iterate.flag.root": "專案根目錄路徑（預設為當前目錄）",

		// iterate prompt output
		"iterate.prompt.intro":                "你現在是本 repo 的主責實作 AI，請進入「持續迭代直到完成」模式，直接執行，不要只給建議。",
		"iterate.prompt.project-info":         "專案資訊",
		"iterate.prompt.root-label":           "專案根目錄：",
		"iterate.prompt.name-label":           "名稱：",
		"iterate.prompt.idea-label":           "Idea：",
		"iterate.prompt.stack-label":          "技術棧：",
		"iterate.prompt.no-context":           "（尚未執行 apa init，未找到 .architect/context.json）",
		"iterate.prompt.docs-status":          "設計文件狀態",
		"iterate.prompt.exists":               "存在",
		"iterate.prompt.missing":              "缺少",
		"iterate.prompt.no-docs":              "（執行 apa init 可建立含設計文件的完整專案）",
		"iterate.prompt.phase-warning":        "需要先重寫階段式文件",
		"iterate.prompt.phase-warning-items":  "以下既有文件不是以優先度排序的 `Phase 0`、`Phase 1`... 階段式寫法：",
		"iterate.prompt.phase-warning-action": "請先呼叫 `apa-docs` skill，將這些文件重寫成對齊的階段式內容，再繼續實作。",
		"iterate.prompt.tasks":                "任務清單",
		"iterate.prompt.workflow":             "工作方式（循環執行，直到 DONE）",
		"iterate.prompt.workflow-steps": `  1. 先盤點現況：讀取 docs、tasks、測試、CI 狀態，列出未完成項目與風險。
  2. 每輪只做 1~3 個最高優先任務（以可驗證結果為準）。
  3. 實作後必須執行必要檢查（至少含測試；有 lint 就跑 lint）。
  4. 失敗就修到通過，不可停在半成品。
  5. 更新文件與任務狀態（包含你做了什麼、剩下什麼、下一輪做什麼）。
  6. 立刻進入下一輪，不要等待我確認，除非遇到「不可自行決策」阻塞。`,
		"iterate.prompt.done": "完成定義（DONE 必須全部滿足）",
		"iterate.prompt.done-items": `  - 所有核心需求已實作完成，且與 docs 一致。
  - 全部測試通過（含新增/修正的測試）。
  - 無阻塞性錯誤、無 P0/P1 已知問題。
  - 文件已更新到可交接（README/操作方式/設計決策/待辦）。
  - 列出最終交付摘要與剩餘風險（若有）。`,
		"iterate.prompt.constraints": "執行限制",
		"iterate.prompt.constraints-items": `  - 不要刪除或回滾我未要求刪除的內容。
  - 每輪先小步提交可運行結果，再擴展下一輪。
  - 若需要重大取捨，先提出「選項 + 建議 + 影響」，其餘情況直接做。`,
		"iterate.prompt.start":             "開始執行，先輸出：",
		"iterate.prompt.start-items":       "  A. 現況盤點\n  B. 第一輪要做的 1~3 個任務\n然後直接進入實作循環，直到 DONE。",
		"iterate.prompt.start-items.phase": "  A. 使用 `apa-docs` skill，將這些文件重寫成對齊的階段式內容：%s\n  B. 文件重寫後的現況盤點\n  C. 第一輪要做的 1~3 個任務\n然後直接進入實作循環，直到 DONE。",

		// list-skills
		"list-skills.short": "列出指定目錄下可用的 skills",
		"list-skills.long": `列出 skills 目錄下第一層的所有 skill 子目錄。

Skill 是可複製到專案的功能模組目錄，搭配 "apa init --skills" 使用。
未指定 --path 時預設使用當前 repo 的 ./skills。`,

		"list-skills.flag.path": "Skills 目錄路徑（預設：./skills）",

		// doctor
		"doctor.short": "檢查執行環境（Go 版本、寫入權限、skills 路徑）",
		"doctor.long": `檢查 apa 執行所需的環境條件。

檢查項目：
  Go executable    確認 go 指令存在於 PATH
  Go version       顯示已安裝的 Go 版本
  Filesystem write 確認當前目錄有寫入權限
  Skills path      驗證 local skills 目錄是否存在（選填）

輸出格式：
  [PASS] 檢查名稱: 說明
  [FAIL] 檢查名稱: 說明`,

		"doctor.flag.skills-path": "要驗證的 skills 目錄路徑（選填，預設建議 ./skills）",
		"doctor.flag.check-write": "是否檢查當前目錄的寫入權限",

		// version
		"version.short": "顯示版本資訊（version / commit / build date）",

		// architect_helpers
		"prompt.multiline.hint": "結束請按兩次 Enter 或 Ctrl+D",

		// wizard / interactive descriptions
		"wizard.project-type.web-app":          "面向使用者的 Web 應用，偏 UI、API、auth、user flow。",
		"wizard.project-type.ai-app":           "以 prompt、model workflow、tool use、evaluation 為核心的 AI 應用。",
		"wizard.project-type.devops-tool":      "偏 CI/CD、部署、回滾、自動化、安全操作的工程工具。",
		"wizard.project-type.internal-tool":    "偏內部流程效率、維運支援、資料管理、onboarding 的工具。",
		"wizard.project-type.platform-service": "偏共享能力、基礎服務、平台化、可觀測性與安全治理的服務。",

		"wizard.ai-feature.none":            "No AI functionality",
		"wizard.ai-feature.prompt-workflow": "Simple LLM app using prompt templates",
		"wizard.ai-feature.rag":             "Retrieval Augmented Generation system",
		"wizard.ai-feature.agent-system":    "Autonomous AI agent with tool usage and planning",

		"wizard.ai-agent.codex":       "OpenAI Codex style project structure",
		"wizard.ai-agent.claude-code": "Claude Code style project structure",

		"wizard.architecture.cli-tool":          "CLI tool",
		"wizard.architecture.backend-service":   "Backend service",
		"wizard.architecture.frontend-app":      "Frontend app",
		"wizard.architecture.fullstack-web-app": "Fullstack web app",
		"wizard.architecture.frontend-backend":  "Frontend + backend",

		"wizard.fullstack.next": "React fullstack framework",
		"wizard.fullstack.nuxt": "Vue fullstack framework",
	})
}
