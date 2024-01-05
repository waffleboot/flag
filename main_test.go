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
			s.T().Log(file)

			cmdLine, output, status := s.parseTestData(file)

			result, err := exec.CommandContext(context.Background(), cmdLine[0], cmdLine[1:]...).Output()
			s.Equal(output, string(result))
			if err != nil {
				var errExit *exec.ExitError
				if errors.As(err, &errExit) {
					var exp int
					_, err := fmt.Sscanf(status, "%d", &exp)
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

func (s *FooSuite) parseTestData(file string) (cmdLine []string, output, status string) {
	a, err := txtar.ParseFile(file)
	s.NoError(err)
	s.Len(a.Files, 3)
	s.Equal("input", a.Files[0].Name)
	s.Equal("output", a.Files[1].Name)
	s.Equal("status", a.Files[2].Name)

	cmdLine = strings.Fields(string(a.Files[0].Data))
	s.True(len(cmdLine) > 0)

	return cmdLine, string(a.Files[1].Data), string(a.Files[2].Data)
}
