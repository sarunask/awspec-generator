package resources

import (
	"fmt"
	"strings"
)

type Structure interface {
	String() string
}

type RuleType int

const (
	Ingress RuleType = iota
	Egress
)

type SG_rule struct {
	Type RuleType
	Port int64
	Protocol string
	CIDR_blocks []string
	Other_SG []string
}

func (t *SG_rule) String(dependencies *map[int]*Resource) (ret string) {
	sg_type := ":inbound"
	if t.Type == Egress {
		sg_type = ":outbound"
	}
	sg_protocol := ""
	if t.Protocol != "" {
		sg_protocol = fmt.Sprintf(".protocol('%v')", t.Protocol)
	}
	for_cidr := ""
	if len(t.CIDR_blocks) > 0 {
		for _, value := range t.CIDR_blocks {
			for_cidr += fmt.Sprintf(".for('%v')", value)
		}
	}
	if len(t.Other_SG) > 0 {
		for _, value := range t.Other_SG {
			//Search in dependencies for name of this group
			for _, dep := range *dependencies {
				if dep.Type == SG {
					id := get_attribute_by_name(dep.Attrs, "id")
					if id.String() == value {
						for_cidr += fmt.Sprintf(".for('%v')", dep.Name)
					}
				}
			}
		}
	}

	ret = fmt.Sprintf("  its(%v) { should be_opened(%v)%v%v }\n", sg_type,
		t.Port, sg_protocol, for_cidr)
	return
}

type ELB_HealthCheck struct {
	Healthy_Threshold int64
	Unhealthy_Threshold int64
	Interval int64
	Target string
	Timeout int64
}

func (t *ELB_HealthCheck) String() (ret string) {
	if t.Target != "" {
		ret += fmt.Sprintf("  its(:health_check_target) {should eq '%v'}\n", t.Target)
	}
	ret += fmt.Sprintf("  its(:health_check_interval) {should eq %v}\n", t.Interval)
	ret += fmt.Sprintf("  its(:health_check_timeout) {should eq %v}\n", t.Timeout)
	ret += fmt.Sprintf("  its(:health_check_unhealthy_threshold) {should eq %v}\n", t.Unhealthy_Threshold)
	ret += fmt.Sprintf("  its(:health_check_healthy_threshold) {should eq %v}\n", t.Healthy_Threshold)
	return
}

type ELB_Listener struct {
	Instance_port int64
	Instance_protocol string
	Lb_port int64
	Lb_protocol string
}

func (t *ELB_Listener) String() (ret string) {
	ret = fmt.Sprintf("  it { should have_listener(protocol: '%v', port: %v, " +
		"instance_protocol: '%v', instance_port: %v) }\n",
		strings.ToUpper(t.Lb_protocol),
		t.Lb_port,
		strings.ToUpper(t.Instance_protocol),
		t.Instance_port)
	return
}

type Tag struct {
	Name string
	Value string
}

func (t *Tag) String (ret string) {
	ret = fmt.Sprintf("  it { should have_tag('%v').value('%v') }",
		t.Name, t.Value)
	return
}

func get_tag(arr *map[string]*Tag, id string) (ret *Tag) {
	value, ok := (*arr)[id]
	if ok == false {
		ret = new(Tag)
		(*arr)[id] = ret
	} else {
		ret = value
	}
	return
}

type Alias struct {
	Name string
	Zone_id string
}

func (t *Alias) String() (ret string) {
	ret = fmt.Sprintf("alias('%v.', '%v')", t.Name, t.Zone_id)
	return
}

func get_alias(arr *map[string]*Alias, id string) (ret *Alias) {
	value, ok := (*arr)[id]
	if ok == false {
		ret = new(Alias)
		(*arr)[id] = ret
	} else {
		ret = value
	}
	return
}