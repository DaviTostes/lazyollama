package commands

import (
	"errors"
	"lazyollama/ollama"
	"os/exec"
	"regexp"
	"strings"
)

func CopyCode(message ollama.MessageChat) error {
	regex := regexp.MustCompile("(?s)```(?:\\w+)?\\s*(.*?)\\s*```")
	matches := regex.FindStringSubmatch(message.Content)

	if len(matches) < 2 {
		return errors.New("Nothing found")
	}

	code := strings.TrimSpace(matches[1])
	code = strings.ReplaceAll(code, "\x1b[0m", "")

	if _, err := exec.LookPath("xclip"); err == nil {
		cmd := exec.Command("xclip", "-selection", "clipboard")
		stdin, _ := cmd.StdinPipe()
		defer stdin.Close()
		_ = cmd.Start()
		_, _ = stdin.Write([]byte(code))
	}

	if _, err := exec.LookPath("wl-copy"); err == nil {
		cmd := exec.Command("wl-copy")
		stdin, _ := cmd.StdinPipe()
		defer stdin.Close()
		_ = cmd.Start()
		_, _ = stdin.Write([]byte(code))
	}

	return nil
}
