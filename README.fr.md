# apa - AI Project Architect

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

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![CLI](https://img.shields.io/badge/Type-CLI-111111)](#commandes)
[![Skills](https://img.shields.io/badge/Repo%20Skills-apa--*-2f855a)](#commands--skills)

## `apa簡要說明`

- 可用觸發詞：`apa簡要說明`、`apa 簡要說明`、`apa說明`、`apa 說明`
- Démarrer la boucle d'implémentation : `apa prompt --reviewer agent-self` -> coller la sortie dans l'agent -> demander `apa-loop` + `apa-implement`
- Réviser d'abord les documents : `apa prompt --docs-only` -> coller la sortie dans l'agent -> demander uniquement `apa-doc-review`
- Wrapper de boucle terminal : `bash scripts/apa-loop-setup.sh --max-iterations 30 --reviewer agent-self`
- Arrêter le wrapper : `bash scripts/apa-loop-cancel.sh`

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
  -> apa prompt
  -> agent implements
  -> make test
  -> repeat
```

## Ce que fait `apa`

- Initialise un nouveau projet à partir d'une idée produit en langage naturel
- Déduit une stack technique pragmatique, modifiable ensuite via des flags
- Génère une documentation par phases à partir de `Phase 0`, incluant PRD, SPEC, ARCHITECTURE, API, DB Schema et plan d'implémentation
- Crée un starter runnable, les tests, un Makefile et la configuration d'agent
- Produit avec `apa prompt` un prompt structuré pour continuer la livraison avec un agent IA
- S'associe naturellement à `apa-loop` pour une livraison par tours : lire l'état, choisir 1 à 3 tâches, vérifier, mettre à jour l'état, répéter

## Workflow recommandé

```bash
# 1. Compiler apa
go build -o apa .

# 2. Créer un nouveau projet en dehors du repo cible
./apa init --idea "Plateforme SaaS de reporting" --name report-platform --path ~/projects

# 3. Entrer dans le repo généré et itérer
cd ~/projects/report-platform
./apa list-skills
./apa prompt --reviewer agent-self
make test
```

Boucle principale :

1. Exécuter `apa init` une fois pour créer la première version du projet.
2. Maintenir une documentation par phases à partir de `Phase 0` afin d'aligner périmètre, tests, gates et rapports entre PRD/API/SPEC.
3. Utiliser par défaut `apa-loop` avec `apa-implement` comme boucle de livraison.
4. Lancer `apa prompt --reviewer agent-self`, laisser l'agent travailler, puis valider avec `make test`.
5. Répéter jusqu'à obtenir un dépôt livrable.

## État de la boucle de livraison et utilisation de `apa-loop`

Les dépôts générés doivent garder `docs/IMPLEMENTATION_STATUS.md` ou `TASKS.md` à jour.
Utilise `apa-loop` avec `apa-implement` pour que l'agent continue à enchaîner implémentation, tests, corrections et mises à jour de documentation jusqu'à ce que le gate de fin soit satisfait.
`apa-loop` est le skill repo-local qui force la boucle de livraison par tours : lire le fichier d'état, choisir 1 à 3 tâches vérifiables, exécuter les tests/contrôles, mettre à jour l'état, puis répéter jusqu'à satisfaire le gate de fin.
Utilisation :
- Si tu veux d'abord itérer sur les documents avant l'implémentation, lance `apa prompt --docs-only` et demande à l'agent d'utiliser uniquement `apa-doc-review`.
- Flux principal via agent (commun à Codex et Claude Code) : lance `apa prompt --reviewer agent-self`, ou lance `apa prompt` puis choisis le reviewer en interactif, puis demande explicitement à l'agent d'utiliser `apa-loop` avec `apa-implement`
- Wrapper terminal optionnel (pour les environnements qui exposent un hook ou une slash command générés, comme Claude Code) : `bash scripts/apa-loop-setup.sh --max-iterations 30 --reviewer agent-self`
- Slash command optionnelle : `/apa-loop --max-iterations 30 --reviewer agent-self`
- Commande d'annulation optionnelle : `/cancel-apa-loop`
- Review policy: indique le reviewer une seule fois via `--reviewer` ou dans la première consigne, puis réutilise-le aux tours suivants sauf si tu demandes explicitement un changement (`agent-self`, `apa-codex-review`, `apa-claude-review`).

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
./apa prompt --reviewer agent-self > prompt.md
make test
```

## Commandes

| Commande | Rôle |
|---|---|
| `apa init` | Créer un nouveau projet à partir d'une idée |
| `apa prompt` | Générer le prompt IA structuré pour la suite |
| `apa list-skills` | Afficher les skills repo-local disponibles |
| `apa doctor` | Vérifier l'environnement local et le chemin des skills |
| `apa version` | Afficher les informations de version du build |

Options complètes : `apa <command> --help`

## Commands & Skills

Commandes CLI actuelles :

- `init`
- `prompt`
- `list-skills`
- `doctor`
- `version`

Skills repo-local actuels :

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
# Lance directement le Wizard, avec Project idea comme première question.

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
| `--type` | `cli`, `server`, `web-app-server`, `mobile-app-server`, `web-app`, `mobile-app` |
| `--agent` | `codex`, `claude-code` ou `universal` |
| `--backend` | `go`, `python`, `node` ou `none` |
| `--frontend` | `react`, `next`, `nuxt`, `vue`, `pure-typescript` ou `none` |
| `--unit-test` `--api-test` `--integration-test` `--e2e-test` | `yes` or `no` |
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

## `apa prompt`

`apa prompt` lit l'état actuel du dépôt et imprime un prompt structuré pour l'agent IA.

La commande est utile avant l'implémentation, pendant le développement ou après une régression, afin que l'agent continue à travailler en tenant compte des documents, des tâches et des contraintes du repo.

La commande vérifie aussi si les documents existants utilisent des sections alignées par priorité `Phase 0`, `Phase 1`, ... Si ce n'est pas le cas, `apa prompt` avertit l'utilisateur et demande à l'agent de les réécrire avec `apa-docs` avant de continuer.

```bash
./apa prompt
./apa prompt --reviewer agent-self
./apa prompt --docs-only
./apa prompt --root ~/projects/report-platform
./apa prompt > prompt.md
```

## Skills repo-local

Ce dépôt utilise la série de noms `apa-*` pour les skills repo-local.

Exemples actuels :

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

`apa-docs` rédige la documentation par phases alignées par priorité (`Phase 0`, `Phase 1`, ...). `Phase 0` est toujours la phase la plus prioritaire. Chaque phase doit préciser le périmètre, les contenus PRD/API/SPEC alignés, les tests requis, les points de contrôle, les critères d'achèvement, un gate explicite vers la phase suivante et un rapport de fin de phase.
`apa-doc-review` sert à une boucle de révision documentaire avec l'utilisateur : ne modifier que les docs, s'arrêter après chaque révision pour recueillir un retour, et ne lancer l'implémentation qu'après approbation explicite des documents.

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
