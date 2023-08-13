package j2log

import (
	"bufio"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// the main struct for the application, hold the command-line options.
type J2Log struct {
	File *os.File `arg:"" default:"-" help:"The file to parse."`

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
	cli.prolouge()
	defer cli.epilogue()

	if err := cli.run(); err != nil {
		log.Error().Err(err).Msg("failed to run")
		return 1
	}

	return 0
}

// the exactly function to run the J2Log
func (cli *J2Log) run() (err error) {
	defer cli.File.Close()

	scanner := bufio.NewScanner(cli.File)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \t")
		log.Debug().Str("line", line).Msg("read line ...")
	}

	return
}

// setup the necessary before running
func (cli *J2Log) prolouge() {
	switch cli.Verbose {
	case -1:
		zerolog.SetGlobalLevel(zerolog.Disabled)
	case 0:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case 1:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case 2:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case 3:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// make loggers pretty
	writter := zerolog.ConsoleWriter{Out: os.Stderr}
	log.Logger = zerolog.New(writter).With().Timestamp().Logger()

	log.Info().Msg("finished prologue ...")
}

// clean up the necessary after runned
func (cli *J2Log) epilogue() {
	log.Info().Msg("starting epilogue ...")
	log.Info().Msg("finished epilogue ...")
}
