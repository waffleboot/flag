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
	f := flag.NewFlagSet("", flag.ContinueOnError)
	f.Usage = func() {
		fmt.Println(usage)
	}

	err := f.Parse(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if f.NArg() == 0 {
		f.Usage()
		os.Exit(2)
	}

	var port int
	var testDuration time.Duration

	f1 := flag.NewFlagSet("run", flag.ContinueOnError)
	f1.IntVar(&port, "p", 80, "http port")
	f1.IntVar(&port, "port", 80, "http port")
	f1.Usage = func() {
		fmt.Println(run)
	}

	f2 := flag.NewFlagSet("test", flag.ContinueOnError)
	f2.DurationVar(&testDuration, "t", 1*time.Minute, "test duration")
	f2.DurationVar(&testDuration, "time", 1*time.Minute, "test duration")
	f2.Usage = func() {
		fmt.Println(test)
	}

	switch cmd := f.Arg(0); cmd {
	case CMDRUN:
		err := f1.Parse(f.Args()[1:])
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Println("http port is", port)
	case CMDTEST:
		err = f2.Parse(f.Args()[1:])
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Println(testDuration)
	case CMDHELP:
		if f.NArg() > 1 {
			switch cmd = f.Arg(1); cmd {
			case CMDRUN:
				f1.Usage()
				os.Exit(0)
			case CMDTEST:
				f2.Usage()
				os.Exit(0)
			default:
				fmt.Printf(`foo: '%s' is not a foo command. See 'foo --help'.`, cmd)
				fmt.Println()
				os.Exit(2)
			}
		}

		f.Usage()
		os.Exit(0)
	default:
		fmt.Printf(`foo: '%s' is not a foo command. See 'foo --help'.`, cmd)
		fmt.Println()
		os.Exit(2)
	}
}
