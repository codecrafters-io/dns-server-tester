package internal

import (
	"github.com/codecrafters-io/tester-utils/tester_definition"
)

var testerDefinition = tester_definition.TesterDefinition{
	AntiCheatTestCases: []tester_definition.TestCase{},
	ExecutableFileName: "your_server.sh",
	TestCases: []tester_definition.TestCase{
		{
			Slug:     "init",
			TestFunc: testInit,
		},
		{
			Slug:     "write-headers",
			TestFunc: testReceiveEmptyResponse,
		},
		{
			Slug:     "write-question-section",
			TestFunc: testReceiveQuestionInResponse,
		},
		{
			Slug:     "write-answer-section",
			TestFunc: testReceiveAnswerInResponse,
		},
		{
			Slug:     "parse-headers",
			TestFunc: testHeaderParsing,
		},
		{
			Slug:     "parse-question",
			TestFunc: testBasicQuestionParsing,
		},
		{
			Slug:     "parse-compressed-packet",
			TestFunc: testCompressedPacketParsing,
		},
		{
			Slug:     "forwarding-server",
			TestFunc: testForwarding,
		},
	},
}
