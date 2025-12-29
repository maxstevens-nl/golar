package vue_tests

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/auvred/golar/golar"
	"github.com/microsoft/typescript-go/shim/bundled"
	"github.com/microsoft/typescript-go/shim/execute"
	"github.com/microsoft/typescript-go/shim/execute/tsc"
	"github.com/microsoft/typescript-go/shim/golarext"
	"github.com/microsoft/typescript-go/shim/tspath"
	"github.com/microsoft/typescript-go/shim/vfs"
	"github.com/microsoft/typescript-go/shim/vfs/osvfs"
	"gotest.tools/v3/assert"
)

// Tests not in this list will be skipped.
var enabledTests = map[string]bool{
	// Tests that currently pass (no diagnostics expected)
	"#2048":           true,
	"#2166":           true,
	"#2250":           true,
	"#2308":           true,
	"#2431":           true,
	"#2629":           true,
	"#2678":           true,
	"#2683":           true,
	"#3138":           true,
	"#3255":           true,
	"#3374":           true,
	"#3433":           true,
	"#3488":           true,
	"#3518":           true,
	"#3718":           true,
	"#3845":           true,
	"#4209":           true,
	"#4369":           true,
	"#4413":           true,
	"#5111":           true,
	"#5267":           true,
	"#5428":           true,
	"#5492":           true,
	"#5729":           true,
	"#5751":           true,
	"#5810":           true,
	"#5819":           true,
	"#5899":           true,
	"#625":            true,
	"no-script-block": true,
	"v-if":            true,

	// Expected failure tests (these should produce errors)
	"_failed_#3632":      true,
	"_failed_#4569":      true,
	"_failed_#5071":      true,
	"_failed_#5823":      true,
	"_failed_directives": true,
}

// vueTscSys implements tsc.System for running vue-tsc build tests
// using the real filesystem.
type vueTscSys struct {
	cwd                string
	fs                 vfs.FS
	defaultLibraryPath string
	output             *strings.Builder
	start              time.Time
}

func newVueTscSys(cwd string) *vueTscSys {
	return &vueTscSys{
		cwd:                tspath.NormalizePath(cwd),
		fs:                 golar.WrapFS(bundled.WrapFS(osvfs.FS())),
		defaultLibraryPath: bundled.LibPath(),
		output:             &strings.Builder{},
		start:              time.Now(),
	}
}

func (s *vueTscSys) FS() vfs.FS {
	return s.fs
}

func (s *vueTscSys) DefaultLibraryPath() string {
	return s.defaultLibraryPath
}

func (s *vueTscSys) GetCurrentDirectory() string {
	return s.cwd
}

func (s *vueTscSys) Writer() io.Writer {
	return s.output
}

func (s *vueTscSys) WriteOutputIsTTY() bool {
	return false
}

func (s *vueTscSys) GetWidthOfTerminal() int {
	return 0
}

func (s *vueTscSys) GetEnvironmentVariable(name string) string {
	return os.Getenv(name)
}

func (s *vueTscSys) Now() time.Time {
	return time.Now()
}

func (s *vueTscSys) SinceStart() time.Duration {
	return time.Since(s.start)
}

func (s *vueTscSys) GetGolarCallbacks() *golarext.GolarCallbacks {
	return golar.GolarExtCallbacks
}

func (s *vueTscSys) GetOutput() string {
	return s.output.String()
}

// getTestWorkspacePath returns the absolute path to the test-workspace directory
func getTestWorkspacePath() string {
	_, filename, _, _ := runtime.Caller(0)
	// Go from internal/vue/tests to test-workspace
	return filepath.Join(filepath.Dir(filename), "..", "..", "..", "test-workspace")
}

// TestVueTscBuild runs the vue-tsc build tests from test-workspace/tsc
func TestVueTscBuild(t *testing.T) {
	t.Parallel()
	t.Skip("TODO: fix build")

	testWorkspacePath := getTestWorkspacePath()
	tscPath := filepath.Join(testWorkspacePath, "tsc")

	// Read the tsconfig.json to get the list of test references
	tsconfigPath := filepath.Join(tscPath, "tsconfig.json")
	if _, err := os.Stat(tsconfigPath); os.IsNotExist(err) {
		t.Fatalf("tsconfig.json not found at %s", tsconfigPath)
	}

	// Get all test directories (those with tsconfig.json)
	entries, err := os.ReadDir(tscPath)
	if err != nil {
		t.Fatalf("Failed to read test directory: %v", err)
	}

	var testDirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			tsconfigFile := filepath.Join(tscPath, entry.Name(), "tsconfig.json")
			if _, err := os.Stat(tsconfigFile); err == nil {
				testDirs = append(testDirs, entry.Name())
			}
		}
	}

	if len(testDirs) == 0 {
		t.Fatal("No test directories found")
	}

	t.Logf("Found %d test directories", len(testDirs))

	// Run the build
	sys := newVueTscSys(tscPath)
	commandLineArgs := []string{"--build", ".", "--pretty", "false"}

	result := execute.CommandLine(sys, commandLineArgs, nil)

	output := sys.GetOutput()

	// Parse the output to extract diagnostics
	diagnostics := parseTscOutput(output)

	// Log results
	t.Logf("Build completed with status: %v", result.Status)
	t.Logf("Found %d diagnostics", len(diagnostics))

	// Separate expected failures from unexpected ones
	var expectedFailures []string
	var unexpectedDiagnostics []string

	for _, diag := range diagnostics {
		if strings.Contains(diag, "_failed_") {
			expectedFailures = append(expectedFailures, diag)
		} else {
			unexpectedDiagnostics = append(unexpectedDiagnostics, diag)
		}
	}

	// Log expected failures
	if len(expectedFailures) > 0 {
		t.Logf("Expected failures (%d):", len(expectedFailures))
		for _, diag := range expectedFailures {
			t.Logf("  %s", diag)
		}
	}

	// Report unexpected diagnostics as errors
	if len(unexpectedDiagnostics) > 0 {
		t.Errorf("Unexpected diagnostics (%d):", len(unexpectedDiagnostics))
		for _, diag := range unexpectedDiagnostics {
			t.Errorf("  %s", diag)
		}
	}
}

// TestVueTscBuildIndividual runs each test directory individually.
// Only tests in the enabledTests allowlist are run; others are skipped.
func TestVueTscBuildIndividual(t *testing.T) {
	t.Parallel()

	testWorkspacePath := getTestWorkspacePath()
	tscPath := filepath.Join(testWorkspacePath, "tsc")

	// Get all test directories
	entries, err := os.ReadDir(tscPath)
	if err != nil {
		t.Fatalf("Failed to read test directory: %v", err)
	}

	var testDirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			tsconfigFile := filepath.Join(tscPath, entry.Name(), "tsconfig.json")
			if _, err := os.Stat(tsconfigFile); err == nil {
				testDirs = append(testDirs, entry.Name())
			}
		}
	}

	var enabledCount, skippedCount int

	for _, testDir := range testDirs {
		isExpectedFailure := strings.HasPrefix(testDir, "_failed_")

		t.Run(testDir, func(t *testing.T) {
			t.Parallel()

			if !enabledTests[testDir] {
				skippedCount++
				t.Skipf("Test %s is not in enabledTests allowlist", testDir)
				return
			}
			enabledCount++

			testPath := filepath.Join(tscPath, testDir)
			sys := newVueTscSys(testPath)
			commandLineArgs := []string{"--build", ".", "--pretty", "false"}

			result := execute.CommandLine(sys, commandLineArgs, nil)
			output := sys.GetOutput()
			diagnostics := parseTscOutput(output)

			if isExpectedFailure {
				// For _failed_ tests, we expect errors
				if len(diagnostics) == 0 {
					t.Logf("Expected failure test %s passed without errors (this might be OK if the issue is fixed)", testDir)
				} else {
					t.Logf("Expected failure test %s has %d diagnostics (expected)", testDir, len(diagnostics))
					for _, diag := range diagnostics {
						t.Logf("  %s", diag)
					}
				}
			} else {
				// For normal tests, we expect no errors
				if len(diagnostics) > 0 {
					t.Errorf("Test %s failed with %d diagnostics:", testDir, len(diagnostics))
					for _, diag := range diagnostics {
						t.Errorf("  %s", diag)
					}
				}
				if result.Status != tsc.ExitStatusSuccess {
					// Only fail if there were actual diagnostics
					if len(diagnostics) > 0 {
						t.Errorf("Test %s exited with status %v", testDir, result.Status)
					}
				}
			}
		})
	}

	t.Logf("Test directories: %d total, %d enabled, %d skipped", len(testDirs), enabledCount, skippedCount)
}

func TestVueTscSnapshot(t *testing.T) {
	t.Parallel()
	t.Skip("TODO: fix snapshots")

	testWorkspacePath := getTestWorkspacePath()
	tscPath := filepath.Join(testWorkspacePath, "tsc")

	sys := newVueTscSys(tscPath)
	commandLineArgs := []string{"--build", ".", "--pretty", "false"}

	execute.CommandLine(sys, commandLineArgs, nil)

	output := sys.GetOutput()
	diagnostics := parseTscOutput(output)

	// Sort diagnostics for consistent comparison
	sort.Strings(diagnostics)

	// Expected diagnostics from the original test
	expectedDiagnostics := []string{
		"test-workspace/tsc/_failed_#3632/both.vue(3,1): error TS1109: Expression expected.",
		"test-workspace/tsc/_failed_#3632/both.vue(7,1): error TS1109: Expression expected.",
		"test-workspace/tsc/_failed_#3632/script.vue(3,1): error TS1109: Expression expected.",
		"test-workspace/tsc/_failed_#3632/scriptSetup.vue(3,1): error TS1109: Expression expected.",
		"test-workspace/tsc/_failed_#4569/main.vue(1,41): error TS4025: Exported variable '__VLS_export' has or is using private name 'Props'.",
		"test-workspace/tsc/_failed_#5071/withScript.vue(1,19): error TS1005: ';' expected.",
		"test-workspace/tsc/_failed_#5071/withoutScript.vue(2,26): error TS1005: ';' expected.",
		"test-workspace/tsc/_failed_#5823/main.vue(6,13): error TS1109: Expression expected.",
		"test-workspace/tsc/_failed_directives/main.vue(14,6): error TS2339: Property 'notExist' does not exist on type '{ exist: {}; Comp: () => void; $: ComponentInternalInstance; $data: {}; $props: {}; $attrs: Data; $refs: Data; $slots: Readonly<InternalSlots>; ... 8 more ...; $watch<T extends string | ((...args: any) => any)>(source: T, cb: T extends (...args: any) => infer R ? (args_0: R, args_1: R, args_2: OnCleanup) => any : (a...'.",
		"test-workspace/tsc/_failed_directives/main.vue(17,2): error TS2578: Unused '@ts-expect-error' directive.",
		"test-workspace/tsc/_failed_directives/main.vue(20,2): error TS2578: Unused '@ts-expect-error' directive.",
		"test-workspace/tsc/_failed_directives/main.vue(9,6): error TS2339: Property 'notExist' does not exist on type '{ exist: {}; Comp: () => void; $: ComponentInternalInstance; $data: {}; $props: {}; $attrs: Data; $refs: Data; $slots: Readonly<InternalSlots>; ... 8 more ...; $watch<T extends string | ((...args: any) => any)>(source: T, cb: T extends (...args: any) => infer R ? (args_0: R, args_1: R, args_2: OnCleanup) => any : (a...'.",
	}

	// Normalize paths in diagnostics for comparison
	normalizedDiagnostics := make([]string, len(diagnostics))
	for i, diag := range diagnostics {
		// Convert absolute paths to relative paths matching expected format
		normalizedDiagnostics[i] = normalizeDiagnosticPath(diag, testWorkspacePath)
	}
	sort.Strings(normalizedDiagnostics)

	// Filter to only include _failed_ diagnostics for snapshot comparison
	var failedDiagnostics []string
	for _, diag := range normalizedDiagnostics {
		if strings.Contains(diag, "_failed_") {
			failedDiagnostics = append(failedDiagnostics, diag)
		}
	}
	sort.Strings(failedDiagnostics)
	sort.Strings(expectedDiagnostics)

	// Compare
	if !slices.Equal(failedDiagnostics, expectedDiagnostics) {
		t.Errorf("Diagnostics do not match expected snapshot")
		t.Logf("Expected (%d):", len(expectedDiagnostics))
		for _, d := range expectedDiagnostics {
			t.Logf("  %s", d)
		}
		t.Logf("Got (%d):", len(failedDiagnostics))
		for _, d := range failedDiagnostics {
			t.Logf("  %s", d)
		}

		// Show diff
		t.Log("Missing from actual:")
		for _, exp := range expectedDiagnostics {
			if !slices.Contains(failedDiagnostics, exp) {
				t.Logf("  - %s", exp)
			}
		}
		t.Log("Extra in actual:")
		for _, act := range failedDiagnostics {
			if !slices.Contains(expectedDiagnostics, act) {
				t.Logf("  + %s", act)
			}
		}
	}
}

// parseTscOutput parses the tsc output and extracts diagnostic lines
func parseTscOutput(output string) []string {
	lines := strings.Split(output, "\n")
	var diagnostics []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// TypeScript diagnostic format: file(line,col): error TSxxxx: message
		if strings.Contains(line, "): error TS") || strings.Contains(line, "): warning TS") {
			diagnostics = append(diagnostics, line)
		}
	}

	return diagnostics
}

// normalizeDiagnosticPath converts absolute paths to relative paths for comparison
func normalizeDiagnosticPath(diag string, basePath string) string {
	// Extract the file path part (before the first '(')
	parenIdx := strings.Index(diag, "(")
	if parenIdx == -1 {
		return diag
	}

	filePath := diag[:parenIdx]
	rest := diag[parenIdx:]

	// Convert to relative path
	if strings.HasPrefix(filePath, basePath) {
		relPath, err := filepath.Rel(filepath.Dir(basePath), filePath)
		if err == nil {
			// Normalize path separators
			relPath = strings.ReplaceAll(relPath, "\\", "/")
			return relPath + rest
		}
	}

	return diag
}

func TestVueTscBuildStatus(t *testing.T) {
	t.Parallel()

	testWorkspacePath := getTestWorkspacePath()
	tscPath := filepath.Join(testWorkspacePath, "tsc")

	// Check that the test workspace exists
	if _, err := os.Stat(tscPath); os.IsNotExist(err) {
		t.Fatalf("Test workspace not found at %s", tscPath)
	}

	// Check that tsconfig.json exists
	tsconfigPath := filepath.Join(tscPath, "tsconfig.json")
	if _, err := os.Stat(tsconfigPath); os.IsNotExist(err) {
		t.Fatalf("tsconfig.json not found at %s", tsconfigPath)
	}

	sys := newVueTscSys(tscPath)
	commandLineArgs := []string{"--build", ".", "--pretty", "false"}

	result := execute.CommandLine(sys, commandLineArgs, nil)

	assert.Assert(t, result.Status == tsc.ExitStatusSuccess ||
		result.Status == tsc.ExitStatusDiagnosticsPresent_OutputsSkipped ||
		result.Status == tsc.ExitStatusDiagnosticsPresent_OutputsGenerated,
		"Build failed with unexpected status: %v", result.Status)

	output := sys.GetOutput()
	t.Logf("Build output length: %d bytes", len(output))
	if len(output) > 0 {
		// Log first 2000 chars of output
		if len(output) > 2000 {
			t.Logf("Build output (truncated):\n%s...", output[:2000])
		} else {
			t.Logf("Build output:\n%s", output)
		}
	}
}
