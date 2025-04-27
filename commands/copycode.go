package commands

import (
	"errors"
	"lazyollama/ollama"
	"os/exec"
	"regexp"
	"strings"
)

func CopyCode(message ollama.MessageChat, number int) error {
	regex := regexp.MustCompile("```(\\w+)?\\s*([\\s\\S]*?)```")
	matches := regex.FindAllStringSubmatch(message.Content, -1)

	if len(matches) == 0 {
		return errors.New("Nothing to copy")
	}

	match := matches[number]

	code := strings.TrimSpace(match[2])
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
