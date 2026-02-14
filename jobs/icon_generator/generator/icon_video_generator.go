package generator

import "os/exec"

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
