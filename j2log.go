package j2log

import (
	"regexp"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	PROJ_NAME = "j2log"
	PROJ_DESC = "A tool to convert JSON log to human-readable log."

	MAJOR = 0
	MINOR = 1
	MACRO = 0
)

var (
	// the docker-compose log format
	RE_JSON_WITH_NAME = regexp.MustCompile(`^([a-zA-Z0-9_]+)\s*\|\s*(\{.*\})$`)
)

// the main struct for the application, hold the command-line options.
type J2Log struct {
	File *os.File `arg:"" default:"-" help:"The file to parse."`

	// show version and exit
	Version kong.VersionFlag `short:"V" name:"version" help:"Print version info and quit"`

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
	kong.Parse(
		cli,
		kong.Name(PROJ_NAME),
		kong.Description(PROJ_DESC),
		kong.Vars{
			"version": fmt.Sprintf("%v (v%d.%d.%d)", PROJ_NAME, MAJOR, MINOR, MACRO),
		},
	)

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

	tmpl := DefaultTmpl()

	scanner := bufio.NewScanner(cli.File)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \t")
		log.Debug().Str("line", line).Msg("read line ...")

		switch encoded_line, ok := cli.trans(line, tmpl); ok {
		case true:
			fmt.Println(encoded_line)
		case false:
			log.Warn().Str("line", line).Msg("cannot translate line")
		}
	}

	return
}

// convert encoded JSON data to human-readable log
func (cli *J2Log) trans(raw string, tmpl *Template) (line string, ok bool) {
	switch data, err := cli.unmarchal(raw); err {
	case nil:
		line, ok = tmpl.Extract(data)
	default:
		log.Debug().Err(err).Msg("failed to unmarshal from JSON")
	}

	return
}

// unmarchal the JSON string to JSON object
func (cli *J2Log) unmarchal(raw string) (data map[string]interface{}, err error) {
	if err = json.Unmarshal([]byte(raw), &data); err != nil {
		// give second chance that the raw string is `NAME | { ... }`
		if !RE_JSON_WITH_NAME.MatchString(raw) {
			log.Debug().Err(err).Msg("failed to unmarshal from JSON")
			return
		}

		raw = RE_JSON_WITH_NAME.FindStringSubmatch(raw)[2]
		raw = strings.Trim(raw, " \t")

		log.Debug().Str("raw", raw).Msg("try to unmarshal again ...")
		data, err = cli.unmarchal(raw)
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
