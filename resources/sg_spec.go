package resources

import (
	"fmt"
	"strings"
	"regexp"
	"github.com/tidwall/gjson"
)

func (t Resource) aws_sg_spec() string {
	return fmt.Sprintf("require 'awspec'\n\n" +
		"describe security_group('%v') do\n"+
		"  it { should exist }\n" +
		"  its('group_name') { should start_with('%v')}\n" +
		t.tags() +
		t.sg_rules() +
		"end\n", t.Name, t.Name)
}

func (t Resource) sg_rules() (ret string)  {
	regexp_pattern := regexp.MustCompile(`^(ingress|egress)\.([0-9]+)\.(to_port|protocol)$`)
	cidr_pattern := regexp.MustCompile(`^(ingress|egress)\.([0-9]+)\.cidr_blocks\.[0-9]+$`)
	sg_pattern := regexp.MustCompile(`^(ingress|egress)\.([0-9]+)\.security_groups\.[0-9]+$`)
	ingress_rules := make(map[string]*SG_rule, 100)
	egress_rules := make(map[string]*SG_rule, 100)

	(*t.Attrs).ForEach(func(key, value gjson.Result) bool {
		key_string := key.String()
		value_string := value.String()
		if strings.EqualFold(key_string, "egress.#") {
			ret += fmt.Sprintf("  its(:outbound_rule_count) { should eq %v }\n", value_string)
			return true
		}
		if strings.EqualFold(key_string, "ingress.#") {
			ret += fmt.Sprintf("  its(:inbound_rule_count) { should eq %v }\n", value_string)
			return true
		}
		if regexp_pattern.MatchString(key_string) {
			pattern_matches := regexp_pattern.FindStringSubmatch(key_string)
			sg_rule := get_sg_rule(&ingress_rules, pattern_matches[2])
			sg_rule.Type = Ingress
			if strings.EqualFold(pattern_matches[1], "egress") {
				sg_rule.Type = Egress
			}
			if strings.EqualFold(pattern_matches[3], "to_port") {
				sg_rule.Port = value.Int()
			} else {
				sg_rule.Protocol = strings.ToLower(value_string)
			}

		}
		if cidr_pattern.MatchString(key_string) {
			pattern_matches := cidr_pattern.FindStringSubmatch(key_string)
			sg_rule := get_sg_rule(&ingress_rules, pattern_matches[2])
			sg_rule.CIDR_blocks = append(sg_rule.CIDR_blocks, value_string)
		}
		if sg_pattern.MatchString(key_string) {
			pattern_matches := sg_pattern.FindStringSubmatch(key_string)
			sg_rule := get_sg_rule(&ingress_rules, pattern_matches[2])
			sg_rule.Other_SG = append(sg_rule.Other_SG, value_string)
		}
		return true
	})
	for _, value := range ingress_rules {
		ret += value.String(&t.Dependent)
	}
	for _, value := range egress_rules {
		ret += value.String(&t.Dependent)
	}
	return
}

func get_sg_rule(arr *map[string]*SG_rule, id string) (ret *SG_rule) {
	value, ok := (*arr)[id]
	if ok == false {
		ret = new(SG_rule)
		(*arr)[id] = ret
	} else {
		ret = value
	}
	return
}


