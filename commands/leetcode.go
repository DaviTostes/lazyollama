package commands

import (
	"crypto/rand"
	"lazyollama/ollama"
	"os"
	"os/exec"
	"time"
)

func GenerateLeetHack(
	ollamaClient *ollama.Ollama,
) (*ollama.ResponseChat, *ollama.MessageChat, error) {
	var tempPath = os.Getenv("LAZY_DIR") + "/_temp/"

	if _, err := os.Stat(tempPath); err != nil {
		os.Mkdir(tempPath, 0777)
	}

	imageName := rand.Text() + ".png"
	outputPath := tempPath + imageName

	cmd := exec.Command("hyprshot", "-m", "window", "-o", tempPath, "-f", imageName)
	cmd.Run()

	time.Sleep(500 * time.Millisecond)

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return nil, nil, err
	}

	cmd = exec.Command(
		"tesseract",
		outputPath,
		tempPath+"/out",
		"--tessdata-dir",
		"/usr/local/share/tessdata",
	)
	if err := cmd.Run(); err != nil {
		return nil, nil, err
	}

	file, err := os.ReadFile(tempPath + "/out.txt")
	if err != nil {
		return nil, nil, err
	}

	prompt := `You are a JavaScript coding assistant. I will provide you with the description of a LeetCode problem. Your task is to:

  Analyze the problem.

  Solve it.

  Return the final JavaScript solution.

	Don't add any print ou comment in the code.

	You won't return any explanation, just the code.

Here is the problem description:	` + string(
		file,
	)

	gen, userMsg, err := ollamaClient.Generate(prompt, []ollama.MessageChat{})
	if err != nil {
		return nil, nil, err
	}

	os.Remove(outputPath)
	os.Remove(tempPath + "/out.txt")

	return gen, userMsg, err
}
