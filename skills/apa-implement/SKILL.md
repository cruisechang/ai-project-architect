# APA Implement Skill

## Goal

根據目前 repository 內既有文件與規格，直接推進可驗證的實作。

## Inputs

- top-level `README*.md`
- `docs/`
- `SPEC.md`
- `PRD.md`
- `API_ROUTES.md`
- `DB_SCHEMA.md`
- 任何現有設計、介面、測試或範例資料

## Steps

1. 盤點目前文件與程式碼狀態。
2. 萃取可直接實作的需求、限制與驗收標準。
3. 列出 Unknowns 與 Assumptions。
4. 先用一句話定義本輪要交付的行為。
5. 先寫或更新測試，讓測試先失敗，確認需求被正確表達。
6. 若資訊足夠，實作最小可行變更讓測試通過。
7. 重跑測試，確認從 RED 進到 GREEN。
8. 在不改變行為前提下整理程式碼，再次驗證。
9. 更新受影響文件。
10. 執行驗證指令，例如 `make test`、`go test ./...`。
11. 回報完成項、未完成項、風險、阻塞點與下一步。

## TDD Cycle

`RED -> GREEN -> REFACTOR`

- RED: 先寫失敗測試，描述目標行為。
- GREEN: 只寫讓測試通過所需的最小程式碼。
- REFACTOR: 整理命名、重複與結構，不改變行為，然後再跑測試。

## Rules

- 不可只停留在分析；文件足夠時必須直接開始 coding。
- 優先實作最小可交付版本。
- 先寫測試，再寫實作；不可先寫 production code 再補測試。
- 若新測試一開始就通過，代表測試不足或寫錯，需先修正測試。
- 不可順手加入與本輪驗證目標無關的 speculative code。
- mocks 只放在系統邊界，例如 DB、HTTP、filesystem；不要在內部模組之間過度 mock。
- 所有假設都要明確標示。
- 若文件與程式碼衝突，先以目前 repo 內最新文件為準，並標記衝突點。
- 若需求不足以安全實作，需明確列出缺口與建議補件內容。

## Output

- 實作內容摘要
- 變更檔案清單
- 測試與驗證結果
- Unknowns / Assumptions
- 風險 / 阻塞點
- 建議下一步
