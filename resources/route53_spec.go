package resources

import (
	"fmt"
	"regexp"
	"github.com/tidwall/gjson"
	"strings"
)

func (t Resource) aws_route53_zone_record_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe route53_hosted_zone('%v') do\n"+
		"  it { should exist }\n" +
		t.route53_attrs() +
		"end\n", get_attribute_by_name(t.Attrs, "zone_id").String())
}

func (t Resource) route53_attrs() (ret string)  {
	//Route53 records attributes for zone
	var ttl int64
	var fqdn string
	var record_type string
	records_pattern := regexp.MustCompile(`^records.[0-9]+$`)
	records := make([]string, 0)
	alias_pattern := regexp.MustCompile(`^alias\.([0-9]+)\.(.+)$`)
	alias_arr := make(map[string]*Alias)
	(*t.Attrs).ForEach(func(key, value gjson.Result) bool {
		key_string := key.String()
		value_string := value.String()
		switch strings.ToLower(key_string) {
		case "fqdn":
			fqdn = value_string
		case "ttl":
			ttl = value.Int()
		case "type":
			record_type = strings.ToLower(value_string)
		}
		if records_pattern.MatchString(key_string) {
			records = append(records, value_string)
		}
		if alias_pattern.MatchString(key_string) {
			pattern_matches := alias_pattern.FindStringSubmatch(key_string)
			alias := get_alias(&alias_arr, pattern_matches[1])
			switch strings.ToLower(pattern_matches[2]) {
			case "name":
				alias.Name = value_string
			case "zone_id":
				alias.Zone_id = value_string
			}
		}
		return true
	})
	records_str := ""
	if len(records) != 0 {
		if strings.EqualFold(record_type, "A") {
			records_str += fmt.Sprintf("a(EC2Helper::GetRoute53AFromZoneAndName('%v', '%v.'))",
				get_attribute_by_name(t.Attrs, "zone_id").String(), fqdn)
		} else {
			records_str += fmt.Sprintf("%v(", record_type)
			for _, value := range records {
				records_str += fmt.Sprintf("'%v' ", value)
			}
			records_str += ")"
		}
	}
	if len(alias_arr) != 0 {
		for _, val := range alias_arr {
			records_str += fmt.Sprintf("%v", val.String())
		}
	}

	ret = fmt.Sprintf("  it { should have_record_set('%v.').%v",
		fqdn, records_str)
	if ttl != 0 {
		ret += fmt.Sprintf(".ttl(%v) }\n", ttl)
	} else {
		ret += "}\n"
	}
	return
}
