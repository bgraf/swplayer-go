package main

import (
	"fmt"
	"github.com/SvantjeJung/swplayer/logfile"
	"github.com/spf13/pflag"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
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

	// Apply history from logfile and remove recently played titles.
	// If not enough titles remain after history application, we keep the original list.
	files = filterByPlayLog(files)

	if len(files) < numTitles {
		fmt.Fprintf(os.Stderr, "requested to play %d titles, but only %d available\n", numTitles, len(files))
		numTitles = len(files)
	}
	if numTitles == 0 {
		fmt.Fprintf(os.Stderr, "no files to play!\n")
		os.Exit(1)
	}

	shuffleStringSlice(files)
	playFiles(files[:numTitles])

	if !avoidShutdown {
		err := performShutdown()
		if err != nil {
			panic(err)
		}
	}
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


func playFiles(files []string) {
	for _, file := range files {
		// Write to play log file
		if playLog != nil {
			err := playLog.AppendTitle(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not write to logfile: %s", err)
			}
		}

		// Play file
		fmt.Println("Playing:", file)
		err := playFile(playerName, file)
		if err != nil {
			panic(err)
		}
	}
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

func filterByPlayLog(files []string) []string {
	if playLog == nil {
		return files
	}
	logEntries, err := playLog.ReadEntries()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read log entries: %s", err)
	}

	numEntries := len(logEntries)
	if histMax > numEntries {
		histMax = numEntries
	}
	logEntries = logEntries[numEntries-histMax : numEntries]

	// Create a map from string to empty-struct and use it as a set of history entries.
	history := make(map[string]struct{})
	for _, entry := range logEntries {
		history[entry] = struct{}{}
	}

	// Create new list of files without those in the current history.
	var remainingFiles []string
	for _, path := range files {
		if _, ok := history[path]; !ok {
			remainingFiles = append(remainingFiles, path)
		}
	}

	// Only use the non history files if there are sufficiently many.
	if len(remainingFiles) >= numTitles {
		files = remainingFiles
	} else {
		fmt.Fprintf(os.Stderr, "not enough titles available if history was applied, thus, the history is ignored.\n")
	}

	return files
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

func shuffleStringSlice(s []string) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}
