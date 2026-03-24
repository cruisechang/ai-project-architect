#!/usr/bin/env bash
set -euo pipefail

TARGET="${1:-}"
OUT_DIR="${OUT_DIR:-release}"
APP_NAME="${APP_NAME:-apa}"
DRY_RUN="${DRY_RUN:-0}"
VERSION="${VERSION:-dev}"
COMMIT="${COMMIT:-$(git rev-parse --short HEAD 2>/dev/null || echo unknown)}"
BUILD_DATE="${BUILD_DATE:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}"
COPY_TO_USR_LOCAL_BIN="${COPY_TO_USR_LOCAL_BIN:-ask}"
BUILT_ARTIFACTS=()

usage() {
  cat <<'USAGE'
Usage: ./build.sh [mac|linux|windows|all|interactive]

When no argument is provided, interactive mode is used.

Environment variables:
  OUT_DIR   Output directory (default: release)
  APP_NAME  Binary name prefix (default: apa)
  DRY_RUN   Set to 1 to skip go build and generate placeholder artifacts
  VERSION   Build version injected by ldflags (default: dev)
  COMMIT    Git commit injected by ldflags (default: current git short sha or unknown)
  BUILD_DATE UTC date injected by ldflags (default: current UTC timestamp)
  COPY_TO_USR_LOCAL_BIN  ask|yes|no (default: ask, ask only when stdin is a TTY)
USAGE
}

choose_target_interactive() {
  local choice=""
  echo "請選擇 build 目標："
  echo "  1) mac"
  echo "  2) linux"
  echo "  3) windows"
  echo "  4) all"

  while true; do
    read -r -p "輸入選項 [1-4] (預設 1): " choice
    case "${choice:-1}" in
      1)
        TARGET="mac"
        return 0
        ;;
      2)
        TARGET="linux"
        return 0
        ;;
      3)
        TARGET="windows"
        return 0
        ;;
      4)
        TARGET="all"
        return 0
        ;;
      *)
        echo "無效選項：${choice}. 請輸入 1, 2, 3 或 4。"
        ;;
    esac
  done
}

mkdir -p "${OUT_DIR}"

build_target() {
  local goos="$1"
  local goarch="$2"
  local ext="$3"
  local out_name="$4"

  echo "Building ${goos}/${goarch} -> ${OUT_DIR}/${out_name}${ext}"
  if [[ "${DRY_RUN}" == "1" ]]; then
    printf 'dry-run artifact for %s/%s\n' "${goos}" "${goarch}" >"${OUT_DIR}/${out_name}${ext}"
    return 0
  fi
  local ldflags="-X project-generator/internal/buildinfo.Version=${VERSION} -X project-generator/internal/buildinfo.Commit=${COMMIT} -X project-generator/internal/buildinfo.BuildDate=${BUILD_DATE}"
  CGO_ENABLED=0 GOOS="${goos}" GOARCH="${goarch}" go build -ldflags "${ldflags}" -o "${OUT_DIR}/${out_name}${ext}" .
  BUILT_ARTIFACTS+=("${OUT_DIR}/${out_name}${ext}")
}

copy_binaries_to_usr_local_bin() {
  local artifact=""
  local dest=""
  local failed=0

  for artifact in "${BUILT_ARTIFACTS[@]}"; do
    dest="/usr/local/bin/$(basename "${artifact}")"
    echo "Copying ${artifact} -> ${dest}"

    if cp "${artifact}" "${dest}" 2>/dev/null; then
      chmod +x "${dest}" 2>/dev/null || true
      continue
    fi

    if command -v sudo >/dev/null 2>&1 && sudo cp "${artifact}" "${dest}"; then
      sudo chmod +x "${dest}" 2>/dev/null || true
      continue
    fi

    echo "Failed to copy ${artifact} to ${dest}" >&2
    failed=1
  done

  if [[ "${failed}" -ne 0 ]]; then
    return 1
  fi
}

maybe_copy_binaries() {
  local answer=""

  if [[ "${DRY_RUN}" == "1" || "${#BUILT_ARTIFACTS[@]}" -eq 0 ]]; then
    return 0
  fi

  case "${COPY_TO_USR_LOCAL_BIN}" in
    yes)
      copy_binaries_to_usr_local_bin
      return $?
      ;;
    no)
      return 0
      ;;
    ask)
      if [[ ! -t 0 ]]; then
        return 0
      fi
      ;;
    *)
      echo "Invalid COPY_TO_USR_LOCAL_BIN value: ${COPY_TO_USR_LOCAL_BIN} (expected ask|yes|no)" >&2
      return 1
      ;;
  esac

  read -r -p "是否要將 binary 複製到 /usr/local/bin？(y/N): " answer
  case "${answer,,}" in
    y|yes)
      copy_binaries_to_usr_local_bin
      ;;
    *)
      ;;
  esac
}

if [[ "${TARGET}" == "-h" || "${TARGET}" == "--help" ]]; then
  usage
  exit 0
fi

if [[ -z "${TARGET}" || "${TARGET}" == "interactive" ]]; then
  choose_target_interactive
fi

TARGET="$(printf '%s' "${TARGET}" | tr '[:upper:]' '[:lower:]')"

case "${TARGET}" in
  mac)
    build_target "darwin" "arm64" "" "${APP_NAME}"
    ;;
  linux)
    build_target "linux" "amd64" "" "${APP_NAME}-linux"
    ;;
  windows)
    build_target "windows" "amd64" ".exe" "${APP_NAME}-windows"
    ;;
  all)
    build_target "darwin" "arm64" "" "${APP_NAME}"
    build_target "linux" "amd64" "" "${APP_NAME}-linux"
    build_target "windows" "amd64" ".exe" "${APP_NAME}-windows"
    ;;
  *)
    echo "Unknown target: ${TARGET}"
    usage
    exit 1
    ;;
esac

maybe_copy_binaries

echo "Done. Artifacts are in ${OUT_DIR}/"
