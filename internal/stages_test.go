package internal

import (
<<<<<<< HEAD
	tester_utils "github.com/codecrafters-io/tester-utils"
	"testing"
)

func TestStages(t *testing.T) {
	testCases := map[string]tester_utils.TesterOutputTestCase{
		"literal_character": {
			StageName:           "init",
=======
	"testing"

	tester_utils_testing "github.com/codecrafters-io/tester-utils/testing"
)

func TestStages(t *testing.T) {
	testCases := map[string]tester_utils_testing.TesterOutputTestCase{
		"literal_character": {
			UntilStageSlug:      "init",
>>>>>>> upstream/main
			CodePath:            "./test_helpers/scenarios/init/failure",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/init/failure",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
	}

<<<<<<< HEAD
	tester_utils.TestTesterOutput(t, testerDefinition, testCases)
=======
	tester_utils_testing.TestTesterOutput(t, testerDefinition, testCases)
>>>>>>> upstream/main
}

func normalizeTesterOutput(testerOutput []byte) []byte {
	return testerOutput
}
