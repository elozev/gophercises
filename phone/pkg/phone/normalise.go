package normalise

import (
	"regexp"
	"strings"
)

var rStartsWithPlus = regexp.MustCompile("^\\+")

func removePlusAtStart(phone string) string {
	r := regexp.MustCompile("^\\+")
	if r.MatchString(phone) {
		parts := r.Split(phone, -1)
		return "00" + parts[1]
	}

	return phone
}

func Clean(phone string) string {
	phoneTrimmed := removePlusAtStart(strings.TrimSpace(phone))
	rNonNumeric := regexp.MustCompile("\\D")
	pSplit := rNonNumeric.Split(phoneTrimmed, -1)
	return strings.Join(pSplit, "")
}
