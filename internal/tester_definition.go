package internal

import (
	testerutils "github.com/codecrafters-io/tester-utils"
)

var testerDefinition = testerutils.TesterDefinition{
	AntiCheatTestCases: []testerutils.TestCase{},
	ExecutableFileName: "your_server.sh",
	TestCases: []testerutils.TestCase{
		{
			Slug:     "init",
			TestFunc: testInitEmpty,
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
			Slug:     "parse-question-compressed",
			TestFunc: testCompressedQuestionParsing,
		},
		{
			Slug:     "forwarding-server",
			TestFunc: testForwarding,
		},
		{
			Slug:     "more-record-types",
			TestFunc: testMoreRecords,
		},
		{
			Slug:     "dns-resolution",
			TestFunc: testMoreRecords,
		},
	},
}
