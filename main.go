package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"time"
)

//go:embed usage/root.txt
var rootUsage string

//go:embed usage/run.txt
var runUsage string

//go:embed usage/test.txt
var testUsage string

const (
	CMDRUN  = "run"
	CMDTEST = "test"
	CMDHELP = "help"
)

func main() {
	var rootDir string

	// root flag set не имеет флагов, только один аргумент - выполняемую команду
	rootFS := flag.NewFlagSet("", flag.ExitOnError)
	rootFS.StringVar(&rootDir, "d", ".", "root dir")
	rootFS.StringVar(&rootDir, "dir", ".", "root dir")
	rootFS.SetOutput(os.Stdout)
	rootFS.Usage = func() {
		fmt.Println(rootUsage)
	}

	args := os.Args[1:]

	err := rootFS.Parse(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if rootFS.NArg() == 0 { // нет ни одной команды
		rootFS.Usage()
		os.Exit(2)
	}

	var runPort int

	runFS := flag.NewFlagSet(CMDRUN, rootFS.ErrorHandling())
	runFS.IntVar(&runPort, "p", 80, "http port")
	runFS.IntVar(&runPort, "port", 80, "http port")
	runFS.SetOutput(rootFS.Output())
	runFS.Usage = func() {
		fmt.Println(runUsage)
	}

	var testDuration time.Duration

	testFS := flag.NewFlagSet(CMDTEST, rootFS.ErrorHandling())
	testFS.DurationVar(&testDuration, "t", 1*time.Minute, "test duration")
	testFS.DurationVar(&testDuration, "time", 1*time.Minute, "test duration")
	testFS.SetOutput(rootFS.Output())
	testFS.Usage = func() {
		fmt.Println(testUsage)
	}

	switch cmd := rootFS.Arg(0); cmd {
	case CMDRUN:
		err := runFS.Parse(rootFS.Args()[1:])
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		run(rootDir, runPort)
	case CMDTEST:
		err = testFS.Parse(rootFS.Args()[1:])
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		test(rootDir, testDuration)
	case CMDHELP:
		if rootFS.NArg() > 1 {
			switch subcmd := rootFS.Arg(1); subcmd {
			case CMDRUN:
				runFS.Usage()
				os.Exit(0)
			case CMDTEST:
				testFS.Usage()
				os.Exit(0)
			default:
				fmt.Printf(`foo: '%s' is not a foo command. See 'foo --help'.`, subcmd)
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

func run(dir string, port int) {
	fmt.Println("dir is", dir)
	fmt.Println("http port is", port)
}

func test(dir string, duration time.Duration) {
	fmt.Println("dir is", dir)
	fmt.Println("test duration is", duration)
}
