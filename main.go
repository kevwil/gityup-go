package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func handleTilde(path string) string {
	if strings.Contains(path, "~") {
		return strings.Replace(path, "~", os.Getenv("HOME"), 1)
	}
	return path
}

func parseArgs() string {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Usage: <app> <root-dir>")
		os.Exit(1)
	}
	argVal := handleTilde(flag.Arg(0))
	absDir, err := filepath.Abs(argVal)
	if err != nil {
		_ = fmt.Errorf("Error resolving absolute path to '%s': %v\n", argVal, err)
		os.Exit(1)
	}
	return absDir
}

func checkExecExists(executable string) (string, error) {
	if len(strings.TrimSpace(executable)) < 1 {
		return executable, errors.New("empty executable arg not allowed")
	}
	return exec.LookPath(executable)
}

func isGit(path string) bool {
	gitSubDir, err := filepath.Abs(filepath.Join(path, ".git"))
	if err != nil {
		_ = fmt.Errorf("Error resolving absolute path to '%s': %v\n", path, err)
	}
	info, err := os.Stat(gitSubDir)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			_ = fmt.Errorf("Error checking if git subdir exists: %v\n", err)
		}
	}
	return info.IsDir()
}

func gitStatus(dir string) bool {
	c1 := exec.Command("git", "status")
	c1.Dir = dir
	c2 := exec.Command("grep", "nothing to commit")
	c2.Stdin, _ = c1.StdoutPipe()
	c2.Stdout = io.Discard
	c2.Stderr = io.Discard
	err := c2.Start()
	if err != nil {
		_ = fmt.Errorf("error running grep: %v", err)
	}
	err = c1.Run()
	if err != nil {
		_ = fmt.Errorf("error running git status: %v", err)
	}
	err = c2.Wait()
	if err != nil {
		_ = fmt.Errorf("error waiting for grep to finish: %v", err)
	}
	return c2.ProcessState.Success()
}

func getBranchName(dir string) (string, error) {
	buf := bytes.NewBuffer(nil)
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = dir
	cmd.Stdout = buf
	cmd.Stderr = buf
	err := cmd.Run()
	return strings.TrimSpace(buf.String()), err
}

func gitRemote(dir, branchName string) bool {
	branchRemote := fmt.Sprintf("branch.%s.remote", branchName)
	cmd := exec.Command("git", "config", branchRemote)
	cmd.Dir = dir
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	err := cmd.Run()
	if err != nil {
		_ = fmt.Errorf("error checking git config: %v", err)
	}
	return cmd.ProcessState.Success()
}

func gitSync(dir string) {
	c1 := exec.Command("git", "smart-pull")
	c1.Dir = dir
	c1.Stdout = os.Stdout
	c1.Stderr = os.Stderr
	err := c1.Run()
	if err != nil {
		_ = fmt.Errorf("error running git smart-pull: %v", err)
	}
	//boolean and (&&)
	if c1.ProcessState.Success() {
		c2 := exec.Command("git", "remote", "update", "origin", "--prune")
		c2.Dir = dir
		c2.Stdout = os.Stdout
		c2.Stderr = os.Stderr
		err := c2.Run()
		if err != nil {
			_ = fmt.Errorf("error running git remote update: %v", err)
		}
	}
}

func updateProjects(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			fullFilePath := filepath.Join(path, file.Name())
			if isGit(fullFilePath) {
				if !gitStatus(fullFilePath) {
					fmt.Println("local changes detected, skipping", fullFilePath)
				} else {
					branchName, err := getBranchName(fullFilePath)
					if err != nil {
						_ = fmt.Errorf("error checking git branch name: %v", err)
					}
					if len(branchName) == 0 {
						fmt.Println("#### detached HEAD state, skipping ####")
					} else {
						if !gitRemote(fullFilePath, branchName) {
							fmt.Println("no remote to pull from, skipping ", fullFilePath)
						} else {
							fmt.Printf("#### pulling %s ####\n", fullFilePath)
							gitSync(fullFilePath)
							fmt.Println("")
						}
					}
				}
			}
		}
	}
}

func main() {
	rootDir := parseArgs()
	_, err := checkExecExists("git")
	if err != nil {
		log.Fatal(err)
	}
	_, err = checkExecExists("git-smart-pull")
	if err != nil {
		log.Fatal(err)
	}
	updateProjects(rootDir)
}
