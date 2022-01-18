package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/ranna-go/cli/config"
	"github.com/ranna-go/ranna/pkg/client"
	"github.com/ranna-go/ranna/pkg/models"
)

func readFileOrStdin(path string) (code string, err error) {
	var data []byte

	if path != "" {
		data, err = os.ReadFile(path)
	} else {
		data, err = io.ReadAll(os.Stdin)
	}

	code = string(data)
	return
}

func vlogger(verbose bool) func(msg string, v ...interface{}) {
	if verbose {
		return func(msg string, v ...interface{}) {
			log.Printf("> "+msg, v...)
		}
	}
	return func(msg string, v ...interface{}) {}
}

func main() {
	log.SetFlags(0)

	godotenv.Load()
	cfg, err := config.Parse()
	if err != nil {
		log.Fatalf("ranna: failed parsing config: %s", err.Error())
	}

	var (
		filePath string
		spec     string
		inline   bool
		verbose  bool
	)

	flag.StringVar(&cfg.Endpoint, "endpoint", cfg.Endpoint, "ranna API endpoint")
	flag.StringVar(&cfg.Version, "version", cfg.Version, "ranna API version")
	flag.StringVar(&cfg.Authorization, "auth", cfg.Authorization, "ranna API auth token")
	flag.StringVar(&spec, "f", "", "code file to read (reads from stdin if not set)")
	flag.StringVar(&spec, "s", "", "spec to use")
	flag.BoolVar(&inline, "i", false, "Execute code inline")
	flag.BoolVar(&verbose, "v", false, "Verbose logging")
	flag.Parse()

	vlogf := vlogger(verbose)

	code, err := readFileOrStdin(filePath)
	if err != nil {
		log.Fatalf("ranna: failed to read code: %s", err.Error())
	}

	vlogf("endpoint: %s/%s", cfg.Endpoint, cfg.Version)
	vlogf("spec: %s", spec)
	vlogf("inline: %t", inline)

	cfg.Options.UserAgent = "ranna-cli"
	c, err := client.New(cfg.Options)
	if err != nil {
		log.Fatalf("ranna: failed creating ranna client: %s", err.Error())
	}

	vlogf("client created")

	if spec == "" {
		specMap, err := c.Spec()
		if err != nil {
			log.Fatalf("ranna: you must specify a spec")
		}

		var msg strings.Builder
		msg.WriteString("ranna: you must specify a spec\n")
		msg.WriteString("       available specs are:\n")
		for n, s := range specMap {
			msg.WriteString("         - " + n + ": ")
			if s.Use != "" {
				msg.WriteString("-> " + s.Use)
			} else {
				msg.WriteString("" + s.Image)
			}
			msg.WriteRune('\n')
		}

		log.Fatal(msg.String())
	}

	vlogf("executing code ...")

	res, err := c.Exec(models.ExecutionRequest{
		Language:         spec,
		Code:             code,
		InlineExpression: inline,
	})
	if err != nil {
		log.Fatalf("ranna: failed code execution: %s", err.Error())
	}

	vlogf("execution time: %s", (time.Duration(res.ExecTimeMS) * time.Millisecond).Round(time.Millisecond).String())
	os.Stderr.WriteString(res.StdErr)
	os.Stdout.WriteString(res.StdOut)
}
