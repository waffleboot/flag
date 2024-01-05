package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"time"
)

//go:embed usage.txt
var usage string

//go:embed test.txt
var test string

//go:embed run.txt
var run string

const (
	CMDRUN  = "run"
	CMDTEST = "test"
	CMDHELP = "help"
)

func main() {
	rootFS := flag.NewFlagSet("", flag.ContinueOnError)
	rootFS.Usage = func() {
		fmt.Println(usage)
	}

	err := rootFS.Parse(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if rootFS.NArg() == 0 {
		rootFS.Usage()
		os.Exit(2)
	}

	var port int
	var testDuration time.Duration

	runFS := flag.NewFlagSet("run", flag.ContinueOnError)
	runFS.IntVar(&port, "p", 80, "http port")
	runFS.IntVar(&port, "port", 80, "http port")
	runFS.Usage = func() {
		fmt.Println(run)
	}

	testFS := flag.NewFlagSet("test", flag.ContinueOnError)
	testFS.DurationVar(&testDuration, "t", 1*time.Minute, "test duration")
	testFS.DurationVar(&testDuration, "time", 1*time.Minute, "test duration")
	testFS.Usage = func() {
		fmt.Println(test)
	}

	switch cmd := rootFS.Arg(0); cmd {
	case CMDRUN:
		err := runFS.Parse(rootFS.Args()[1:])
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Println("http port is", port)
	case CMDTEST:
		err = testFS.Parse(rootFS.Args()[1:])
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Println(testDuration)
	case CMDHELP:
		if rootFS.NArg() > 1 {
			switch cmd = rootFS.Arg(1); cmd {
			case CMDRUN:
				runFS.Usage()
				os.Exit(0)
			case CMDTEST:
				testFS.Usage()
				os.Exit(0)
			default:
				fmt.Printf(`foo: '%s' is not a foo command. See 'foo --help'.`, cmd)
				fmt.Println()
				os.Exit(2)
			}
		}

		rootFS.Usage()
		os.Exit(0)
	default:
		fmt.Printf(`foo: '%s' is not a foo command. See 'foo --help'.`, cmd)
		fmt.Println()
		os.Exit(2)
	}
}
