// Package main is the gityup tool. It is a simple tool to update a folder of Git repositories.
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

// ErrEmptyExec is thrown when the executable arg to checkExecExists is empty.
var ErrEmptyExec = errors.New("empty executable arg not allowed")

func handleTilde(path string) string {
	if strings.Contains(path, "~") {
		return strings.Replace(path, "~", os.Getenv("HOME"), 1)
	}

	return path
}

func parseArgs() string {
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("Usage: <app> <root-dir>")
	}

	argVal := handleTilde(flag.Arg(0))
	absDir, err := filepath.Abs(argVal)

	if err != nil {
		log.Fatalf("failed to determine absolute path of %s: %s", argVal, err)
	}

	return absDir
}

func checkExecExists(executable string) (string, error) {
	if len(strings.TrimSpace(executable)) < 1 {
		return "", ErrEmptyExec
	}

	path, err := exec.LookPath(executable)
	if err != nil {
		return path, fmt.Errorf("error looking up executable '%s': %w", executable, err)
	}

	return path, nil
}

func isGit(path string) bool {
	// ensure subdirectory ".git" exists in project
	gitSubDir, err := filepath.Abs(filepath.Join(path, ".git"))
	if err != nil {
		_ = fmt.Errorf("error resolving absolute path to '%s': %w", path, err)
	}

	info, err := os.Stat(gitSubDir)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}

		_ = fmt.Errorf("error checking if git subdir exists: %w", err)
	}

	return info.IsDir()
}

func gitStatus(dir string) bool {
	command1 := exec.Command("git", "status")
	command1.Dir = dir
	command2 := exec.Command("grep", "nothing to commit")
	command2.Stdin, _ = command1.StdoutPipe()
	command2.Stdout = io.Discard
	command2.Stderr = io.Discard
	err := command2.Start()

	if err != nil {
		_ = fmt.Errorf("error running grep: %w", err)
	}

	err = command1.Run()
	if err != nil {
		_ = fmt.Errorf("error running git status: %w", err)
	}

	err = command2.Wait()
	if err != nil {
		_ = fmt.Errorf("error waiting for grep to finish: %w", err)
	}

	return command2.ProcessState.Success()
}

func getBranchName(dir string) (string, error) {
	// buffer to save output
	buf := bytes.NewBuffer(nil)
	cmd := exec.Command("git", "branch", "--show-current") // #nosec G204
	cmd.Dir = dir
	// output to buffer
	cmd.Stdout = buf
	cmd.Stderr = buf
	err := cmd.Run()
	// return output as string
	return strings.TrimSpace(buf.String()), err
}

func gitRemote(dir, branchName string) bool {
	branchRemote := fmt.Sprintf("branch.%s.remote", branchName)
	cmd := exec.Command("git", "config", branchRemote) // #nosec G204
	cmd.Dir = dir
	cmd.Stdout = io.Discard // >/dev/null
	cmd.Stderr = io.Discard // 2>&1
	err := cmd.Run()

	if err != nil {
		_ = fmt.Errorf("error checking git config: %w", err)
	}

	return cmd.ProcessState.Success()
}

func gitSync(dir string) {
	command1 := exec.Command("git", "smart-pull")
	command1.Dir = dir
	// show output in shell
	command1.Stdout = os.Stdout
	command1.Stderr = os.Stderr
	err := command1.Run()

	if err != nil {
		_ = fmt.Errorf("error running git smart-pull: %w", err)
	}
	// boolean and (&&)
	if command1.ProcessState.Success() {
		command2 := exec.Command("git", "remote", "update", "origin", "--prune")
		command2.Dir = dir
		command2.Stdout = os.Stdout
		command2.Stderr = os.Stderr
		err := command2.Run()

		if err != nil {
			_ = fmt.Errorf("error running git remote update: %w", err)
		}
	}
}

func updateProjects(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() { //nolint
			fullFilePath := filepath.Join(path, file.Name())
			if isGit(fullFilePath) {
				if !gitStatus(fullFilePath) {
					fmt.Println("local changes detected, skipping", fullFilePath)
				} else {
					branchName, err := getBranchName(fullFilePath)
					if err != nil {
						_ = fmt.Errorf("error checking git branch name: %w", err)
					}

					if len(branchName) == 0 {
						fmt.Println("#### detached HEAD state, skipping ####")
					} else {
						if !gitRemote(fullFilePath, branchName) {
							fmt.Println("no remote to pull from, skipping ", fullFilePath)
						} else {
							fmt.Println("#### pulling", fullFilePath, "####")
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
