package resources

import (
	"github.com/tidwall/gjson"
	"strings"
	"regexp"
	"fmt"
)

func get_attribute_by_name(attrs *gjson.Result, name string) (ret gjson.Result) {
	attrs.ForEach(func(key, value gjson.Result) bool {
		if strings.EqualFold(key.String(), name) {
			ret = value
			return false
		}
		return true
	})
	return
}

func get_list_items_by_pattern(attrs *gjson.Result, pattern string) *[]string {
	ret := make([]string, 0)
	regexp_pattern := regexp.MustCompile(pattern)
	attrs.ForEach(func(key, value gjson.Result) bool {
		key_string := key.String()
		value_string := value.String()
		append_array_when_regexp_match(regexp_pattern, &ret, key_string, value_string)
		return true
	})
	return &ret
}

func append_array_when_regexp_match(reg *regexp.Regexp, arr *[]string, key string, val string) {
	//This function would append to provided array of string (arr) val if key match Regexp reg
	if reg.MatchString(key) {
		*arr = append( *arr, val)
	}
}

func create_ruby_string_array(a *[]string, pattern string) (ret string) {
	if len(*a) == 0 {
		return
	}
	var ruby_array string
	for _, str := range *a {
		ruby_array += fmt.Sprintf("'%v',", str)
	}
	ret = fmt.Sprintf(pattern, ruby_array)
	return
}

func create_ruby_string(str string) (ret string) {
	ret = "nil"
	if strings.EqualFold(str, "true") || strings.EqualFold(str, "false") {
		ret = fmt.Sprintf("%v", str)
	} else if len(str) != 0 {
		ret = fmt.Sprintf("'%v'", str)
	}
	return
}

func create_ruby_hash(m *map[string]string) (ret string) {
	if len(*m) == 0 {
		ret = "nil"
		return
	}
	ret = "{"
	for key, value := range *m {
		ret += fmt.Sprintf("'%v'=>'%v',", key, value)
	}
	ret += "}"
	return
}

func create_ruby_array(a *[]string) (ret string) {
	if len(*a) == 0 {
		ret = "nil"
		return
	}
	ret = "["
	for _, value := range *a {
		ret += fmt.Sprintf("'%v',", value)
	}
	ret += "]"
	return
}