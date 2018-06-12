package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// declare sub test setup func
type SetupSubTest func(t *testing.T) func(t *testing.T)

// This is the directory where our test fixtures are.
const fixtureDir = "./text-fixtures"

func TestGen(t *testing.T) {

	tt := []struct {
		cfg          *config
		file         string
		setupSubTest SetupSubTest
	}{
		{
			cfg: &config{
				typeNames:  "ErrorCode",
				apiVersion: "1.0",
			},
			file: "custom_error1_custom_string",
			setupSubTest: func(t *testing.T) func(t *testing.T) {
				os.Chdir(filepath.Join(fixtureDir, "custom_error1"))
				return func(t *testing.T) {
					os.Chdir("../..")
				}
			},
		},
		{
			cfg: &config{
				typeNames:  "ErrorCode",
				apiVersion: "1.1",
			},
			file: "custom_error2_custom_string",
			setupSubTest: func(t *testing.T) func(t *testing.T) {
				os.Chdir(filepath.Join(fixtureDir, "custom_error2"))
				return func(t *testing.T) {
					os.Chdir("../..")
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.file, func(t *testing.T) {
			teardownSubTest := tc.setupSubTest(t)
			defer teardownSubTest(t)

			for _, v := range tc.cfg.ParsePackage() {
				given := tc.cfg.genString(v)

				want, err := ioutil.ReadFile(fmt.Sprintf("%s.golden", tc.file))
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, given, want)
			}
		})
	}
}
