package protocmp

import (
	"fmt"
	"github.com/nbaztec/protocmp/cmp"
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

type suite struct {
	ProtocVersion      string
	ProtocGenGoVersion string
}

var tests = []struct{
	name string
	fn   func(*testing.T)
} {
	{
		name: "assert equal",
		fn:   cmp.TestAssertEqual,
	},
	{
		name: "string",
		fn:   cmp.TestAssertString,
	},
	{
		name: "int",
		fn:   cmp.TestAssertInt,
	},
	{
		name: "bool",
		fn:   cmp.TestAssertBool,
	},
	{
		name: "double",
		fn:   cmp.TestAssertDouble,
	},
	{
		name: "bytes",
		fn:   cmp.TestAssertBytes,
	},
	{
		name: "repeated",
		fn:   cmp.TestAssertRepeated,
	},
	{
		name: "repeated simple",
		fn:   cmp.TestAssertRepeatedSimple,
	},
}

func TestVersions(t *testing.T) {
	suites := []suite{
		{
			ProtocVersion:      "3.5.1",
			ProtocGenGoVersion: "1.2.0",
		},
		{
			ProtocVersion:      "3.5.1",
			ProtocGenGoVersion: "1.4.2",
		},
		{
			ProtocVersion:      "3.12.4",
			ProtocGenGoVersion: "1.4.2",
		},
	}

	for _, tt := range suites {
		tt := tt
		t.Run(fmt.Sprintf("protoc@v%s + protoc_gen_go@v%s", tt.ProtocVersion, tt.ProtocGenGoVersion), func(t *testing.T) {
			cmd := exec.Command("make", "clean", "protos")
			cmd.Env = append(
				cmd.Env,
				"PATH="+os.Getenv("PATH"),
				"HOME="+os.Getenv("HOME"),
				"GOPATH="+os.Getenv("GOPATH"),
				"GOCACHE="+os.Getenv("GOCACHE"),
				"PATH="+os.Getenv("PATH"),
				"PROTOC_VERSION="+tt.ProtocVersion,
				"PROTOC_GEN_GO="+tt.ProtocGenGoVersion,
			)
			cmd.Stdin = os.Stdin
			cmd.Stdout = nil
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				t.Error(err)
			}

			for _, test := range tests {
				test := test
				t.Run(test.name, func(t *testing.T) {
					test.fn(t)
				})
			}
		})
	}
}
