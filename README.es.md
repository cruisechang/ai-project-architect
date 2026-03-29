# apa - AI Project Architect

## Idiomas

| Idioma | README |
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
[![CLI](https://img.shields.io/badge/Type-CLI-111111)](#comandos)
[![Skills](https://img.shields.io/badge/Repo%20Skills-apa--*-2f855a)](#commands--skills)

## `apa簡要說明`

- 可用觸發詞：`apa簡要說明`、`apa 簡要說明`、`apa說明`、`apa 說明`
- Iniciar el bucle de implementación: `apa prompt --reviewer agent-self` -> pega la salida en el agente -> pide `apa-loop` + `apa-implement`
- Revisar primero la documentación: `apa prompt --docs-only` -> pega la salida en el agente -> pide solo `apa-doc-review`
- Wrapper de bucle en terminal: `bash scripts/apa-loop-setup.sh --max-iterations 30 --reviewer agent-self`
- Detener el wrapper: `bash scripts/apa-loop-cancel.sh`

`apa` es una CLI en Go que convierte rápidamente una idea de producto en un punto de partida real para un proyecto.

Genera contexto del proyecto, documentos de diseño, scaffolding ejecutable y un flujo repetible de iteración con IA para acortar el paso entre idea e implementación.

Sirve bien cuando quieres que una sola herramienta cubra el bootstrap del proyecto, el arranque de documentación, el código inicial y el paso a un ciclo de entrega asistido por IA.

## Why `apa`

- Quieres que un proyecto nuevo empiece con documentación y código, no solo con carpetas vacías
- Quieres que la IA trabaje desde un estado de repo estructurado y no desde prompts improvisados
- Quieres usar skills repo-local `apa-*` y comandos nativos como `make test` para avanzar

## Installation

```bash
go build -o apa .
```

Opcional:

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

## Qué hace `apa`

- Inicializa un proyecto nuevo a partir de una idea de producto en lenguaje natural
- Infiere una stack técnica razonable y permite sobrescribirla con flags
- Genera documentación por fases desde `Phase 0`, incluyendo PRD, SPEC, ARCHITECTURE, API, DB Schema y plan de implementación
- Crea código inicial ejecutable, estructura de tests, objetivos de Makefile y configuración de agente
- Produce con `apa prompt` un prompt estructurado para que un agente de IA siga entregando trabajo
- Se combina de forma natural con `apa-loop` para una entrega por rondas: leer estado, elegir 1-3 tareas, verificar, actualizar estado y repetir

## Flujo recomendado

```bash
# 1. Compilar apa
go build -o apa .

# 2. Crear un proyecto nuevo fuera del repo objetivo
./apa init --idea "Plataforma SaaS de reporting" --name report-platform --path ~/projects

# 3. Entrar en el repo generado e iterar
cd ~/projects/report-platform
./apa list-skills
./apa prompt --reviewer agent-self
make test
```

Bucle principal:

1. Ejecutar `apa init` una vez para crear la primera versión del proyecto.
2. Mantener documentación por fases desde `Phase 0` para alinear alcance, pruebas, gates e informes entre PRD/API/SPEC.
3. Usar `apa-loop` con `apa-implement` como bucle de entrega por defecto.
4. Ejecutar `apa prompt --reviewer agent-self`, dejar trabajar al agente y validar con `make test`.
5. Repetir hasta que el repo esté listo para entregar.

## Estado del ciclo de entrega y uso de `apa-loop`

Los repos generados deben mantener `docs/IMPLEMENTATION_STATUS.md` o `TASKS.md` actualizados.
Usa `apa-loop` con `apa-implement` para que el agente siga rotando entre implementación, pruebas, correcciones y actualizaciones de documentación hasta que se cumpla el gate de finalización.
`apa-loop` es el skill repo-local que fuerza el bucle de entrega por rondas: leer el archivo de estado, elegir 1-3 tareas verificables, ejecutar pruebas o checks, actualizar el estado y repetir hasta cumplir el gate de finalización.
Uso:
- Si quieres iterar primero sobre la documentación antes de implementar, ejecuta `apa prompt --docs-only` y pide al agente que use solo `apa-doc-review`.
- Flujo principal por agente (común para Codex y Claude Code): ejecuta `apa prompt --reviewer agent-self`, o ejecuta `apa prompt` y elige el reviewer de forma interactiva, y luego pide explícitamente al agente que use `apa-loop` con `apa-implement`
- Wrapper opcional en terminal (para entornos con hook o slash command generado, como Claude Code): `bash scripts/apa-loop-setup.sh --max-iterations 30 --reviewer agent-self`
- Slash command opcional: `/apa-loop --max-iterations 30 --reviewer agent-self`
- Comando opcional de cancelación: `/cancel-apa-loop`
- Review policy: especifica el reviewer una sola vez con `--reviewer` o en la primera instrucción y reutilízalo en las siguientes rondas; solo cámbialo si lo pides explícitamente (`agent-self`, `apa-codex-review`, `apa-claude-review`).

## Ejemplo rápido

```bash
./apa init \
  --idea "Panel interno de soporte con búsqueda de tickets y resúmenes con IA" \
  --name support-hub \
  --path ~/projects \
  --backend go \
  --frontend next \
  --agent codex

cd ~/projects/support-hub
./apa list-skills
./apa prompt --reviewer agent-self > prompt.md
make test
```

## Comandos

| Comando | Propósito |
|---|---|
| `apa init` | Crear un proyecto nuevo desde una idea |
| `apa prompt` | Generar el prompt estructurado de IA para seguir entregando |
| `apa list-skills` | Mostrar los skills repo-local disponibles |
| `apa doctor` | Revisar el entorno local y la ruta de skills |
| `apa version` | Imprimir la versión del build |

Opciones completas: `apa <command> --help`

## Commands & Skills

Comandos CLI actuales:

- `init`
- `prompt`
- `list-skills`
- `doctor`
- `version`

Skills repo-local actuales:

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

`apa init` es el comando principal para el bootstrap inicial.

Flujo típico:

```text
[1/4] infer tech stack from idea
[2/4] generate code scaffold
[3/4] generate design docs
[4/4] done
```

Usos comunes:

```bash
# Modo interactivo
./apa init
# Entra directo al Wizard y la primera pregunta es Project idea.

# Modo no interactivo
./apa init --idea "Plataforma de pedidos de comida online" --name food-platform --path ~/projects

# Solo vista previa
./apa init --dry-run
```

Flags comunes:

| Flag | Descripción |
|---|---|
| `--idea` | Idea de producto usada para inferir la stack |
| `--name` | Nombre del proyecto |
| `--path` | Directorio padre donde se crea el proyecto |
| `--type` | `cli`, `server`, `web-app-server`, `mobile-app-server`, `web-app`, `mobile-app` |
| `--agent` | `codex`, `claude-code` o `universal` |
| `--backend` | `go`, `python`, `node` o `none` |
| `--frontend` | `react`, `next`, `nuxt`, `vue`, `pure-typescript` o `none` |
| `--unit-test` `--api-test` `--integration-test` `--e2e-test` | `yes` or `no` |
| `--skills` | Skills repo-local a copiar, separados por comas |
| `--skills-path` | Directorio fuente de skills, por defecto `./skills` |
| `--force` | Respaldar un directorio existente y reconstruir |
| `--dry-run` | Mostrar el plan sin escribir archivos |

Salida típica generada:

```text
.architect/context.json
docs/
backend/ frontend/ tests/
Makefile
CLAUDE.md or PROMPT.md
agents/ skills/
```

## `apa prompt`

`apa prompt` lee el estado actual del repositorio y genera un prompt estructurado para el agente de IA.

Es útil antes de implementar, durante el desarrollo o después de una regresión, para que el agente siga trabajando alineado con los documentos, tareas y restricciones del repo.

También comprueba si los documentos existentes usan secciones alineadas por prioridad `Phase 0`, `Phase 1`, ... Si no es así, `apa prompt` avisa al usuario y le indica al agente que los reescriba con `apa-docs` antes de continuar.

```bash
./apa prompt
./apa prompt --reviewer agent-self
./apa prompt --docs-only
./apa prompt --root ~/projects/report-platform
./apa prompt > prompt.md
```

## Skills repo-local

Este repositorio usa la serie `apa-*` para los skills repo-local.

Ejemplos actuales:

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

`apa-docs` redacta la documentación por fases alineadas por prioridad (`Phase 0`, `Phase 1`, ...). `Phase 0` siempre es la fase de mayor prioridad. Cada fase debe definir alcance, contenido alineado de PRD/API/SPEC, pruebas requeridas, puntos de verificación, criterios de finalización, un gate explícito para la siguiente fase e informe de cierre de fase.
`apa-doc-review` sirve para una revisión documental iterativa con el usuario: solo cambia documentos, se detiene tras cada revisión para recibir feedback, y no inicia implementación hasta que la documentación quede aprobada explícitamente.

Listarlos con:

```bash
./apa list-skills
```

## Desarrollo y pruebas

```bash
go test ./...
make test
make build TARGET=mac
make build TARGET=linux
make build TARGET=windows
make build-all
```

Ejemplo de metadata de build:

```bash
VERSION=v1.2.3 COMMIT=abc1234 BUILD_DATE=2026-03-19T10:00:00Z ./build.sh mac
```

## Estructura del repositorio

```text
ai-project-architect/
├── apa/           # comandos CLI
├── internal/      # generation, planner, runtime, config, output
├── skills/        # skills repo-local apa-*
├── templates/     # plantillas embebidas
├── scripts/
├── build.sh
└── Makefile
```
