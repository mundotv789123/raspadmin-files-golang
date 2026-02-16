package generator

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type IconEmbedGenerator struct {
	next IconGenerator
}

func (g *IconEmbedGenerator) Generate(filePath string, iconPath string) (bool, error) {
	cmd := exec.Command("ffmpeg", "-i", filePath, "-map", "0:v", "-map", "-0:V", "-c", "copy", iconPath)
	err := cmd.Run()

	if err != nil {
		// Tratar quando for erro de stream, assumir que o vídeo não contem imagem
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 234 {
			return false, nil
		}
		cmdOut, _ := cmd.Output()
		return false, fmt.Errorf("%s\n%s", err, string(cmdOut))
	}

	_, err = os.Stat(iconPath)
	if err != nil {
		if errors.Is(os.ErrNotExist, err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (g *IconEmbedGenerator) SetNext(next IconGenerator) {
	g.next = next
}

func NewIconEmbedGenerator() *IconEmbedGenerator {
	return &IconEmbedGenerator{}
}
