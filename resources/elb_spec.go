package resources

import (
	"fmt"
	"regexp"
	"github.com/tidwall/gjson"
	"strings"
)

func (t Resource) aws_elb_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe elb('%v') do\n" +
		"  it { should exist }\n" +
		"  its(:load_balancer_name) { should eq '%v' }\n" +
		t.elb_attrs() +
		"end\n", t.Name, t.Name)
}

func (t Resource) elb_attrs() (eattrs string)  {
	availability_zones := get_list_items_by_pattern(t.Attrs, `^availability_zones\.[0-9]+$`)
	healthcheck_pattern := regexp.MustCompile(`^health_check\.([0-9]+)\.(.+)$`)
	healthchecks := make(map[string]*ELB_HealthCheck, 10)
	listener_pattern := regexp.MustCompile(`^listener\.([0-9]+)\.(.+)$`)
	listeners := make(map[string]*ELB_Listener, 10)
	(*t.Attrs).ForEach(func(key, value gjson.Result) bool {
		key_string := key.String()
		value_string := value.String()
		value_int64 := value.Int()
		if healthcheck_pattern.MatchString(key_string) {
			pattern_matches := healthcheck_pattern.FindStringSubmatch(key_string)
			healthcheck := get_elb_healthcheck(&healthchecks, pattern_matches[1])
			switch strings.ToLower(pattern_matches[2]) {
			case "healthy_threshold":
				healthcheck.Healthy_Threshold = value_int64
			case "interval":
				healthcheck.Interval = value_int64
			case "target":
				healthcheck.Target = value_string
			case "timeout":
				healthcheck.Timeout = value_int64
			case "unhealthy_threshold":
				healthcheck.Unhealthy_Threshold = value_int64
			}
		}
		if listener_pattern.MatchString(key_string) {
			pattern_matches := listener_pattern.FindStringSubmatch(key_string)
			listener := get_elb_listener(&listeners, pattern_matches[1])
			switch strings.ToLower(pattern_matches[2]) {
			case "instance_port":
				listener.Instance_port = value_int64
			case "instance_protocol":
				listener.Instance_protocol = value_string
			case "lb_port":
				listener.Lb_port = value_int64
			case "lb_protocol":
				listener.Lb_protocol = value_string
			}
		}
		if strings.EqualFold(strings.ToLower(key_string), "zone_id") {
			eattrs += fmt.Sprintf("  its(:canonical_hosted_zone_name_id) { should eq '%v' }\n", value_string)
		}
		return true
	})
	eattrs += create_ruby_string_array(availability_zones, "  its(:availability_zones) { should =~ [%v] }\n")
	for _, val := range healthchecks {
		eattrs += val.String()
	}
	for _, val := range listeners {
		eattrs += val.String()
	}
	return
}

func get_elb_healthcheck(arr *map[string]*ELB_HealthCheck, id string) (ret *ELB_HealthCheck) {
	value, ok := (*arr)[id]
	if ok == false {
		ret = new(ELB_HealthCheck)
		(*arr)[id] = ret
	} else {
		ret = value
	}
	return
}

func get_elb_listener(arr *map[string]*ELB_Listener, id string) (ret *ELB_Listener) {
	value, ok := (*arr)[id]
	if ok == false {
		ret = new(ELB_Listener)
		(*arr)[id] = ret
	} else {
		ret = value
	}
	return
}
