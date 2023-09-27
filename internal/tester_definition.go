package internal

import (
	testerutils "github.com/codecrafters-io/tester-utils"
)

var testerDefinition = testerutils.TesterDefinition{
	AntiCheatTestCases:    []testerutils.TestCase{},
	ExecutableFileName: "script.sh",
	TestCases: []testerutils.TestCase{
		{
			Slug:                    "init",
			TestFunc:                testInit,
		},
	},
}
