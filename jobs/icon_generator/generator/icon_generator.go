package generator

import (
	"regexp"
)

var (
	videoTypeRegex = regexp.MustCompile(`^video/(mp4|mkv|webm)$`)
	audioTypeRegex = regexp.MustCompile(`^audio/(mpeg)$`)
)

type IconGenerator interface {
	Generate(filePath string, iconPath string) (bool, error)
	SetNext(IconGenerator)
}

func GenerateIcon(filePath string, iconPath string, generator IconGenerator) (bool, error) {
	return generator.Generate(filePath, iconPath)
}

func GetGenerator(contentType string) (IconGenerator, bool) {
	if videoTypeRegex.MatchString(contentType) {
		return NewIconVideoGenerator(), true
	}
	if audioTypeRegex.MatchString(contentType) {
		return NewIconAudioGenerator(), true
	}
	return nil, false
}
