package generator

type IconAudioGenerator struct {
	next IconGenerator
}

func (g *IconAudioGenerator) Generate(filePath string, iconPath string) (bool, error) {
	return g.next.Generate(filePath, iconPath)
}

func (g *IconAudioGenerator) SetNext(next IconGenerator) {
	g.next = next
}

func NewIconAudioGenerator() *IconAudioGenerator {
	g := &IconAudioGenerator{}
	g.SetNext(NewIconEmbedGenerator())
	return g
}
