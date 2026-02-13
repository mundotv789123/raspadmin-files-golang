package generator

import "fmt"

type IconEmbedGenerator struct{}

func (g *IconEmbedGenerator) Generate(filePath string, iconPath string) error {
	fmt.Println("Gerado ícone de embed para:", filePath, "salvando em:", iconPath)
	return nil
}
