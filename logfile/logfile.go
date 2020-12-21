package logfile

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Wessie/appdirs"
)

// Logfile provides an interface to an on-disk logfile that keeps track of the titles
// that have been played.
type Logfile struct {
	path string
}

// NewLogfile creates a new logfile at the specified path.
func NewLogfile(path string) Logfile {
	return Logfile{
		path: path,
	}
}

// NewDefaultLogfile creates a new logfile at the XDG_USER_CACHE default location.
// Usually that is "~/.cache/swplayer/<version>/history.log".
func NewDefaultLogfile() Logfile {
	dir := appdirs.UserCacheDir("swplayer", "takaputo", "0.0.1", false)
	historyPath := filepath.Join(dir, "history.log")
	return NewLogfile(historyPath)
}

// AppendTitle appends the absolute path of the given file path to the logfile.
func (lf *Logfile) AppendTitle(file string) error {
	file, err := filepath.Abs(file)
	if err != nil {
		return fmt.Errorf("AppendTitle: %w", err)
	}

	f, err := os.OpenFile(lf.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("AppendTitle: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("%s\n", file)); err != nil {
		return fmt.Errorf("AppendTitle: %w", err)
	}

	return nil
}

// Path of the logfile.
func (lf *Logfile) Path() string {
	return lf.path
}

// EnsureDirectory creates the directories that the logfile is located in if they do not already exist.
func (lf *Logfile) EnsureDirectory() error {
	dir := filepath.Dir(lf.path)
	err := os.MkdirAll(dir, 0700)
	return err
}

// ReadEntries provides all entries stored in the logfile.
func (lf *Logfile) ReadEntries() ([]string, error) {
	f, err := os.Open(lf.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var entries []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		entries = append(entries, scanner.Text())
	}
	return entries, nil
}