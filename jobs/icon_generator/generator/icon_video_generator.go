package generator

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type IconVideoGenerator struct {
	next IconGenerator
}

func (g *IconVideoGenerator) Generate(filePath string, iconPath string) (bool, error) {
	ok, err := g.next.Generate(filePath, iconPath)
	if err != nil || ok {
		return ok, err
	}

	cmd := exec.Command("ffmpegthumbnailer", "-i", filePath, "-o", iconPath, "-s", "512")
	err = cmd.Run()

	if err != nil {
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

func (g *IconVideoGenerator) SetNext(next IconGenerator) {
	g.next = next
}

func NewIconVideoGenerator() *IconVideoGenerator {
	g := &IconVideoGenerator{}
	g.SetNext(NewIconEmbedGenerator())
	return g
}
