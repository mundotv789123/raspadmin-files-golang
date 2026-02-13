package generator

import "fmt"

type IconVideoGenerator struct{}

func (g *IconVideoGenerator) Generate(filePath string, iconPath string) error {
	fmt.Println("Gerado ícone de vídeo para:", filePath, "salvando em:", iconPath)
	return nil
}
