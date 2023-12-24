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
)

func TestFoo(t *testing.T) {
	files, err := filepath.Glob("testdata/**/*.txt")
	assert.NoError(t, err)

	for _, file := range files {
		file := file
		t.Run(filepath.Base(file), func(t *testing.T) {

			a, err := txtar.ParseFile(file)
			assert.NoError(t, err)
			assert.Len(t, a.Files, 3)
			assert.Equal(t, "input", a.Files[0].Name)
			assert.Equal(t, "output", a.Files[1].Name)
			assert.Equal(t, "status", a.Files[2].Name)

			cmdLine := strings.Fields(string(a.Files[0].Data))
			assert.True(t, len(cmdLine) > 0)

			output, err := exec.CommandContext(context.Background(), "./foo", cmdLine[1:]...).Output()
			assert.Equal(t, string(a.Files[1].Data), string(output))
			if err != nil {
				var errExit *exec.ExitError
				if errors.As(err, &errExit) {
					var exp int
					_, err := fmt.Sscanf(string(a.Files[2].Data), "%d", &exp)
					assert.NoError(t, err)
					assert.Equal(t, exp, errExit.ExitCode())
					return
				}
			}
		})
	}

}
