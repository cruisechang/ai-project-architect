# apa - AI Project Architect

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![CLI](https://img.shields.io/badge/Type-CLI-111111)](#commandes)
[![Skills](https://img.shields.io/badge/Repo%20Skills-apa--*-2f855a)](#commands--skills)

`apa` est une CLI Go qui transforme rapidement une idée produit en point de départ concret pour un projet.

Elle génère le contexte du projet, les documents de conception, un scaffold exécutable et un flux d'itération IA reproductible pour accélérer le passage de l'idée à l'implémentation.

L'outil est utile si tu veux centraliser le bootstrap du projet, l'amorçage de la documentation, la génération du code initial et la transition vers une boucle de livraison assistée par IA.

## Why `apa`

- Tu veux qu'un nouveau projet démarre avec de la documentation et du code, pas seulement des dossiers vides
- Tu veux que l'IA parte d'un état de repo structuré plutôt que de prompts improvisés
- Tu veux t'appuyer sur des skills repo-local `apa-*` et des commandes natives comme `make test`

## Installation

```bash
go build -o apa .
```

Optionnel :

```bash
cp ./apa /usr/local/bin/apa
```

## Flow

```text
idea
  -> apa init
  -> docs + scaffold + context
  -> apa list-skills
  -> apa iterate
  -> agent implements
  -> make test
  -> repeat
```

## Langues

| Langue | README |
|---|---|
| English | [README.md](README.md) |
| 简体中文 | [README.zh-CN.md](README.zh-CN.md) |
| 繁體中文 | [README.zh-TW.md](README.zh-TW.md) |
| Deutsch | [README.de.md](README.de.md) |
| Français | [README.fr.md](README.fr.md) |
| Español | [README.es.md](README.es.md) |
| 日本語 | [README.ja.md](README.ja.md) |
| 한국어 | [README.ko.md](README.ko.md) |

## Ce que fait `apa`

- Initialise un nouveau projet à partir d'une idée produit en langage naturel
- Déduit une stack technique pragmatique, modifiable ensuite via des flags
- Génère PRD, SPEC, ARCHITECTURE, API, DB Schema et plan d'implémentation
- Crée un starter runnable, les tests, un Makefile et la configuration d'agent
- Produit avec `apa iterate` un prompt structuré pour continuer la livraison avec un agent IA

## Workflow recommandé

```bash
# 1. Compiler apa
go build -o apa .

# 2. Créer un nouveau projet en dehors du repo cible
./apa init --idea "Plateforme SaaS de reporting" --name report-platform --path ~/projects

# 3. Entrer dans le repo généré et itérer
cd ~/projects/report-platform
./apa list-skills
./apa iterate
make test
```

Boucle principale :

1. Exécuter `apa init` une fois pour créer la première version du projet.
2. Utiliser les skills repo-local `apa-*` pour guider l'implémentation.
3. Lancer `apa iterate`, laisser l'agent travailler, puis valider avec `make test`.
4. Répéter jusqu'à obtenir un dépôt livrable.

## Exemple rapide

```bash
./apa init \
  --idea "Tableau de bord support interne avec recherche de tickets et résumés IA" \
  --name support-hub \
  --path ~/projects \
  --backend go \
  --frontend next \
  --agent codex

cd ~/projects/support-hub
./apa list-skills
./apa iterate > prompt.md
make test
```

## Commandes

| Commande | Rôle |
|---|---|
| `apa init` | Créer un nouveau projet à partir d'une idée |
| `apa iterate` | Générer le prompt IA structuré pour la suite |
| `apa list-skills` | Afficher les skills repo-local disponibles |
| `apa doctor` | Vérifier l'environnement local et le chemin des skills |
| `apa version` | Afficher les informations de version du build |

Options complètes : `apa <command> --help`

## Commands & Skills

Commandes CLI actuelles :

- `init`
- `iterate`
- `list-skills`
- `doctor`
- `version`

Skills repo-local actuels :

- `apa-catalog`
- `apa-debug`
- `apa-devops`
- `apa-docs`
- `apa-feature`
- `apa-integration`
- `apa-review`
- `apa-tdd`

## `apa init`

`apa init` est la commande principale pour le bootstrap initial.

Flux typique :

```text
[1/4] infer tech stack from idea
[2/4] generate code scaffold
[3/4] generate design docs
[4/4] done
```

Usages courants :

```bash
# Mode interactif
./apa init

# Mode non interactif
./apa init --idea "Plateforme de commande de repas en ligne" --name food-platform --path ~/projects

# Aperçu uniquement
./apa init --dry-run
```

Flags utiles :

| Flag | Description |
|---|---|
| `--idea` | Idée produit utilisée pour déduire la stack |
| `--name` | Nom du projet |
| `--path` | Répertoire parent de création |
| `--agent` | `codex` ou `claude-code` |
| `--backend` | `go`, `python`, `node` ou `none` |
| `--frontend` | `react`, `next`, `nuxt`, `vue`, `pure-typescript` ou `none` |
| `--skills` | Skills repo-local à copier, séparés par des virgules |
| `--skills-path` | Répertoire source des skills, par défaut `./skills` |
| `--force` | Sauvegarder un répertoire existant puis reconstruire |
| `--dry-run` | Afficher la sortie prévue sans écrire |

Sortie générée typique :

```text
.architect/context.json
docs/
backend/ frontend/ tests/
Makefile
CLAUDE.md or PROMPT.md
agents/ skills/
```

## `apa iterate`

`apa iterate` lit l'état actuel du dépôt et imprime un prompt structuré pour l'agent IA.

La commande est utile avant l'implémentation, pendant le développement ou après une régression, afin que l'agent continue à travailler en tenant compte des documents, des tâches et des contraintes du repo.

```bash
./apa iterate
./apa iterate --root ~/projects/report-platform
./apa iterate > prompt.md
```

## Skills repo-local

Ce dépôt utilise la série de noms `apa-*` pour les skills repo-local.

Exemples actuels :

- `apa-catalog`
- `apa-debug`
- `apa-devops`
- `apa-docs`
- `apa-feature`
- `apa-integration`
- `apa-review`
- `apa-tdd`

Liste avec :

```bash
./apa list-skills
```

## Développement et tests

```bash
go test ./...
make test
make build TARGET=mac
make build TARGET=linux
make build TARGET=windows
make build-all
```

Exemple de métadonnées de build :

```bash
VERSION=v1.2.3 COMMIT=abc1234 BUILD_DATE=2026-03-19T10:00:00Z ./build.sh mac
```

## Structure du dépôt

```text
ai-project-architect/
├── apa/           # commandes CLI
├── internal/      # generation, planner, runtime, config, output
├── skills/        # skills repo-local apa-*
├── templates/     # templates embarqués
├── scripts/
├── build.sh
└── Makefile
```
