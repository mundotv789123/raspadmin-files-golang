package generator

import (
	"os/exec"
)

type IconEmbedGenerator struct {
	next IconGenerator
}

func (g *IconEmbedGenerator) Generate(filePath string, iconPath string) (bool, error) {
	cmd := exec.Command("ffmpeg", "-i", filePath, "-map", "0:v", "-map", "-0:V", "-c", "copy", iconPath)
	err := cmd.Run()

	if err != nil {
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
