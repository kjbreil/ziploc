package option

import (
	"log"
	"regexp"
)

func (o *Option) excluded(current string) bool {
	for _, item := range o.Exclude {
		reg, err := regexp.Compile("(?i)" + item)
		if err != nil {
			log.Panicf("regex could not compile %s: %v", item, err)
		}

		if reg.MatchString(current) {
			return true
		}

	}

	// for _, exclude := range o.Exclude {
	// 	if strings.EqualFold(exclude, current) {
	// 		return true
	// 	}
	// }
	return false
}
