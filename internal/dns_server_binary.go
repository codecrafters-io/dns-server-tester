package internal

import (
	executable "github.com/codecrafters-io/tester-utils/executable"
	logger "github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

type DnsServerBinary struct {
	executable *executable.Executable
	logger     *logger.Logger
}

func NewDnsServerBinary(stageHarness *test_case_harness.TestCaseHarness) *DnsServerBinary {
	b := &DnsServerBinary{
		executable: stageHarness.Executable,
		logger:     stageHarness.Logger,
	}

	stageHarness.RegisterTeardownFunc(func() { b.Kill() })

	return b
}

func (b *DnsServerBinary) Run(args ...string) error {
	b.logger.Debugf("Running program")
	if err := b.executable.Start(args...); err != nil {
		return err
	}

	return nil
}

func (b *DnsServerBinary) HasExited() bool {
	return b.executable.HasExited()
}

func (b *DnsServerBinary) Kill() error {
	if b.HasExited() {
		return nil
	}

	b.logger.Debugf("Terminating program")
	if err := b.executable.Kill(); err != nil {
		b.logger.Debugf("Error terminating program: '%v'", err)
		return err
	}

	b.logger.Debugf("Program terminated successfully")
	return nil // When does this happen?
}
