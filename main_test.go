package main

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rogpeppe/go-internal/txtar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FooSuite struct {
	suite.Suite
}

func TestFooSuite(t *testing.T) {
	suite.Run(t, &FooSuite{})
}

func (s *FooSuite) TestFoo() {
	for _, file := range s.testDataFiles() {
		file := file
		s.Run(filepath.Base(file), func() {

			a, err := txtar.ParseFile(file)
			s.NoError(err)
			s.Len(a.Files, 3)
			s.Equal("input", a.Files[0].Name)
			s.Equal("output", a.Files[1].Name)
			s.Equal("status", a.Files[2].Name)

			cmdLine := strings.Fields(string(a.Files[0].Data))
			s.True(len(cmdLine) > 0)

			output, err := exec.CommandContext(context.Background(), "./foo", cmdLine[1:]...).Output()
			s.Equal(string(a.Files[1].Data), string(output))
			if err != nil {
				var errExit *exec.ExitError
				if errors.As(err, &errExit) {
					var exp int
					_, err := fmt.Sscanf(string(a.Files[2].Data), "%d", &exp)
					s.NoError(err)
					s.Equal(exp, errExit.ExitCode())
					return
				}
			}
		})
	}

}

func (s *FooSuite) testDataFiles() []string {
	files, err := filepath.Glob("testdata/**/*.txt")
	assert.NoError(s.T(), err)
	return files
}
