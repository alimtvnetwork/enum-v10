#!/usr/bin/env bash
# ─────────────────────────────────────────────────────────────────
# run.sh — Cross-platform equivalent of run.ps1 for Linux/macOS
# Usage: ./run.sh <command> [options]
# ─────────────────────────────────────────────────────────────────
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_LOG_DIR="$SCRIPT_DIR/data/test-logs"

# ── Colors ──────────────────────────────────────────────────────
CYAN='\033[0;36m'
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
GRAY='\033[0;90m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

header()  { echo -e "\n${CYAN}=== $1 ===${NC}"; }
success() { echo -e "  ${GREEN}✓ $1${NC}"; }
fail()    { echo -e "  ${RED}✗ $1${NC}"; }
warn()    { echo -e "  ${YELLOW}$1${NC}"; }

# ── Helpers ─────────────────────────────────────────────────────
ensure_test_log_dir() { mkdir -p "$TEST_LOG_DIR"; }

open_file() {
    local f="$1"
    if command -v open &>/dev/null; then
        open "$f"
    elif command -v xdg-open &>/dev/null; then
        xdg-open "$f"
    else
        echo -e "  ${GRAY}Open manually: $f${NC}"
    fi
}

open_failing_tests() {
    local f="$TEST_LOG_DIR/failing-tests.txt"
    if [[ -f "$f" ]] && ! grep -q '# Count: 0' "$f"; then
        echo ""
        warn "Opening failing tests log..."
        open_file "$f"
    fi
}

has_flag() {
    local flag="$1"; shift
    for arg in "$@"; do
        [[ "$arg" == "$flag" ]] && return 0
    done
    return 1
}

# ── Git + Deps ──────────────────────────────────────────────────
do_git_pull() {
    header "Pulling latest from remote"
    if git pull 2>&1 | sed 's/^/  /'; then
        success "Git pull complete"
    else
        fail "git pull failed (continuing anyway)"
    fi
}

do_fetch_latest() {
    do_git_pull
    header "Fetching latest dependencies"
    if go mod tidy; then
        success "Dependencies up to date"
    else
        fail "go mod tidy failed"
    fi
}

# ── Build Check ─────────────────────────────────────────────────
do_build_check() {
    local build_path="$1"
    header "Build check: $build_path"
    local output
    if output=$(go build "$build_path" 2>&1); then
        success "Build OK"
        return 0
    else
        fail "Build failed — skipping tests"
        ensure_test_log_dir
        local ts
        ts=$(date '+%Y-%m-%d %H:%M:%S')
        {
            echo "# Failing Tests — $ts"
            echo "# Count: 1"
            echo ""
            echo "# Build Failed — tests were NOT run"
            echo ""
            echo "# ── Build Errors ──"
            echo ""
            echo "$output"
        } > "$TEST_LOG_DIR/failing-tests.txt"
        echo "$output" > "$TEST_LOG_DIR/raw-output.txt"
        echo -e "${RED}$output${NC}" | sed 's/^/  /'
        open_failing_tests
        return 1
    fi
}

# ── Write Test Logs ─────────────────────────────────────────────
write_test_logs() {
    local raw_file="$1"
    ensure_test_log_dir
    local passing_file="$TEST_LOG_DIR/passing-tests.txt"
    local failing_file="$TEST_LOG_DIR/failing-tests.txt"
    local ts
    ts=$(date '+%Y-%m-%d %H:%M:%S')

    # Extract passing/failing test names
    local pass_names fail_names
    pass_names=$(grep -oP '(?<=--- PASS: )\S+' "$raw_file" 2>/dev/null | sort -u || true)
    fail_names=$(grep -oP '(?<=--- FAIL: )\S+' "$raw_file" 2>/dev/null | sort -u || true)

    local pass_count fail_count
    pass_count=$(echo "$pass_names" | grep -c . 2>/dev/null || echo 0)
    fail_count=$(echo "$fail_names" | grep -c . 2>/dev/null || echo 0)

    {
        echo "# Passing Tests — $ts"
        echo "# Count: $pass_count"
        echo ""
        echo "$pass_names"
    } > "$passing_file"

    {
        echo "# Failing Tests — $ts"
        echo "# Count: $fail_count"
        echo ""
        if [[ $fail_count -gt 0 ]]; then
            echo "# ── Summary ──"
            echo "$fail_names" | sed 's/^/  - /'
            echo ""
            echo "# ── Details ──"
            echo ""
            # Include lines around FAIL markers
            grep -A5 '--- FAIL:' "$raw_file" 2>/dev/null || true
        fi
    } > "$failing_file"

    echo ""
    if [[ $pass_count -gt 0 ]]; then success "$pass_count passing test(s) → $passing_file"; fi
    if [[ $fail_count -gt 0 ]]; then fail "$fail_count failing test(s) → $failing_file"
    else success "No failing tests"; fi
    echo -e "  ${GRAY}Raw output → $raw_file${NC}"
}

# ── Lint Checks (safeTest boundaries, autofix, bracecheck) ─────
run_safetest_check() {
    local ctx="$1"
    local script="$SCRIPT_DIR/scripts/check-safetest-boundaries.sh"
    if [[ -f "$script" ]]; then
        echo ""
        warn "Running safeTest boundary + empty-if lint check..."
        if ! bash "$script"; then
            fail "safeTest boundary check failed. Fix reported issues before $ctx."
            exit 1
        fi
    fi
}

run_autofix() {
    local ctx="$1"; shift
    local skip_autofix=false skip_brace=false dry_run=false

    for arg in "$@"; do
        case "$arg" in
            --no-autofix)      skip_autofix=true ;;
            --skip-bracecheck) skip_brace=true ;;
            --dry-run)         dry_run=true ;;
        esac
    done

    # ── Go auto-fixer ──
    if $skip_brace; then
        echo -e "  ${YELLOW}Skipping Go auto-fixer and syntax pre-check (--skip-bracecheck)${NC}"
    elif $skip_autofix; then
        echo -e "  ${YELLOW}Skipping Go auto-fixer (--no-autofix)${NC}"
    else
        local dry_label=""
        local dry_flag=""
        if $dry_run; then
            dry_label=" (dry-run)"
            dry_flag="--dry-run"
        fi
        warn "Running Go auto-fixer${dry_label}..."
        local fix_out
        if fix_out=$(go run ./scripts/autofix/ $dry_flag 2>&1); then
            if [[ -n "$fix_out" ]]; then success "$fix_out"; fi
        else
            echo -e "${RED}$fix_out${NC}"
            fail "Go auto-fixer encountered errors."
        fi
    fi

    # ── Go syntax pre-check (bracecheck) ──
    if $skip_brace; then
        : # already logged above
    else
        warn "Running Go syntax pre-check (bracecheck)..."
        local brace_out
        if brace_out=$(go run ./scripts/bracecheck/ 2>&1); then
            success "$(echo "$brace_out" | tr '\n' ' ' | xargs)"
        else
            echo -e "${RED}$brace_out${NC}"
            fail "Go syntax check failed. Fix reported issues before $ctx."
            exit 1
        fi

        # ── Write syntax-issues.txt report ──
        local syntax_dir="$SCRIPT_DIR/data/coverage"
        mkdir -p "$syntax_dir"
        local syntax_file="$syntax_dir/syntax-issues.txt"
        local syntax_ts
        syntax_ts=$(date '+%Y-%m-%d %H:%M:%S')
        local brace_summary
        brace_summary=$(echo "$brace_out" | tr '\n' ' ' | xargs)

        if [[ -f "$syntax_file" ]]; then
            # Append bracecheck results after autofix section
            {
                echo ""
                echo "────────────────────────────────────────────────────────────────────────────────"
                echo " BRACECHECK RESULTS"
                echo "────────────────────────────────────────────────────────────────────────────────"
                echo ""
                echo "  $brace_summary"
                echo ""
                echo "================================================================================"
            } >> "$syntax_file"
        else
            {
                echo "================================================================================"
                echo "  Syntax Issues Report — $syntax_ts"
                echo "  Generated by: autofix + bracecheck pipeline"
                echo "================================================================================"
                echo ""
                echo "────────────────────────────────────────────────────────────────────────────────"
                echo " BRACECHECK RESULTS"
                echo "────────────────────────────────────────────────────────────────────────────────"
                echo ""
                echo "  $brace_summary"
                echo ""
                echo "================================================================================"
            } > "$syntax_file"
        fi
    fi
}

# ── Commands ────────────────────────────────────────────────────

cmd_test() {
    header "Running all tests"
    do_fetch_latest
    pushd tests > /dev/null
    if ! do_build_check "./..."; then popd > /dev/null; return; fi

    local raw_file="$TEST_LOG_DIR/raw-output.txt"
    ensure_test_log_dir
    set +e
    go test -v -count=1 ./... 2>&1 | grep -v '^warning: no packages being tested' | tee "$raw_file"
    local ec=$?
    set -e
    popd > /dev/null

    write_test_logs "$raw_file"
    if [[ $ec -eq 0 ]]; then success "All tests passed"
    else fail "Some tests failed (exit code: $ec)"; fi
    open_failing_tests
}

cmd_test_pkg() {
    local pkg="${1:-}"
    if [[ -z "$pkg" ]]; then
        fail "Package name required. Usage: ./run.sh tp <package>"
        echo -e "  ${YELLOW}Available packages:${NC}"
        ls -1 tests/integratedtests/ 2>/dev/null | sed 's/^/    - /'
        return
    fi

    header "Running tests for package: $pkg"
    do_fetch_latest
    pushd tests > /dev/null
    if ! do_build_check "./integratedtests/$pkg/..."; then popd > /dev/null; return; fi

    local raw_file="$TEST_LOG_DIR/raw-output.txt"
    ensure_test_log_dir
    set +e
    go test -v -count=1 "./integratedtests/$pkg/..." 2>&1 | grep -v '^warning: no packages being tested' | tee "$raw_file"
    local ec=$?
    set -e
    popd > /dev/null

    write_test_logs "$raw_file"
    if [[ $ec -eq 0 ]]; then success "Package tests passed"
    else fail "Package tests failed (exit code: $ec)"; fi
    open_failing_tests
}

cmd_test_coverage() {
    header "Running tests with coverage"
    do_fetch_latest

    # Clean data folder
    local data_dir="$SCRIPT_DIR/data"
    if [[ -d "$data_dir" ]]; then
        rm -rf "$data_dir"
        warn "Cleaned data/ folder"
    fi

    local cover_dir="$SCRIPT_DIR/data/coverage"
    local partial_dir="$cover_dir/partial"
    mkdir -p "$partial_dir"

    local cover_profile="$cover_dir/coverage.out"
    local cover_html="$cover_dir/coverage.html"
    local cover_summary="$cover_dir/coverage-summary.txt"

    # Build coverpkg list (source packages only, excluding tests/)
    local cov_pkg_list
    cov_pkg_list=$(go list ./... 2>&1 | grep -v '/tests/' | paste -sd, -)

    # Get all test packages
    local all_test_pkgs
    all_test_pkgs=$(go list ./tests/integratedtests/... 2>&1 | grep -v '^warning:' | sort)

    # ── Lint checks ──
    run_safetest_check "TC"
    run_autofix "TC" "$@"

    # ── Pre-coverage compile check ──
    local total_pkgs
    total_pkgs=$(echo "$all_test_pkgs" | wc -l | xargs)
    header "Pre-coverage compile check ($total_pkgs packages)"

    local blocked_pkgs=()
    local test_pkgs=()

    while IFS= read -r pkg; do
        [[ -z "$pkg" ]] && continue
        set +e
        go test -count=1 -run '^$' -gcflags=all=-e "-coverpkg=$cov_pkg_list" "$pkg" > /dev/null 2>&1
        local ec=$?
        set -e

        local short_name
        short_name=$(echo "$pkg" | sed 's|.*integratedtests/\?||')
        [[ -z "$short_name" ]] && short_name="(root)"

        if [[ $ec -eq 0 ]]; then
            test_pkgs+=("$pkg")
        else
            blocked_pkgs+=("$short_name")
        fi
    done <<< "$all_test_pkgs"

    # Print blocked summary
    if [[ ${#blocked_pkgs[@]} -gt 0 ]]; then
        echo ""
        echo -e "  ${RED}┌─────────────────────────────────────────────────${NC}"
        echo -e "  ${RED}│ BLOCKED PACKAGES (${#blocked_pkgs[@]} failed to compile)${NC}"
        echo -e "  ${RED}│${NC}"
        for bp in "${blocked_pkgs[@]}"; do
            echo -e "  ${RED}│   ✗ $bp${NC}"
        done
        echo -e "  ${RED}│${NC}"
        echo -e "  ${YELLOW}│ These packages will be SKIPPED in coverage.${NC}"
        echo -e "  ${YELLOW}│ Fix their build errors to include them.${NC}"
        echo -e "  ${RED}└─────────────────────────────────────────────────${NC}"
        echo ""
    else
        echo ""
        success "All ${#test_pkgs[@]} packages compiled successfully"
    fi

    if [[ ${#test_pkgs[@]} -eq 0 ]]; then
        fail "No packages compiled — aborting coverage run"
        return
    fi

    # ── Coverage run ──
    echo ""
    warn "Running ${#test_pkgs[@]} test packages..."

    local all_output="$TEST_LOG_DIR/raw-output.txt"
    ensure_test_log_dir
    > "$all_output"
    local overall_exit=0
    local pkg_index=0

    for pkg in "${test_pkgs[@]}"; do
        pkg_index=$((pkg_index + 1))
        local partial_profile="$partial_dir/cover-${pkg_index}.out"
        set +e
        go test -count=1 "-coverprofile=$partial_profile" "-coverpkg=$cov_pkg_list" "$pkg" 2>&1 >> "$all_output"
        local ec=$?
        set -e
        if [[ $ec -ne 0 ]]; then overall_exit=$ec; fi
    done

    write_test_logs "$all_output"

    # ── Merge partial profiles (MAX count per line) ──
    {
        echo "mode: set"
        cat "$partial_dir"/cover-*.out 2>/dev/null | grep -v '^mode:' | sort
    } > "$cover_profile"

    # ── Generate reports ──
    if [[ -f "$cover_profile" ]]; then
        local func_output
        func_output=$(go tool cover "-func=$cover_profile" 2>&1)

        # HTML report
        if ! go tool cover "-html=$cover_profile" "-o=$cover_html" 2>/dev/null; then
            warn "Failed to generate HTML report"
        fi

        # Summary
        local total_line
        total_line=$(echo "$func_output" | grep '^total:' | tail -1)
        {
            echo "# Coverage Summary — $(date '+%Y-%m-%d %H:%M:%S')"
            echo ""
            if [[ -n "$total_line" ]]; then
                echo "## Total Coverage"
                echo "  $total_line"
                echo ""
            fi
        } > "$cover_summary"

        # Console summary
        if [[ -n "$total_line" ]]; then
            echo ""
            echo -e "  ${CYAN}$total_line${NC}"
        fi

        echo ""
        success "Coverage profile:  $cover_profile"
        success "HTML report:       $cover_html"
        success "Summary:           $cover_summary"
        local syntax_file="$cover_dir/syntax-issues.txt"
        if [[ -f "$syntax_file" ]]; then
            success "Syntax issues:     $syntax_file"
        fi

        # Open HTML if --open flag
        if has_flag "--open" "$@" && [[ -f "$cover_html" ]]; then
            echo ""
            warn "Opening HTML coverage report..."
            open_file "$cover_html"
        fi
    fi

    open_failing_tests
}

cmd_test_coverage_pkg() {
    local pkg="${1:-}"
    shift || true
    if [[ -z "$pkg" ]]; then
        fail "Usage: ./run.sh tcp <package-name>"
        echo -e "  ${GRAY}Example: ./run.sh tcp regexnewtests${NC}"
        return
    fi

    header "Running coverage for package: $pkg"
    do_fetch_latest

    # Clean data folder
    local data_dir="$SCRIPT_DIR/data"
    if [[ -d "$data_dir" ]]; then rm -rf "$data_dir"; warn "Cleaned data/ folder"; fi

    # ── Lint checks ──
    run_safetest_check "TCP"
    run_autofix "TCP" "$@"

    # Build check
    pushd tests > /dev/null
    if ! do_build_check "./integratedtests/$pkg/..."; then popd > /dev/null; return; fi
    popd > /dev/null

    local cover_dir="$SCRIPT_DIR/data/coverage"
    mkdir -p "$cover_dir"
    local cover_profile="$cover_dir/coverage-$pkg.out"
    local cover_html="$cover_dir/coverage-$pkg.html"
    local cover_summary="$cover_dir/coverage-$pkg-summary.txt"

    local cov_pkg_list
    cov_pkg_list=$(go list ./... 2>&1 | grep -v '/tests/' | paste -sd, -)

    local raw_file="$TEST_LOG_DIR/raw-output.txt"
    ensure_test_log_dir
    set +e
    go test -v -count=1 "-coverprofile=$cover_profile" "-coverpkg=$cov_pkg_list" \
        "./tests/integratedtests/$pkg/..." 2>&1 | tee "$raw_file"
    local ec=$?
    set -e

    write_test_logs "$raw_file"

    if [[ -f "$cover_profile" ]]; then
        go tool cover "-func=$cover_profile" 2>&1 | tail -1
        go tool cover "-html=$cover_profile" "-o=$cover_html" 2>/dev/null || true

        echo ""
        success "Coverage profile:  $cover_profile"
        success "HTML report:       $cover_html"
        local syntax_file="$(dirname "$cover_profile")/../syntax-issues.txt"
        [[ -f "$cover_dir/syntax-issues.txt" ]] && success "Syntax issues:     $cover_dir/syntax-issues.txt"

        if has_flag "--open" "$@" && [[ -f "$cover_html" ]]; then
            warn "Opening HTML coverage report..."
            open_file "$cover_html"
        fi
    fi

    open_failing_tests
}

cmd_test_integrated() {
    header "Running integrated tests only"
    do_fetch_latest
    pushd tests > /dev/null
    if ! do_build_check "./integratedtests/..."; then popd > /dev/null; return; fi

    local raw_file="$TEST_LOG_DIR/raw-output.txt"
    ensure_test_log_dir
    set +e
    go test -v -count=1 ./integratedtests/... 2>&1 | tee "$raw_file"
    local ec=$?
    set -e
    popd > /dev/null

    write_test_logs "$raw_file"
    if [[ $ec -eq 0 ]]; then success "Integrated tests passed"
    else fail "Integrated tests failed (exit code: $ec)"; fi
    open_failing_tests
}

cmd_run() {
    header "Running main application"
    go run ./cmd/main/*.go
}

cmd_build() {
    header "Building binary"
    mkdir -p build
    if go build -o build/cli ./cmd/main/; then
        success "Build complete: build/cli"
    else
        fail "Build failed"
    fi
}

cmd_build_run() {
    cmd_build
    if [[ -f build/cli ]]; then
        header "Running built binary"
        ./build/cli
    fi
}

cmd_fmt() {
    header "Formatting Go files"
    gofmt -w -s .
    success "Formatting complete"
}

cmd_vet() {
    header "Running go vet"
    if go vet ./...; then
        success "No issues found"
    else
        fail "Issues found"
    fi
}

cmd_tidy() {
    header "Running go mod tidy"
    go mod tidy
    success "Tidy complete"
}

cmd_goconvey() {
    header "Launching GoConvey"
    if ! command -v goconvey &>/dev/null; then
        warn "GoConvey not found. Installing..."
        go install github.com/smartystreets/goconvey@latest
    fi
    local port="${1:-8080}"
    warn "Starting GoConvey on http://localhost:$port"
    echo -e "  ${GRAY}Press Ctrl+C to stop${NC}"
    pushd tests > /dev/null
    goconvey -port "$port"
    popd > /dev/null
}

cmd_precommit() {
    local single_pkg="${1:-}"
    shift || true
    header "Pre-commit API mismatch checker"

    # Regression guard
    local regression_script="$SCRIPT_DIR/scripts/check-integrated-regressions.sh"
    if [[ -f "$regression_script" ]]; then
        warn "Running regression guard scan..."
        if [[ -n "$single_pkg" ]]; then
            bash "$regression_script" "$single_pkg"
        else
            bash "$regression_script"
        fi
        if [[ $? -ne 0 ]]; then
            fail "Regression guard failed. Fix reported issues before PC."
            exit 1
        fi
    fi

    # safeTest boundary check
    run_safetest_check "PC"
    run_autofix "PC" "$@"

    # Discover test packages with Coverage* files
    local test_base_dir="$SCRIPT_DIR/tests/integratedtests"
    local pkgs_with_coverage=()

    if [[ -n "$single_pkg" ]]; then
        local target_dir="$test_base_dir/$single_pkg"
        if [[ ! -d "$target_dir" ]]; then
            fail "Package not found: $single_pkg"
            return
        fi
        if ls "$target_dir"/Coverage* &>/dev/null 2>&1; then
            pkgs_with_coverage+=("$target_dir")
        fi
    else
        for dir in "$test_base_dir"/*/; do
            if ls "$dir"Coverage* &>/dev/null 2>&1; then
                pkgs_with_coverage+=("$dir")
            fi
        done
    fi

    if [[ ${#pkgs_with_coverage[@]} -eq 0 ]]; then
        success "No Coverage* files found to check"
        return
    fi

    warn "Checking ${#pkgs_with_coverage[@]} packages with Coverage* files..."
    echo ""

    local passed=0 failed=0
    for dir in "${pkgs_with_coverage[@]}"; do
        local rel_path
        rel_path=$(echo "$dir" | sed "s|^$SCRIPT_DIR/||" | sed 's|/$||' | tr '\\' '/')
        local go_pkg="github.com/alimtvnetwork/enum-v10/$rel_path"

        set +e
        go test -c -gcflags=all=-e -o /dev/null "$go_pkg" > /dev/null 2>&1
        local ec=$?
        set -e

        local short_name
        short_name=$(echo "$rel_path" | sed 's|.*integratedtests/||')

        if [[ $ec -eq 0 ]]; then
            passed=$((passed + 1))
        else
            failed=$((failed + 1))
            echo -e "  ${RED}✗ $short_name${NC}"
        fi
    done

    echo ""
    if [[ $failed -eq 0 ]]; then
        echo -e "  ${GREEN}┌─────────────────────────────────────────────────${NC}"
        echo -e "  ${GREEN}│ ✓ ALL $passed PACKAGES PASSED API CHECK${NC}"
        echo -e "  ${GREEN}└─────────────────────────────────────────────────${NC}"
    else
        echo -e "  ${RED}┌─────────────────────────────────────────────────${NC}"
        echo -e "  ${RED}│ ✗ $failed PACKAGE(S) HAVE API MISMATCHES${NC}"
        echo -e "  ${RED}│${NC}"
        echo -e "  ${YELLOW}│ Fix these before committing Coverage* files.${NC}"
        echo -e "  ${RED}└─────────────────────────────────────────────────${NC}"
        exit 1
    fi
}

cmd_show_fail_log() {
    local f="$TEST_LOG_DIR/failing-tests.txt"
    if [[ ! -f "$f" ]]; then
        header "No failing tests log found"
        warn "Run tests first: ./run.sh t"
        return
    fi
    header "Last Failing Tests"
    if grep -q '# Count: 0' "$f"; then
        success "No failing tests in last run"
    else
        cat "$f"
    fi
    echo ""
    echo -e "  ${GRAY}Log file: $f${NC}"
}

cmd_clean() {
    header "Cleaning build artifacts"
    [[ -d build ]] && rm -rf build
    [[ -f tests/coverage.out ]] && rm tests/coverage.out
    local cover_dir="$SCRIPT_DIR/data/coverage"
    [[ -d "$cover_dir" ]] && rm -rf "$cover_dir" && success "Removed coverage reports"
    local pc_dir="$SCRIPT_DIR/data/precommit"
    [[ -d "$pc_dir" ]] && rm -rf "$pc_dir" && success "Removed precommit reports"
    success "Clean complete"
}

cmd_help() {
    echo ""
    echo -e "  ${CYAN}Project Runner — ./run.sh <command> [options]${NC}"
    echo ""
    echo -e "  ${YELLOW}Testing:${NC}"
    echo "    t   | test          Run all tests (verbose)"
    echo "    tp  | test-pkg      Run tests for a specific package"
    echo "    tc  | test-cover    Run tests with coverage (HTML + summary)"
    echo "    tcp | test-cover-pkg Run coverage for a specific package"
    echo "    ti  | test-int      Run integrated tests only"
    echo "    tf  | test-fail     Show last failing tests log"
    echo "    gc  | goconvey      Launch GoConvey (browser test runner)"
    echo ""
    echo -e "  ${YELLOW}Build & Run:${NC}"
    echo "    r   | run           Run the main application"
    echo "    b   | build         Build the binary"
    echo "    br  | build-run     Build then run"
    echo ""
    echo -e "  ${YELLOW}Code Quality:${NC}"
    echo "    f   | fmt           Format all Go files"
    echo "    l   | lint          Run go vet"
    echo "    v   | vet           Run go vet"
    echo "    ty  | tidy          Run go mod tidy"
    echo "    pc  | pre-commit    Check Coverage* files for API mismatches"
    echo ""
    echo -e "  ${YELLOW}Other:${NC}"
    echo "    c   | clean         Clean build artifacts"
    echo "    h   | help          Show this help"
    echo ""
    echo -e "  ${YELLOW}Mode Options (for tc/tcp/pc):${NC}"
    echo "    --open              Open HTML coverage report in browser"
    echo "    --skip-bracecheck   Skip Go syntax pre-check"
    echo "    --no-autofix        Skip Go auto-fixer before bracecheck"
    echo "    --dry-run           Run auto-fixer in preview mode"
    echo ""
    echo -e "  ${GRAY}Examples:${NC}"
    echo "    ./run.sh t"
    echo "    ./run.sh tp regexnewtests"
    echo "    ./run.sh tc"
    echo "    ./run.sh tc --open"
    echo "    ./run.sh tcp regexnewtests"
    echo "    ./run.sh pc"
    echo "    ./run.sh pc corejsontests"
    echo "    ./run.sh gc 9090"
    echo ""
}

# ── Dispatch ────────────────────────────────────────────────────
CMD="${1:-help}"
shift || true

case "$CMD" in
    t|test)              cmd_test "$@" ;;
    tp|test-pkg)         cmd_test_pkg "$@" ;;
    tc|test-cover)       cmd_test_coverage "$@" ;;
    tcp|test-cover-pkg)  cmd_test_coverage_pkg "$@" ;;
    ti|test-int)         cmd_test_integrated "$@" ;;
    tf|test-fail)        cmd_show_fail_log ;;
    gc|goconvey)         cmd_goconvey "$@" ;;
    r|run)               cmd_run ;;
    b|build)             cmd_build ;;
    br|build-run)        cmd_build_run ;;
    f|fmt)               cmd_fmt ;;
    l|lint|v|vet)        cmd_vet ;;
    ty|tidy)             cmd_tidy ;;
    pc|pre-commit)       cmd_precommit "$@" ;;
    c|clean)             cmd_clean ;;
    h|help)              cmd_help ;;
    *)
        fail "Unknown command: '$CMD'"
        cmd_help
        exit 1
        ;;
esac
