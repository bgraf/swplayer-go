package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/SvantjeJung/swplayer/logfile"
	"github.com/spf13/pflag"
)

var numTitles int = 1
var histMax int = 100
var avoidShutdown bool = false
var fileExtensions []string
var playerName string

var playLog *logfile.Logfile = nil

func main() {
	parseProgramArguments()
	setupPlayLog()
	ensurePlayerAvailable()
	convertExtentionsLowercase()
	searchPaths := getSearchPaths()

	// Collect files from search paths
	files, err := collectFiles(searchPaths, fileExtensions)
	if err != nil {
		panic(err)
	}

	for i := 0; i < numTitles; i++ {
		performSingleFile(files)
	}

	if !avoidShutdown {
		err := performShutdown()
		if err != nil {
			panic(err)
		}
	}
}

func performSingleFile(files []string) error {
	lg := logfile.NewDefaultLogfile()

	historyEntries, err := lg.ReadEntries()
	if err != nil {
		fmt.Printf("no history")
	}

	nextFile, err := chooseFile(files, historyEntries)
	if err != nil {
		return err
	}

	if err := lg.AppendTitle(nextFile); err != nil {
		return err
	}

	fmt.Println("Playing:", nextFile)
	if err := playFile(playerName, nextFile); err != nil {
		return err

	}

	return nil
}

func performShutdown() error {
	cmd := exec.Command("shutdown", "-h", "now")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	return err
}

func parseProgramArguments() {
	pflag.IntVarP(&numTitles, "num", "n", 1, "number of titles to play")
	pflag.IntVarP(&histMax, "history", "l", 100, "maximum number of titles to ignore from history")
	pflag.BoolVar(&avoidShutdown, "no-shutdown", false, "...")
	pflag.StringSliceVarP(
		&fileExtensions,
		"formats",
		"f",
		[]string{"mp3", "mp4", "m4a", "flac", "webm"},
		"file extensions",
	)
	pflag.StringVarP(&playerName, "player", "p", "mpv", "player name")
	pflag.Parse()
}

func setupPlayLog() {
	playLog = new(logfile.Logfile)
	*playLog = logfile.NewDefaultLogfile()
	err := playLog.EnsureDirectory()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not ensure directory of logfile:\n  %s\n", playLog.Path())
		playLog = nil
	}
}

func ensurePlayerAvailable() {
	if !checkPlayerAvailable(playerName) {
		fmt.Fprintf(os.Stderr, "player '%s' is not available\n", playerName)
		os.Exit(1)
	}
}

func convertExtentionsLowercase() {
	// Convert all file extensions to lowercase
	for i, ext := range fileExtensions {
		fileExtensions[i] = strings.ToLower(ext)
	}
}

func getSearchPaths() []string {
	// Add search path "." if none were given
	searchPaths := pflag.Args()
	if len(searchPaths) == 0 {
		searchPaths = append(searchPaths, ".")
	}
	return searchPaths
}
