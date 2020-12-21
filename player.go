package main

import (
	"os"
	"os/exec"
)

func checkPlayerAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func playFile(player string, path string) error {
	cmd := exec.Command(player, path)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	return err
}

