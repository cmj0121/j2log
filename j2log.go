package j2log

import (
	"os"

	"github.com/alecthomas/kong"
)

// the main struct for the application, hold the command-line options.
type J2Log struct {
	File *os.File `arg:"" default:"-" help:"The file to parse."`

	// the logger options
	Quiet   bool `short:"q" group:"logger" xor:"verbose,quiet" help:"Disable all logger."`
	Verbose int  `short:"v" group:"logger" xor:"verbose,quiet" type:"counter" help:"Show the verbose logger."`
}

// create a new J2Log struct with default settings, and return a pointer to it.
func New() *J2Log {
	return &J2Log{}
}

// the main function, called from main.go, parses the command-line options and
// calls the appropriate function.
func (cli *J2Log) ParseAndRun() {
	kong.Parse(cli)
	cli.RunAndExit()
}

// execute the J2Log and exit with the appropriate exit code.
func (cli *J2Log) RunAndExit() {
	code := cli.Run()
	os.Exit(code)
}

// execute the J2Log and return the appropriate exit code.
func (cli *J2Log) Run() int {
	return 0
}
