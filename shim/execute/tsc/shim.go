package tsc

import (
	"io"
	"time"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/execute/incremental"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
	"github.com/microsoft/typescript-go/internal/golarext"
	"github.com/microsoft/typescript-go/internal/locale"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

// System interface for tsc execution
type System interface {
	GetGolarCallbacks() *golarext.GolarCallbacks

	Writer() io.Writer
	FS() vfs.FS
	DefaultLibraryPath() string
	GetCurrentDirectory() string
	WriteOutputIsTTY() bool
	GetWidthOfTerminal() int
	GetEnvironmentVariable(name string) string

	Now() time.Time
	SinceStart() time.Duration
}

// ExitStatus represents the exit status of the compiler
type ExitStatus = tsc.ExitStatus

const (
	ExitStatusSuccess                              = tsc.ExitStatusSuccess
	ExitStatusDiagnosticsPresent_OutputsGenerated  = tsc.ExitStatusDiagnosticsPresent_OutputsGenerated
	ExitStatusDiagnosticsPresent_OutputsSkipped    = tsc.ExitStatusDiagnosticsPresent_OutputsSkipped
	ExitStatusInvalidProject_OutputsSkipped        = tsc.ExitStatusInvalidProject_OutputsSkipped
	ExitStatusProjectReferenceCycle_OutputsSkipped = tsc.ExitStatusProjectReferenceCycle_OutputsSkipped
	ExitStatusNotImplemented                       = tsc.ExitStatusNotImplemented
)

// Watcher interface for watch mode
type Watcher = tsc.Watcher

// CommandLineResult represents the result of a command line execution
type CommandLineResult = tsc.CommandLineResult

// CommandLineTesting interface for testing hooks
type CommandLineTesting interface {
	OnEmittedFiles(result *compiler.EmitResult, mTimesCache *collections.SyncMap[tspath.Path, time.Time])
	OnListFilesStart(w io.Writer)
	OnListFilesEnd(w io.Writer)
	OnStatisticsStart(w io.Writer)
	OnStatisticsEnd(w io.Writer)
	OnBuildStatusReportStart(w io.Writer)
	OnBuildStatusReportEnd(w io.Writer)
	OnWatchStatusReportStart()
	OnWatchStatusReportEnd()
	GetTrace(w io.Writer, locale locale.Locale) func(msg *diagnostics.Message, args ...any)
	OnProgram(program *incremental.Program)
}

// CompileTimes contains timing information for compilation
type CompileTimes = tsc.CompileTimes

// CompileAndEmitResult contains the result of compilation
type CompileAndEmitResult struct {
	Diagnostics []*ast.Diagnostic
	EmitResult  *compiler.EmitResult
	Status      ExitStatus
}
