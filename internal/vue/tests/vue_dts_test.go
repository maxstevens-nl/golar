package vue_tests

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/auvred/golar/golar"
	"github.com/microsoft/typescript-go/shim/bundled"
	"github.com/microsoft/typescript-go/shim/compiler"
	"github.com/microsoft/typescript-go/shim/core"
	"github.com/microsoft/typescript-go/shim/tsoptions"
	"github.com/microsoft/typescript-go/shim/tspath"
	"github.com/microsoft/typescript-go/shim/vfs/osvfs"
	"gotest.tools/v3/assert"
)

func TestVueDTS(t *testing.T) {
	t.Parallel()

	testWorkspacePath := getTestWorkspacePath()
	componentMetaPath := filepath.Join(testWorkspacePath, "component-meta")
	testdataPath := filepath.Join(filepath.Dir(testWorkspacePath), "internal", "vue", "tests", "testdata")

	inputFiles := readDTSFilesRecursive(componentMetaPath, componentMetaPath)
	if len(inputFiles) == 0 {
		t.Fatal("No input files found")
	}

	rootDir := tspath.NormalizePath(componentMetaPath)
	compilerOptions := &core.CompilerOptions{
		RootDir:              rootDir,
		Declaration:          core.TSTrue,
		EmitDeclarationOnly:  core.TSTrue,
		AllowNonTsExtensions: core.TSTrue,
	}

	fs := golar.WrapFS(bundled.WrapFS(osvfs.FS()))

	baseHost := compiler.NewCompilerHost(
		rootDir,
		fs,
		bundled.LibPath(),
		nil,
		nil,
	)
	host := golar.GolarExtCallbacks.WrapCompilerHost(baseHost)

	config := tsoptions.NewParsedCommandLine(
		compilerOptions,
		inputFiles,
		tspath.ComparePathsOptions{
			CurrentDirectory:          rootDir,
			UseCaseSensitiveFileNames: fs.UseCaseSensitiveFileNames(),
		},
	)

	program := compiler.NewProgram(compiler.ProgramOptions{
		Host:   host,
		Config: config,
	})

	for _, inputFile := range inputFiles {
		expectedOutputFile := getExpectedDTSOutputFile(inputFile)
		testName := shortenDTSPath(inputFile) + " -> " + shortenDTSPath(expectedOutputFile)

		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			sourceFile := program.GetSourceFile(inputFile)
			if sourceFile == nil {
				t.Skipf("Could not get source file for %s", inputFile)
				return
			}

			var outputText string
			var outputFile string

			ctx := context.Background()
			program.Emit(ctx, compiler.EmitOptions{
				TargetSourceFile: sourceFile,
				EmitOnly:         compiler.EmitOnlyDts,
				WriteFile: func(fileName string, text string, writeByteOrderMark bool, data *compiler.WriteFileData) error {
					outputFile = fileName
					outputText = text
					return nil
				},
			})

			relPath, err := filepath.Rel(componentMetaPath, expectedOutputFile)
			if err != nil {
				t.Fatalf("Failed to get relative path: %v", err)
			}
			snapshotFile := filepath.Join(testdataPath, relPath)

			expectedContent, err := os.ReadFile(snapshotFile)
			if err != nil {
				if os.IsNotExist(err) {
					t.Skipf("Snapshot file not found: %s", snapshotFile)
					return
				}
				t.Fatalf("Failed to read snapshot: %v", err)
			}

			expected := normalizeDTSNewlines(string(expectedContent))
			actual := normalizeDTSNewlines(outputText)

			if outputFile != "" {
				normalizedOutputFile := normalizeDTSPath(outputFile)
				normalizedExpectedOutputFile := normalizeDTSPath(expectedOutputFile)
				assert.Equal(t, normalizedExpectedOutputFile, normalizedOutputFile, "Output file path mismatch")
			}
			assert.Equal(t, expected, actual, "Output content mismatch for %s", inputFile)
		})
	}
}

func readDTSFilesRecursive(dir string, workspace string) []string {
	relDir, _ := filepath.Rel(workspace, dir)
	if strings.HasPrefix(relDir, "#") || strings.HasPrefix(filepath.Base(dir), "#") {
		return nil
	}

	var result []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return result
	}

	for _, entry := range entries {
		name := entry.Name()

		if name == "tsconfig.json" {
			continue
		}

		// TODO: add support for custom extensions
		if strings.HasSuffix(name, ".cext") {
			continue
		}

		fullPath := filepath.Join(dir, name)

		if entry.IsDir() {
			if strings.HasPrefix(name, "#") {
				continue
			}
			result = append(result, readDTSFilesRecursive(fullPath, workspace)...)
		} else {
			result = append(result, tspath.NormalizePath(fullPath))
		}
	}

	return result
}

func getExpectedDTSOutputFile(inputFile string) string {
	if strings.HasSuffix(inputFile, ".ts") && !strings.HasSuffix(inputFile, ".d.ts") {
		return inputFile[:len(inputFile)-3] + ".d.ts"
	}
	if strings.HasSuffix(inputFile, ".tsx") {
		return inputFile[:len(inputFile)-4] + ".d.ts"
	}
	return inputFile + ".d.ts"
}

func shortenDTSPath(path string) string {
	path = normalizeDTSPath(path)
	segments := strings.Split(path, "/")
	if len(segments) <= 2 {
		return path
	}
	return strings.Join(segments[len(segments)-2:], "/")
}

func normalizeDTSPath(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

func normalizeDTSNewlines(text string) string {
	return strings.ReplaceAll(text, "\r\n", "\n")
}
