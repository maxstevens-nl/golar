package execute

import (
	"github.com/microsoft/typescript-go/internal/execute/tsc"
	_ "unsafe"
)

//go:linkname CommandLine github.com/microsoft/typescript-go/internal/execute.CommandLine
func CommandLine(sys tsc.System, commandLineArgs []string, testing tsc.CommandLineTesting) tsc.CommandLineResult
