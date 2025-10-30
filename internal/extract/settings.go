package extract

import (
	"github.com/klauspost/compress/zip"
	"strings"
)

func isValid(r *zip.ReadCloser) bool {
	if len(r.File) == 0 {
		return false
	}
	return true
}

func isOption(r *zip.ReadCloser) bool {
	if isValid(r) {
		optionOrSample := getOptionSample(r)
		if optionOrSample == "" {
			return false
		}
		return strings.EqualFold(optionOrSample, "OPTIONS")
	}

	return false
}

func getOptionSample(r *zip.ReadCloser) string {
	t := strings.Split(fixBaseStart(r.File[0].Name), "/")
	if len(t) <= 1 {
		return ""
	}
	if strings.EqualFold("options", t[0]) || strings.EqualFold("samples", t[0]) {
		return t[0]
	}
	return ""
}

func getSampleName(r *zip.ReadCloser) string {
	t := strings.Split(fixBaseStart(r.File[0].Name), "/")
	if len(t) <= 2 {
		return ""
	}
	return t[1]
}

func getOptionSampleSampleName(r *zip.ReadCloser) string {
	return getOptionSample(r) + "/" + getSampleName(r) + "/"
}

func getSampleNameLast(sampleName string) string {
	split := strings.Split(sampleName, "_")
	return split[len(split)-1]
}

func fixBaseStart(s string) string {
	return strings.TrimPrefix(s, "/")
}
