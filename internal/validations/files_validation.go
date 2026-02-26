package validations

import "regexp"

var hiddenFilesRegex = regexp.MustCompile(`^[\._\$].*$`)
var tempFilesRegex = regexp.MustCompile(`(\.swp|\.part|\.temp.*)$`)

func IsHiddenFile(fileName string) bool {
	return hiddenFilesRegex.MatchString(fileName)
}

func IsTempFile(fileName string) bool {
	return tempFilesRegex.MatchString(fileName)
}

func IsTempOrHiddenFile(fileName string) bool {
	return IsHiddenFile(fileName) || IsTempFile(fileName)
}
