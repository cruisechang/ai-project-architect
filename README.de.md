# apa - AI Project Architect

## Sprachen

| Sprache | README |
|---|---|
| English | [README.md](README.md) |
| 简体中文 | [README.zh-CN.md](README.zh-CN.md) |
| 繁體中文 | [README.zh-TW.md](README.zh-TW.md) |
| Deutsch | [README.de.md](README.de.md) |
| Français | [README.fr.md](README.fr.md) |
| Español | [README.es.md](README.es.md) |
| 日本語 | [README.ja.md](README.ja.md) |
| 한국어 | [README.ko.md](README.ko.md) |

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![CLI](https://img.shields.io/badge/Type-CLI-111111)](#befehle)
[![Skills](https://img.shields.io/badge/Repo%20Skills-apa--*-2f855a)](#commands--skills)

## `apa簡要說明`

- 可用觸發詞：`apa簡要說明`、`apa 簡要說明`、`apa說明`、`apa 說明`
- Implementierungs-Loop starten: `apa prompt` -> Ausgabe an den Agenten geben -> `apa-loop` + `apa-implement` anweisen
- Erst Dokumente iterativ schärfen: `apa prompt --docs-only` -> Ausgabe an den Agenten geben -> nur `apa-doc-review` anweisen
- Terminal-Loop-Wrapper: `bash scripts/apa-loop-setup.sh --max-iterations 30 --reviewer agent-self`
- Wrapper stoppen: `bash scripts/apa-loop-cancel.sh`

`apa` ist ein Go-CLI, das eine Produktidee schnell in einen brauchbaren Projektstart verwandelt.

Es erzeugt Projektkontext, Designdokumente, lauffähige Scaffolds und einen wiederholbaren KI-Iterationsablauf, damit der Weg von der Idee zur Umsetzung kürzer wird.

Es ist sinnvoll, wenn ein einziges Tool Projekt-Bootstrap, Dokumentenstart, Starter-Code und den Übergang in einen KI-gestützten Delivery-Loop übernehmen soll.

## Why `apa`

- Du willst, dass ein neues Projekt mit Dokumenten und Code startet, nicht nur mit leeren Ordnern
- Du willst, dass KI-Arbeit aus einem strukturierten Repo-Zustand beginnt statt aus spontanen Prompts
- Du willst repo-lokale `apa-*`-Skills und native Befehle wie `make test` als Delivery-Loop nutzen

## Installation

```bash
go build -o apa .
```

Optional:

```bash
cp ./apa /usr/local/bin/apa
```

## Flow

```text
idea
  -> apa init
  -> docs + scaffold + context
  -> apa list-skills
  -> apa prompt
  -> agent implements
  -> make test
  -> repeat
```

## Was `apa` macht

- Erstellt ein neues Projekt aus einer Produktidee in natürlicher Sprache
- Leitet einen sinnvollen Tech-Stack ab, den du per Flags überschreiben kannst
- Generiert phasenbasierte Dokumente ab `Phase 0`, darunter PRD, SPEC, ARCHITECTURE, API, DB Schema und Implementierungsplan
- Erstellt lauffähigen Startcode, Teststruktur, Makefile-Ziele und Agent-Konfiguration
- Gibt mit `apa prompt` einen strukturierten KI-Prompt für die weitere Lieferung aus
- Lässt sich natürlich mit `apa-loop` koppeln, um die Lieferung rundenbasiert voranzutreiben: Status lesen, 1-3 Aufgaben wählen, prüfen, Status aktualisieren, wiederholen

## Empfohlener Ablauf

```bash
# 1. apa bauen
go build -o apa .

# 2. Neues Projekt außerhalb des Ziel-Repos anlegen
./apa init --idea "SaaS Reporting Platform" --name report-platform --path ~/projects

# 3. Ins erzeugte Repo wechseln und iterieren
cd ~/projects/report-platform
./apa list-skills
./apa prompt
make test
```

Kernschleife:

1. Mit `apa init` das erste Projektgerüst erzeugen.
2. Phasenbasierte Dokumente ab `Phase 0` pflegen, damit Umfang, Tests, Gates und Berichte in PRD/API/SPEC ausgerichtet bleiben.
3. Standardmäßig `apa-loop` mit `apa-implement` als Delivery-Schleife verwenden.
4. `apa prompt` ausführen, den Agenten arbeiten lassen und mit `make test` prüfen.
5. Wiederholen, bis das Repo lieferbar ist.

## Status der Delivery-Schleife und `apa-loop` Verwendung

Generierte Repositories sollten `docs/IMPLEMENTATION_STATUS.md` oder `TASKS.md` fortlaufend aktualisieren.
Nutze `apa-loop` zusammen mit `apa-implement`, damit der Agent zwischen Implementierung, Tests, Fehlerbehebung und Dokumentations-Updates weiter rotiert, bis das Abschluss-Gate erfüllt ist.
`apa-loop` ist das repo-lokale Skill für den erzwungenen Rundenbetrieb: Statusdatei lesen, 1-3 überprüfbare Aufgaben wählen, Tests/Checks ausführen, Status aktualisieren und wiederholen, bis das Abschluss-Gate erfüllt ist.
Verwendung:
- Wenn du die Dokumente vor der Implementierung iterativ schärfen willst, starte mit `apa prompt --docs-only` und weise den Agenten an, nur `apa-doc-review` zu verwenden.
- Primärer Agent-Workflow (für Codex und Claude Code): `apa prompt` ausführen und den Agenten dann explizit anweisen, `apa-loop` mit `apa-implement` zu verwenden
- Optionaler Terminal-Wrapper (für Umgebungen mit generiertem Hook oder Slash-Command, z. B. Claude Code): `bash scripts/apa-loop-setup.sh --max-iterations 30 --reviewer agent-self`
- Optionaler Slash-Command: `/apa-loop --max-iterations 30 --reviewer agent-self`
- Optionaler Abbruchbefehl: `/cancel-apa-loop`
- Review policy: interactive per round. Ask which reviewer to use (`agent-self`, `apa-codex-review`, or `apa-claude-review`) before review.

## Schnelles Beispiel

```bash
./apa init \
  --idea "Internes Support-Dashboard mit Ticketsuche und KI-Zusammenfassungen" \
  --name support-hub \
  --path ~/projects \
  --backend go \
  --frontend next \
  --agent codex

cd ~/projects/support-hub
./apa list-skills
./apa prompt > prompt.md
make test
```

## Befehle

| Befehl | Zweck |
|---|---|
| `apa init` | Neues Projekt aus einer Idee erstellen |
| `apa prompt` | Strukturierten KI-Prompt für die weitere Lieferung erzeugen |
| `apa list-skills` | Verfügbare repo-lokale Skills anzeigen |
| `apa doctor` | Lokale Umgebung und Skills-Pfad prüfen |
| `apa version` | Build-Versionsinfo ausgeben |

Vollständige Optionen: `apa <command> --help`

## Commands & Skills

Aktuelle CLI-Befehle:

- `init`
- `prompt`
- `list-skills`
- `doctor`
- `version`

Aktuelle repo-lokale Skills:

- `apa-catalog`
- `apa-debug`
- `apa-devops`
- `apa-docs`
- `apa-doc-review`
- `apa-feature`
- `apa-implement`
- `apa-integration`
- `apa-loop`
- `apa-review`
- `apa-codex-review`
- `apa-claude-review`
- `apa-tdd`

## `apa init`

`apa init` ist der Kernbefehl für das erste Bootstrap.

Typischer Ablauf:

```text
[1/4] infer tech stack from idea
[2/4] generate code scaffold
[3/4] generate design docs
[4/4] done
```

Typische Nutzung:

```bash
# Interaktiv
./apa init
# Startet direkt im Wizard; die erste Frage ist Project idea.

# Nicht interaktiv
./apa init --idea "Online-Essensbestellplattform" --name food-platform --path ~/projects

# Nur Vorschau
./apa init --dry-run
```

Wichtige Flags:

| Flag | Beschreibung |
|---|---|
| `--idea` | Produktidee für die Stack-Ableitung |
| `--name` | Projektname |
| `--path` | Übergeordnetes Zielverzeichnis |
| `--type` | `cli`, `server`, `web-app-server`, `mobile-app-server`, `web-app`, `mobile-app` |
| `--agent` | `codex`, `claude-code` oder `universal` |
| `--backend` | `go`, `python`, `node` oder `none` |
| `--frontend` | `react`, `next`, `nuxt`, `vue`, `pure-typescript` oder `none` |
| `--unit-test` `--api-test` `--integration-test` `--e2e-test` | `yes` or `no` |
| `--skills` | Kommagetrennte repo-lokale Skills zum Kopieren |
| `--skills-path` | Quellverzeichnis der Skills, standardmäßig `./skills` |
| `--force` | Vorhandenes Verzeichnis sichern und neu aufbauen |
| `--dry-run` | Ausgabe planen, aber nichts schreiben |

Typische Artefakte:

```text
.architect/context.json
docs/
backend/ frontend/ tests/
Makefile
CLAUDE.md or PROMPT.md
agents/ skills/
```

## `apa prompt`

`apa prompt` liest den aktuellen Zustand des Repos und erzeugt einen strukturierten Prompt für den KI-Agenten.

Das ist vor der Implementierung, während der Entwicklung oder nach Regressionen nützlich, damit der Agent im Einklang mit Dokumenten, Aufgaben und Randbedingungen weiterarbeitet.

Der Befehl prüft außerdem, ob bestehende Dokumente prioritätsbasiert ausgerichtete Abschnitte `Phase 0`, `Phase 1`, ... verwenden. Falls nicht, warnt `apa prompt` und fordert den Agenten auf, diese Dokumente zuerst mit `apa-docs` umzuschreiben.

```bash
./apa prompt
./apa prompt --docs-only
./apa prompt --root ~/projects/report-platform
./apa prompt > prompt.md
```

## Repo-lokale Skills

Dieses Repository verwendet für repo-lokale Skills die Namensreihe `apa-*`.

Aktuelle Beispiele:

- `apa-catalog`
- `apa-debug`
- `apa-devops`
- `apa-docs`
- `apa-doc-review`
- `apa-feature`
- `apa-implement`
- `apa-integration`
- `apa-loop`
- `apa-review`
- `apa-tdd`

`apa-docs` schreibt Dokumentation in prioritätsbasierten Phasen (`Phase 0`, `Phase 1`, ...). `Phase 0` ist immer die Phase mit der höchsten Priorität. Jede Phase muss Umfang, abgestimmte PRD/API/SPEC-Inhalte, erforderliche Tests, Prüfpunkte, Abschlusskriterien, ein explizites Next-Phase-Gate und einen Phasenbericht enthalten.
`apa-doc-review` ist für eine dokumentzentrierte Review-Schleife mit dem Nutzer: pro Runde nur Doku ändern, dann auf Feedback warten, und erst nach expliziter Freigabe der Doku mit der Implementierung beginnen.

Anzeigen mit:

```bash
./apa list-skills
```

## Entwicklung und Tests

```bash
go test ./...
make test
make build TARGET=mac
make build TARGET=linux
make build TARGET=windows
make build-all
```

Beispiel für Build-Metadaten:

```bash
VERSION=v1.2.3 COMMIT=abc1234 BUILD_DATE=2026-03-19T10:00:00Z ./build.sh mac
```

## Repository-Struktur

```text
ai-project-architect/
├── apa/           # CLI-Befehle
├── internal/      # generation, planner, runtime, config, output
├── skills/        # repo-lokale apa-* Skills
├── templates/     # eingebettete Templates
├── scripts/
├── build.sh
└── Makefile
```
