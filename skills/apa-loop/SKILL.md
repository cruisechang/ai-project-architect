# APA Loop Skill

## Activation

This is a repo-local skill for Codex-style agents.
Codex projects do not generate a `/apa-loop` slash command or a `.claude/settings.json` Stop hook.

Activate it by explicitly telling the agent to use the `apa-loop` skill together with `apa-implement`, usually right after pasting the output of `apa iterate`.

Example prompt:
```
Use the `apa-loop` and `apa-implement` skills. Read `docs/IMPLEMENTATION_STATUS.md`, pick the next 1-3 verifiable tasks, follow RED -> GREEN -> REFACTOR, run validation, update the status file, and continue until the completion gate is met.
```

When the Completion Gate is fully met, output exactly:
```
<promise>COMPLETE</promise>
```

---

## Goal

依據目前 repository 內現有文件與程式碼狀態，持續循環推進實作、測試與修正，直到 app 達到可交付狀態。

## Inputs

- `README*.md`
- `docs/`
- `SPEC.md`
- `PRD.md`
- `API_ROUTES.md`
- `DB_SCHEMA.md`
- 現有程式碼
- 現有測試
- build / test / lint 指令

## Loop

1. 讀取目前文件與程式碼，確認最高優先需求。
2. 讀取目前的進度與任務狀態檔，確認尚未完成、進行中、失敗與阻塞項。
3. 選擇 1-3 個當前最重要且可驗證的工作項。
4. 先為每個工作項定義目標行為與驗收方式。
5. 先寫或更新測試，刻意讓測試失敗。
6. 以主流程與可執行性為優先，做最小可交付變更讓測試轉綠。
7. 執行測試與驗證。
8. 修正失敗、回歸或不一致之處。
9. 在保持行為不變的前提下重構，然後再次驗證。
10. 更新必要文件與進度狀態檔。
11. 根據最新結果決定下一輪工作。
12. 重複以上流程，直到停止條件成立。

## Embedded TDD Policy

每一輪都要內建 `RED -> GREEN -> REFACTOR`：

- RED: 先用 failing test 描述目標行為。
- GREEN: 只補足讓測試通過的最小實作。
- REFACTOR: 清理結構與命名，再次跑測試確認沒有行為回歸。

## Execution State

- 維護 repo-local 狀態檔，例如 `docs/IMPLEMENTATION_STATUS.md` 或 `TASKS.md`。
- 每輪開始前都要先讀取狀態檔。
- 每輪結束後都要更新狀態檔，至少包含：
  - 已完成項目
  - 進行中項目
  - 下一輪 1-3 個工作項
  - failing tests / failing checks
  - blockers
  - assumptions / unknowns
- 若 repo 尚未有狀態檔，第一輪就建立一份。

## Stop Conditions

- 主要功能完成
- 核心使用流程可運作
- 相關測試通過
- 高優先 bug 已處理
- 文件與實作一致
- 沒有明確且可立即推進的高價值工作項

## Completion Gate

不可因為單一功能完成、單一測試通過、或只完成部分畫面就停止。只有在以下條件同時成立時才可結束：

- 所有已文件化的 P0 功能都已實作
- 所有已文件化的核心流程都可執行
- 相關測試、build、lint 或驗證檢查已通過
- 沒有已知高嚴重度缺陷仍未處理
- README、docs 與實作行為一致
- 狀態檔內沒有更高優先且可立即推進的未完成工作

## Decision Policy

**預設：有文件就直接做，不要先問。**

- 若現有文件足以支持實作，直接執行，不要先詢問。
- 若只缺少局部細節，採用最保守且可回退的假設繼續，並明確記錄 Assumptions。
- 不要因為一般實作選型、命名、檔案拆分、測試寫法、UI 微調而提問，應自行做出合理決策並持續推進。
- 除非符合以下條件，否則不得中斷向使用者提問：
  - 規格彼此直接衝突，且任何假設都有破壞性風險
  - 缺少無法推測的關鍵商業規則
  - 需要執行破壞性或不可逆的操作（刪資料、drop table、reset migration）
  - 需要外部憑證、金鑰或部署權限
  - 操作會影響生產環境或付費外部資源

若必須提問，需先整理：
1. 已知事實
2. 缺少的關鍵資訊
3. 若採用預設假設會產生的風險
4. 建議的預設做法

## Default Autonomy

在不涉及 secrets、外部付費資源、部署到正式環境、刪除資料、修改生產設定的前提下，允許自主執行：

- 讀取任何 repo 內檔案
- 修改程式碼
- 新增或修改測試
- 執行本地測試、build、lint
- 更新 repo 內文件

## Risk Tiers

| 風險層級 | 操作類型 | 處理方式 |
|---|---|---|
| **低** | 讀檔、改碼、補測試、跑本地測試、更新文件 | 直接執行，不需確認 |
| **中** | 依 assumption 實作非關鍵細節、新增 migration、調整 config | 採保守方式執行，明確標示 Assumption |
| **高** | 刪資料、改 production 設定、使用 secrets、付費 API、破壞性指令 | 停止，說明風險，等待明確授權 |

## Priority Order

除非被阻塞，必須依照以下順序推進：

1. 先讓 app 可 build、可啟動、可執行。
2. 完成核心主流程。
3. 為核心主流程補齊或修正測試。
4. 修正 blocking bugs、回歸與高優先缺陷。
5. 補齊次要功能。
6. 對齊 README、docs 與實際行為。
7. 最後才做非必要優化與 polish。

## Rules

- 不可只停留在分析或規劃。
- 不可完成單一小功能後就結束整體任務。
- 每輪都必須有可驗證產出。
- 每輪都必須執行測試、檢查或其他驗證。
- 每輪都必須先有測試或其他可執行驗證，再補實作。
- 每輪都必須更新 repo-local 狀態檔，而不是只在回覆中描述進度。
- 若文件不足但仍可安全假設，明確標示 Assumptions 並繼續。
- 若缺口已阻止安全實作，才停止並列出阻塞。
- 優先處理主流程，再處理次要功能與優化。
- 若前一輪驗證失敗，優先修復失敗項，不可跳過後直接做新功能。
- 若測試一開始就通過，需先確認測試是否真的覆蓋目標行為。

## Output Per Round

- 本輪完成項
- 測試／驗證結果
- 剩餘問題
- Assumptions / Unknowns
- 下一輪目標
