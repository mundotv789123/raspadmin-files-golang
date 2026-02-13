package generator

import "fmt"

type IconAudioGenerator struct{}

func (g *IconAudioGenerator) Generate(filePath string, iconPath string) error {
	fmt.Println("Gerado ícone de áudio para:", filePath, "salvando em:", iconPath)
	return nil
}
