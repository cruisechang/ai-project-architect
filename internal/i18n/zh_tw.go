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
  apa prompt                   # 產出「持續迭代直到完成」AI 提示詞
                              # 明確要求 agent 使用 apa-loop + apa-implement
  apa prompt --docs-only       # 在實作前產出文件審閱提示詞
                              # 明確要求 agent 只使用 apa-doc-review
                              # 若環境有 hook / slash wrapper，也可選擇：
                              # /apa-loop --max-iterations 30 --reviewer agent-self
                              # /cancel-apa-loop
  make test                    # 執行測試（repo 原生 Makefile）

` + iterSep + `
apa init — 首次 bootstrap（核心指令）
` + iterSep + `

  # 互動模式（依序詢問 idea / 名稱 / 路徑 / 技術選項）
  apa init

  # 非互動模式（--idea 自動推論技術棧）
  apa init --idea "線上訂餐平台" --name food-platform --path ~/projects

  # 進階：覆蓋推論結果
  apa init --idea "線上訂餐平台" --name food-platform --backend python --agent universal

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

  apa prompt      輸出「持續迭代直到完成」AI 提示詞（任何時段皆可執行）
  apa prompt --docs-only  輸出文件審閱提示詞（實作前）
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
  進入 repo 後    →  apa prompt → 明確要求 agent 使用 apa-loop + apa-implement → make test（循環）

產物：
  .architect/context.json          推論出的技術棧
  docs/                            PRD / SPEC / ARCHITECTURE / API / DB Schema / 實作計畫
  backend/ frontend/ tests/        可執行程式碼起始檔案
  Makefile                         make test / make lint / make build
  AGENTS.md + PROMPT.md / CLAUDE.md agent 設定（依 --agent）
  agents/ skills/                  agent 與 repo-local skill 模板`,

		"init.flag.idea":             "產品 idea（觸發自動推論技術棧）",
		"init.flag.idea-file":        "從檔案讀取 idea（使用 - 從 stdin 讀取）",
		"init.flag.name":             "專案名稱（留空則從 idea 自動推論）",
		"init.flag.path":             "專案建立的父目錄路徑",
		"init.flag.type":             "類型：cli | server | web-app-server | mobile-app-server | web-app | mobile-app",
		"init.flag.ai-feature":       "AI 功能：none | prompt-workflow | rag | agent-system",
		"init.flag.agent":            "AI Agent 類型：codex | claude-code | universal",
		"init.flag.backend":          "後端語言：go | python | node | none",
		"init.flag.frontend":         "前端框架：react | next | nuxt | vue | pure-typescript | none",
		"init.flag.stack":            "技術棧描述",
		"init.flag.docs":             "文件類型：basic | full",
		"init.flag.unit-test":        "加入單元測試：yes | no",
		"init.flag.api-test":         "加入 API 測試：yes | no",
		"init.flag.integration-test": "加入整合測試：yes | no",
		"init.flag.e2e-test":         "加入 E2E 測試：yes | no",
		"init.flag.docker-compose":   "加入 Docker Compose：yes | no",
		"init.flag.skills":           "要複製的 skill 名稱，逗號分隔",
		"init.flag.skills-path":      "skills 來源目錄路徑（預設：當前 repo 的 ./skills）",
		"init.flag.description":      "專案描述",
		"init.flag.force":            "備份既有目錄後覆蓋重建",
		"init.flag.dry-run":          "預覽將建立的檔案與目錄，不實際寫入",

		// prompt
		"prompt.short": "輸出「持續迭代直到完成」AI 提示詞（任何時段皆可執行）",
		"prompt.long": `讀取當前 repo 狀態，輸出給 AI 的「持續迭代直到完成」指令。

將輸出複製並貼到 AI 後，統一以 ` + "`apa-loop`" + ` + ` + "`apa-implement`" + ` 作為主要交付循環。
不論是 Codex 或 Claude Code，主用法都是明確要求 agent 使用這兩個 skills。
若執行環境另外提供 slash command 或 hook，將其視為同一套流程的可選包裝。
AI 會自動：
  1. 盤點現況（docs、tasks、測試、CI 狀態）
  2. 依 ` + "`apa-loop`" + ` 循環：選 1~3 項任務 → RED → GREEN → REFACTOR → 驗證 → 更新文件與狀態
  3. 直到 completion gate 滿足，才能輸出 ` + "`<promise>COMPLETE</promise>`" + `

使用方式：
  apa prompt                        # 在當前 repo 目錄執行
  apa prompt --docs-only            # 在實作前產出文件審閱提示詞
  apa prompt --reviewer agent-self  # 產出已指定 reviewer 的實作提示詞
  apa prompt --root ~/projects/foo  # 指定專案目錄
  apa prompt | pbcopy               # 複製到剪貼簿（macOS）
  apa prompt > prompt.md            # 輸出到檔案`,

		"prompt.flag.root":        "專案根目錄路徑（預設為當前目錄）",
		"prompt.flag.docs-only":   "產出只使用 `apa-doc-review`、且禁止開始實作的文件審閱提示詞",
		"prompt.flag.reviewer":    "在實作提示詞中持續沿用的 reviewer：agent-self | apa-codex-review | apa-claude-review",
		"prompt.mode.label":       "提示詞模式？（implementation/docs-only）",
		"prompt.mode.default":     "implementation",
		"prompt.mode.invalid":     "無效的提示詞模式 %q（請使用 implementation 或 docs-only）",
		"prompt.reviewer.label":   "Reviewer？（agent-self/apa-codex-review/apa-claude-review）",
		"prompt.reviewer.default": "agent-self",
		"prompt.reviewer.invalid": "無效的 reviewer %q（請使用 agent-self、apa-codex-review、或 apa-claude-review）",

		// prompt output
		"prompt.output.intro":                "你現在是本 repo 的主責實作 AI，請進入「持續迭代直到完成」模式。使用 `apa-loop` 與 `apa-implement` skills，直接執行，不要只給建議。",
		"prompt.output.project-info":         "專案資訊",
		"prompt.output.root-label":           "專案根目錄：",
		"prompt.output.reviewer-label":       "Reviewer：",
		"prompt.output.name-label":           "名稱：",
		"prompt.output.idea-label":           "Idea：",
		"prompt.output.stack-label":          "技術棧：",
		"prompt.output.no-context":           "（尚未執行 apa init，未找到 .architect/context.json）",
		"prompt.output.docs-status":          "設計文件狀態",
		"prompt.output.exists":               "存在",
		"prompt.output.missing":              "缺少",
		"prompt.output.no-docs":              "（執行 apa init 可建立含設計文件的完整專案）",
		"prompt.output.phase-warning":        "需要先重寫階段式文件",
		"prompt.output.phase-warning-items":  "以下既有文件不是以優先度排序的 `Phase 0`、`Phase 1`... 階段式寫法：",
		"prompt.output.phase-warning-action": "請先呼叫 `apa-docs` skill，將這些文件重寫成對齊的階段式內容，再繼續實作。",
		"prompt.output.tasks":                "任務清單",
		"prompt.output.workflow":             "工作方式（循環執行，直到 DONE）",
		"prompt.output.workflow-steps": `  1. 先盤點現況：讀取 docs、tasks、測試、CI 狀態，列出未完成項目與風險。
  2. 以 ` + "`apa-loop`" + ` 與 ` + "`apa-implement`" + ` 作為主要工作流程。
  3. 每輪只做 1~3 個最高優先任務，以可驗證結果為準。
  4. 嚴格遵守 RED → GREEN → REFACTOR，先寫失敗測試或其他可執行驗證，再補實作。
  5. 實作後必須執行必要檢查（至少含測試；有 lint 就跑 lint）。
  6. 失敗就修到通過，不可停在半成品。
  7. 每輪 review 一律使用 ` + "`%s`" + `，除非我明確要求切換 reviewer。
  8. 更新 repo-local 狀態檔（` + "`docs/IMPLEMENTATION_STATUS.md`" + ` 或 ` + "`TASKS.md`" + `），記錄已完成、進行中、失敗檢查、blockers、assumptions 與下一輪 1~3 項任務。
  9. 除非遇到無法安全自行決策的阻塞，否則立刻進入下一輪。`,
		"prompt.output.done": "完成定義（DONE 必須全部滿足）",
		"prompt.output.done-items": `  - 所有已文件化的 P0 / 核心需求皆已完成，且與 docs 一致。
  - 所有已文件化的核心流程皆可執行。
  - 測試 / build / lint / 必要檢查全部通過。
  - 無阻塞性錯誤，且沒有未處理的高嚴重度問題。
  - 文件與 repo-local 狀態已更新到可交接。
  - 只有在 completion gate 完整滿足後，才能輸出 ` + "`<promise>COMPLETE</promise>`" + `。`,
		"prompt.output.constraints": "執行限制",
		"prompt.output.constraints-items": `  - 不要刪除或回滾我未要求刪除的內容。
  - 不可做完單一小任務或部分功能就停止。
  - 文件足夠時直接做；只有遇到無法安全假設或不可逆風險時才提問。
  - 若需要重大取捨，先提出「選項 + 建議 + 影響」，其餘情況直接做。`,
		"prompt.output.start":             "開始執行，先輸出：",
		"prompt.output.start-items":       "  A. 現況盤點\n  B. 第一輪要做的 1~3 個任務\n  C. 確認本次 loop 使用 reviewer `%s`\n然後直接進入 apa-loop 實作循環，直到 completion gate 完整滿足。",
		"prompt.output.start-items.phase": "  A. 使用 `apa-docs` skill，將這些文件重寫成對齊的階段式內容：%s\n  B. 文件重寫後的現況盤點\n  C. 第一輪要做的 1~3 個任務\n  D. 確認本次 loop 使用 reviewer `%s`\n然後直接進入 apa-loop 實作循環，直到 completion gate 完整滿足。",
		"prompt.output.docs-only.intro":   "你現在是本 repo 的文件審閱 AI。只使用 `apa-doc-review` skill。不要實作程式碼，直接修改文件，並在每輪修改後停下等待回饋。",
		"prompt.output.docs-only.workflow-steps": `  1. 先做文件盤點：讀取 README、PRD、SPEC、API、DB schema、architecture 與 implementation plan。
  2. 以 ` + "`apa-doc-review`" + ` 作為主要工作流程。
  3. 每一輪只做最小但有價值的文件修訂，優先改善清晰度、一致性與範圍控制。
  4. 只更新本輪需要的文件，不要修改實作程式碼。
  5. 每輪修改後，說明這次改了什麼、還有哪些不清楚、以及下一輪建議聚焦的文件。
  6. 每輪都必須停下等待我的回饋。
  7. 在我明確說出 ` + "`docs approved`" + ` 之前，不得開始 ` + "`apa-loop`" + ` 或 ` + "`apa-implement`" + `。`,
		"prompt.output.docs-only.done-items": `  - PRD、SPEC、API、DB schema、architecture、README 與 implementation plan 在需要的地方已彼此對齊。
  - 所有 assumptions、未解決取捨、與延後決策都已明確寫出，沒有被隱藏。
  - 文件已清楚到足以交接實作。
  - 只有在我明確說出 ` + "`docs approved`" + ` 後，你才能詢問是否切換到 ` + "`apa-loop`" + ` + ` + "`apa-implement`" + `。`,
		"prompt.output.docs-only.constraints-items": `  - 不要撰寫或修改實作程式碼。
  - 不要啟動 ` + "`apa-loop`" + `。
  - 不要使用 ` + "`apa-implement`" + `。
  - 每輪修訂要小而可審閱。
  - 若有不確定之處，請標記為 assumption 或 open question，不要假裝已確定。`,
		"prompt.output.docs-only.start-items":       "  A. 目前文件的問題與不一致處\n  B. 本輪最小但有價值的文件修改\n然後修改文件並停下等待回饋。",
		"prompt.output.docs-only.start-items.phase": "  A. 先使用 `apa-docs` skill，把這些文件重寫成對齊的階段式內容：%s\n  B. 文件重寫後，目前文件的問題與不一致處\n  C. 本輪最小但有價值的文件修改\n然後修改文件並停下等待回饋。",

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
		"wizard.project-type.frontend-only": "Frontend only：以前端為主，不含後端服務。",
		"wizard.project-type.backend-only":  "Backend only：以 API／服務／CLI 為主。",
		"wizard.project-type.full-stack":    "Full stack：前後端同專案。",

		"wizard.ai-feature.none":            "No AI functionality",
		"wizard.ai-feature.prompt-workflow": "Simple LLM app using prompt templates",
		"wizard.ai-feature.rag":             "Retrieval Augmented Generation system",
		"wizard.ai-feature.agent-system":    "Autonomous AI agent with tool usage and planning",

		"wizard.ai-agent.codex":       "OpenAI Codex style project structure",
		"wizard.ai-agent.claude-code": "Claude Code style project structure",
		"wizard.ai-agent.universal":   "共用同一套專案核心，同時產生 Codex 與 Claude Code 包裝",

		"wizard.architecture.cli":               "CLI",
		"wizard.architecture.server":            "Server",
		"wizard.architecture.web-app-server":    "Web app + server",
		"wizard.architecture.mobile-app-server": "Mobile app + server",
		"wizard.architecture.web-app":           "Web app",
		"wizard.architecture.mobile-app":        "Mobile app",

		"wizard.fullstack.next": "React fullstack framework",
		"wizard.fullstack.nuxt": "Vue fullstack framework",
	})
}
