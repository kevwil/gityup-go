package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MainTestSuite struct {
	suite.Suite
}

func (suite *MainTestSuite) TestHandleTilde() {
	suite.Equal(filepath.Join(os.Getenv("HOME"), "Downloads"), handleTilde("~/Downloads"))
}

func (suite *MainTestSuite) TestHandleTilde_NoPath() {
	suite.Empty(handleTilde(""))
}

func (suite *MainTestSuite) TestHandleTilde_NoTilde() {
	suite.Equal("/usr/local", handleTilde("/usr/local"))
}

func (suite *MainTestSuite) TestExecExists() {
	path, err := checkExecExists("bash")
	suite.Require().NoError(err, "received unexpected error")
	suite.NotEmpty(path, "result should be full path to executable")

	possibleBashPaths := []string{"/bin/bash", "/usr/bin/bash", "/usr/local/bin/bash"}
	// suite.Equal("/bin/bash", path)
	suite.Contains(possibleBashPaths, path)
}

func (suite *MainTestSuite) TestExecExists_NoPath() {
	path, err := checkExecExists("")
	suite.Empty(path, "expected exec full path to be empty")
	suite.Error(err, "expected an error")
}

func (suite *MainTestSuite) TestIsGit() {
	suite.True(isGit(handleTilde(".")))
}

func (suite *MainTestSuite) TestIsGit_False() {
	suite.False(isGit(handleTilde("~/Downloads")))
}

func (suite *MainTestSuite) TestGitStatus() {
	suite.True(gitStatus(handleTilde(".")))
}

func (suite *MainTestSuite) TestGitStatus_False() {
	suite.False(gitStatus(handleTilde("~/Downloads")))
}

func (suite *MainTestSuite) TestGetBranchName() {
	branch, err := getBranchName(handleTilde("."))
	suite.Require().NoError(err, "received unexpected error")
	suite.NotEmpty(branch, "branch name should not be empty")
	suite.Equal("main", branch, "expected main branch")
}

func TestMainSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}
