package internal

import (
	"github.com/codecrafters-io/tester-utils/tester_definition"
)

var testerDefinition = tester_definition.TesterDefinition{
	AntiCheatTestCases:       []tester_definition.TestCase{},
	ExecutableFileName:       "your_program.sh",
	LegacyExecutableFileName: "your_server.sh",
	TestCases: []tester_definition.TestCase{
		{
			Slug:     "ux2",
			TestFunc: testInit,
		},
		{
			Slug:     "tz1",
			TestFunc: testReceiveEmptyResponse,
		},
		{
			Slug:     "bf2",
			TestFunc: testReceiveQuestionInResponse,
		},
		{
			Slug:     "xm2",
			TestFunc: testReceiveAnswerInResponse,
		},
		{
			Slug:     "uc8",
			TestFunc: testHeaderParsing,
		},
		{
			Slug:     "hd8",
			TestFunc: testBasicQuestionParsing,
		},
		{
			Slug:     "yc9",
			TestFunc: testCompressedPacketParsing,
		},
		{
			Slug:     "gt1",
			TestFunc: testForwarding,
		},
	},
}
