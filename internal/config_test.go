package internal

import (
	"flag"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	type want struct {
		cfg Config
		err error
	}
	type test struct {
		name  string
		flags []string
		env   func()
		want  want
	}

	tests := []test{
		{
			name: "default read config without envs",
			flags: []string{
				"test",
				"--host", "23.233.43.9", "--port", "777", "--db", "mockDbDSN",
			},
			env: nil,
			want: want{
				cfg: Config{
					Host:      "23.233.43.9",
					Port:      777,
					Debug:     false,
					DBConnStr: "mockDbDSN",
				},
				err: nil,
			},
		},
		{
			name:  "default read config with envs",
			flags: []string{"test", "--debug"},
			env: func() {
				t.Setenv("NOTES_HOST", "73.133.73.97")
				t.Setenv("NOTES_PORT", "1111")
				t.Setenv("NOTES_DB", "mockDbDSN")
			},
			want: want{
				cfg: Config{
					Host:      "73.133.73.97",
					Port:      1111,
					Debug:     true,
					DBConnStr: "mockDbDSN",
				},
				err: nil,
			},
		},
		{
			name:  "call with bad port",
			flags: []string{"test", "--debug"},
			env: func() {
				t.Setenv("NOTES_HOST", "73.133.73.97")
				t.Setenv("NOTES_PORT", "abc")
				t.Setenv("NOTES_DB", "mockDbDSN")
			},
			want: want{
				cfg: Config{},
				err: strconv.ErrSyntax,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			oldsArgs := os.Args
			oldCommandLine := flag.CommandLine
			defer func() {
				os.Args = oldsArgs
				flag.CommandLine = oldCommandLine
			}()
			os.Args = tc.flags
			flag.CommandLine = flag.NewFlagSet(tc.flags[0], flag.ExitOnError)
			if tc.env != nil {
				tc.env()
			}

			cfg, err := ReadConfig()
			if tc.want.err != nil {
				assert.ErrorIs(t, err, tc.want.err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want.cfg, *cfg)
		})
	}
}
