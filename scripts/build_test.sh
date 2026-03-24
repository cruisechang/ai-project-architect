#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUILD_SCRIPT="${ROOT_DIR}/build.sh"
TMP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/project-generator-build-test.XXXXXX")"

cleanup() {
  rm -rf "${TMP_DIR}"
}
trap cleanup EXIT

assert_file_exists() {
  local file="$1"
  if [[ ! -f "${file}" ]]; then
    echo "[FAIL] expected file not found: ${file}"
    exit 1
  fi
}

assert_non_zero() {
  local status="$1"
  if [[ "${status}" -eq 0 ]]; then
    echo "[FAIL] expected non-zero exit code"
    exit 1
  fi
}

echo "[TEST] explicit linux target"
OUT_DIR="${TMP_DIR}/linux" DRY_RUN=1 bash "${BUILD_SCRIPT}" linux >/dev/null
assert_file_exists "${TMP_DIR}/linux/apa-linux"

echo "[TEST] interactive windows target"
OUT_DIR="${TMP_DIR}/interactive" DRY_RUN=1 bash "${BUILD_SCRIPT}" <<< $'3\n' >/dev/null
assert_file_exists "${TMP_DIR}/interactive/apa-windows.exe"

echo "[TEST] interactive retry then mac"
OUT_DIR="${TMP_DIR}/retry" DRY_RUN=1 bash "${BUILD_SCRIPT}" <<< $'9\n1\n' >/dev/null
assert_file_exists "${TMP_DIR}/retry/apa"

echo "[TEST] custom APP_NAME overrides default binary name"
OUT_DIR="${TMP_DIR}/custom" APP_NAME=apa-custom DRY_RUN=1 bash "${BUILD_SCRIPT}" mac >/dev/null
assert_file_exists "${TMP_DIR}/custom/apa-custom"

echo "[TEST] unknown target should fail"
set +e
OUT_DIR="${TMP_DIR}/bad" DRY_RUN=1 bash "${BUILD_SCRIPT}" bad >/dev/null 2>&1
status=$?
set -e
assert_non_zero "${status}"

echo "[PASS] build.sh tests passed"
