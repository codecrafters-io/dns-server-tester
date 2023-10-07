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
			TestFunc: testInit,
		},
		{
			Slug:     "receive-empty-response",
			TestFunc: testReceiveEmptyResponse,
		},
		{
			Slug:     "receive-question-in-response",
			TestFunc: testReceiveQuestionInResponse,
		},
		{
			Slug:     "receive-answer-in-response",
			TestFunc: testReceiveAnswerInResponse,
		},
		{
			Slug:     "parse-headers",
			TestFunc: testHeaderParsing,
		},
		{
			Slug:     "parse-question-basic",
			TestFunc: testBasicQuestionParsing,
		},
	},
}
