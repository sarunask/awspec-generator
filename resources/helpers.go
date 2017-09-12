package resources

import (
	"github.com/tidwall/gjson"
	"strings"
)

func GetAttributeByName(attrs *gjson.Result, name string) (ret gjson.Result) {
	attrs.ForEach(func(key, value gjson.Result) bool {
		if strings.EqualFold(key.String(), name) {
			ret = value
			return false
		}
		return true
	})
	return
}
