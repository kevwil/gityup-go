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
	suite.Equal("", handleTilde(""))
}

func (suite *MainTestSuite) TestHandleTilde_NoTilde() {
	suite.Equal("/usr/local", handleTilde("/usr/local"))
}

func (suite *MainTestSuite) TestExecExists() {
	s, err := checkExecExists("bash")
	suite.NoError(err, "received unexpected error")
	suite.NotEmpty(s, "result should be full path to executable")
	suite.Equal("/bin/bash", s)
}

func (suite *MainTestSuite) TestExecExists_NoPath() {
	s, err := checkExecExists("")
	suite.Empty(s, "exec full path empty")
	suite.Error(err, "expected an error")
}

func (suite *MainTestSuite) TestIsGit() {
	suite.True(isGit(handleTilde("~/dotfiles")))
}

func (suite *MainTestSuite) TestGitStatus() {
	suite.True(gitStatus(handleTilde("~/dotfiles")))
}

func (suite *MainTestSuite) TestGetBranchName() {
	branch, err := getBranchName(handleTilde("~/dotfiles"))
	suite.NoError(err, "received unexpected error")
	suite.NotEmpty(branch, "branch name should not be empty")
	suite.Equal("master", branch, "expected master branch")
}

func TestMainSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}
