package internal

import (
	testerutils "github.com/codecrafters-io/tester-utils"
)

var testerDefinition = testerutils.TesterDefinition{
<<<<<<< HEAD
	AntiCheatStages:    []testerutils.Stage{},
	ExecutableFileName: "script.sh",
	Stages: []testerutils.Stage{
		{
			Number:                  1,
			Slug:                    "init",
			Title:                   "Match a literal character",
			TestFunc:                testInit,
			ShouldRunPreviousStages: true,
=======
	AntiCheatTestCases:    []testerutils.TestCase{},
	ExecutableFileName: "script.sh",
	TestCases: []testerutils.TestCase{
		{
			Slug:                    "init",
			TestFunc:                testInit,
>>>>>>> upstream/main
		},
	},
}
