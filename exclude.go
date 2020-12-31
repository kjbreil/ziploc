package ziploc

import (
	"strings"
)

func (o *Option) excluded(current string) bool {
	for _, exclude := range o.Exclude {
		if strings.EqualFold(exclude, current) {
			return true
		}
	}
	return false
}
