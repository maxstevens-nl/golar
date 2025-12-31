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

// skippedTests contains tests that are known to fail.
// Remove tests from this list as they are fixed.
var skippedTests = map[string]bool{
	"#1886":                         true,
	"#2206":                         true,
	"#2370":                         true,
	"#2472":                         true,
	"#2554":                         true,
	"#2586":                         true,
	"#2590":                         true,
	"#2617":                         true,
	"#2639":                         true,
	"#2646":                         true,
	"#2709":                         true,
	"#2712":                         true,
	"#2720":                         true,
	"#2730":                         true,
	"#2744":                         true,
	"#2758":                         true,
	"#3102":                         true,
	"#3121":                         true,
	"#3123":                         true,
	"#3164":                         true,
	"#3204":                         true,
	"#3289":                         true,
	"#329":                          true,
	"#3295":                         true,
	"#3311":                         true,
	"#3318":                         true,
	"#3353":                         true,
	"#3373":                         true,
	"#3405":                         true,
	"#3414":                         true,
	"#3440":                         true,
	"#3476":                         true,
	"#3561":                         true,
	"#3592":                         true,
	"#3612":                         true,
	"#3615":                         true,
	"#3656":                         true,
	"#3671":                         true,
	"#3672":                         true,
	"#3688":                         true,
	"#3732":                         true,
	"#3748":                         true,
	"#3756":                         true,
	"#3779":                         true,
	"#3820":                         true,
	"#3997":                         true,
	"#4050":                         true,
	"#4327":                         true,
	"#4353":                         true,
	"#4386":                         true,
	"#4433":                         true,
	"#4503":                         true,
	"#4512":                         true,
	"#4537":                         true,
	"#4539":                         true,
	"#4646":                         true,
	"#4668":                         true,
	"#4682":                         true,
	"#4699":                         true,
	"#4785":                         true,
	"#4799":                         true,
	"#4812":                         true,
	"#4820":                         true,
	"#4822":                         true,
	"#4826":                         true,
	"#4827":                         true,
	"#4828":                         true,
	"#4878":                         true,
	"#4972":                         true,
	"#5067":                         true,
	"#5106":                         true,
	"#5120":                         true,
	"#5136":                         true,
	"#5157":                         true,
	"#5159":                         true,
	"#5228":                         true,
	"#5338":                         true,
	"#5474":                         true,
	"#5592":                         true,
	"#5604":                         true,
	"#5617":                         true,
	"#5776":                         true,
	"#5780":                         true,
	"#5840":                         true,
	"#5843":                         true,
	"#5895":                         true,
	"attrs":                         true,
	"components":                    true,
	"cssModule":                     true,
	"defineExpose":                  true,
	"defineModel":                   true,
	"defineOptions":                 true,
	"directiveComments":             true,
	"directives":                    true,
	"events":                        true,
	"fallthroughAttributes":         true,
	"fallthroughAttributes_generic": true,
	"no-script-block":               true,
	"rootEl":                        true,
	"script-setup-scope":            true,
	"script_src":                    true,
	"slots":                         true,
	"templateRef":                   true,
	"type-helpers":                  true,
	"v-for":                         true,
	"v-generic":                     true,
	"withDefaults":                  true,
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

func getTestWorkspacePath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "..", "test-workspace")
}

func TestVueTscBuild(t *testing.T) {
	t.Parallel()
	t.Skip("TODO: fix build")

	testWorkspacePath := getTestWorkspacePath()
	tscPath := filepath.Join(testWorkspacePath, "tsc")

	tsconfigPath := filepath.Join(tscPath, "tsconfig.json")
	if _, err := os.Stat(tsconfigPath); os.IsNotExist(err) {
		t.Fatalf("tsconfig.json not found at %s", tsconfigPath)
	}

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

	sys := newVueTscSys(tscPath)
	commandLineArgs := []string{"--build", ".", "--pretty", "false"}

	result := execute.CommandLine(sys, commandLineArgs, nil)

	output := sys.GetOutput()

	diagnostics := parseTscOutput(output)

	t.Logf("Build completed with status: %v", result.Status)
	t.Logf("Found %d diagnostics", len(diagnostics))

	var expectedFailures []string
	var unexpectedDiagnostics []string

	for _, diag := range diagnostics {
		if strings.Contains(diag, "_failed_") {
			expectedFailures = append(expectedFailures, diag)
		} else {
			unexpectedDiagnostics = append(unexpectedDiagnostics, diag)
		}
	}

	if len(expectedFailures) > 0 {
		t.Logf("Expected failures (%d):", len(expectedFailures))
		for _, diag := range expectedFailures {
			t.Logf("  %s", diag)
		}
	}

	if len(unexpectedDiagnostics) > 0 {
		t.Errorf("Unexpected diagnostics (%d):", len(unexpectedDiagnostics))
		for _, diag := range unexpectedDiagnostics {
			t.Errorf("  %s", diag)
		}
	}
}

func TestVueTscBuildIndividual(t *testing.T) {
	t.Parallel()

	testWorkspacePath := getTestWorkspacePath()
	tscPath := filepath.Join(testWorkspacePath, "tsc")

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
		if skippedTests[testDir] {
			skippedCount++
		} else {
			enabledCount++
		}
	}
	t.Logf("Test directories: %d total, %d enabled, %d skipped", len(testDirs), enabledCount, skippedCount)

	for _, testDir := range testDirs {
		isExpectedFailure := strings.HasPrefix(testDir, "_failed_")

		t.Run(testDir, func(t *testing.T) {
			t.Parallel()

			if skippedTests[testDir] {
				t.Skipf("Test %s is in skippedTests list", testDir)
				return
			}

			testPath := filepath.Join(tscPath, testDir)
			sys := newVueTscSys(testPath)
			commandLineArgs := []string{"--build", ".", "--pretty", "false"}

			result := execute.CommandLine(sys, commandLineArgs, nil)
			output := sys.GetOutput()
			diagnostics := parseTscOutput(output)

			if isExpectedFailure {
				if len(diagnostics) == 0 {
					t.Logf("Expected failure test %s passed without errors (this might be OK if the issue is fixed)", testDir)
				} else {
					t.Logf("Expected failure test %s has %d diagnostics (expected)", testDir, len(diagnostics))
					for _, diag := range diagnostics {
						t.Logf("  %s", diag)
					}
				}
			} else {
				if len(diagnostics) > 0 {
					t.Errorf("Test %s failed with %d diagnostics:", testDir, len(diagnostics))
					for _, diag := range diagnostics {
						t.Errorf("  %s", diag)
					}
				}
				if result.Status != tsc.ExitStatusSuccess {
					if len(diagnostics) > 0 {
						t.Errorf("Test %s exited with status %v", testDir, result.Status)
					}
				}
			}
		})
	}
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

	sort.Strings(diagnostics)

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

	normalizedDiagnostics := make([]string, len(diagnostics))
	for i, diag := range diagnostics {
		normalizedDiagnostics[i] = normalizeDiagnosticPath(diag, testWorkspacePath)
	}
	sort.Strings(normalizedDiagnostics)

	var failedDiagnostics []string
	for _, diag := range normalizedDiagnostics {
		if strings.Contains(diag, "_failed_") {
			failedDiagnostics = append(failedDiagnostics, diag)
		}
	}
	sort.Strings(failedDiagnostics)
	sort.Strings(expectedDiagnostics)

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

func parseTscOutput(output string) []string {
	lines := strings.Split(output, "\n")
	var diagnostics []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(line, "): error TS") || strings.Contains(line, "): warning TS") {
			diagnostics = append(diagnostics, line)
		}
	}

	return diagnostics
}

func normalizeDiagnosticPath(diag string, basePath string) string {
	parenIdx := strings.Index(diag, "(")
	if parenIdx == -1 {
		return diag
	}

	filePath := diag[:parenIdx]
	rest := diag[parenIdx:]

	if strings.HasPrefix(filePath, basePath) {
		relPath, err := filepath.Rel(filepath.Dir(basePath), filePath)
		if err == nil {
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

	if _, err := os.Stat(tscPath); os.IsNotExist(err) {
		t.Fatalf("Test workspace not found at %s", tscPath)
	}

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
		if len(output) > 2000 {
			t.Logf("Build output (truncated):\n%s...", output[:2000])
		} else {
			t.Logf("Build output:\n%s", output)
		}
	}
}
