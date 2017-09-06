package resources

import (
	"fmt"
	"github.com/tidwall/gjson"
	"strings"
	"regexp"
)

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
}

func (t *SG_rule) String() (ret string) {
	sg_type := ":inbound"
	if t.Type == Egress {
		sg_type = ":outbound"
	}
	sg_protocol := ""
	if t.Protocol != "" {
		sg_protocol = fmt.Sprintf(".protocol('%v')", t.Protocol)
	}

	ret = fmt.Sprintf("  its(%v) { should be_opened(%v)%v }\n", sg_type,
		t.Port, sg_protocol)
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

func (t Resource) aws_vpc_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe vpc('%v') do\n"+
		"  it { should exist }\n" +
		"  it { should be_available }\n" +
		t.tags() +
		"  its(:vpc_id) { should eq EC2Helper.GetVPCIdFromName('%v') }\n" +
		"  its('Assigned IGW count'){ expect(EC2Helper.GetIGWsCountForVPCwithName('%v')).to eq 0 }\n" +
		"end\n", t.Name, t.Name, t.Name)
}

func (t Resource) aws_subnet_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe subnet('%v') do\n"+
		"  it { should exist }\n" +
		"  it { should be_available }\n" +
		t.tags() +
		t.subnet_vpc_id_should_be() +
		"end\n", t.Name)
}

func (t Resource) aws_sg_spec() string {
	return fmt.Sprintf("require 'awspec'\n\n" +
		"describe security_group('%v') do\n"+
		"  it { should exist }\n" +
	    "  its('group_name') { should start_with('%v')}\n" +
		t.tags() +
	    t.sg_rules() +
		"end\n", t.Name, t.Name)
}

func (t Resource) aws_vpn_gw_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe vpn_gateway('%v') do\n"+
		"  it { should exist }\n" +
		"  it { should be_available }\n" +
		"  its(:type) {should eq 'ipsec.1'}\n" +
		t.tags() +
		t.vpc_attachments_should_be() +
		"  its(:vpn_gateway_id) { should eq EC2Helper.GetVPNGWIdFromName('%v') }\n" +
		"end\n", t.Name, t.Name)
}

func (t Resource) aws_customer_gw_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe customer_gateway('%v') do\n"+
		"  it { should exist }\n" +
		"  it { should be_available }\n" +
		"  its(:type) {should eq 'ipsec.1'}\n" +
		"  its(:bgp_asn) {should eq '65000'}\n" +
		t.tags() +
		"end\n", t.Name)
}

func (t Resource) aws_vpn_connection_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe vpn_connection('%v') do\n"+
		"  it { should exist }\n" +
		"  it { should be_available }\n" +
		"  its(:type) {should eq 'ipsec.1'}\n" +
		t.tags() +
		t.vpn_gw_id_should_be() +
		"end\n", t.Name)
}

//describe elb('rss-non-prod-kapacitor-lb') do
//it { should exist }
//its(:load_balancer_name) { should eq 'rss-non-prod-kapacitor-lb' }
//its(:health_check_target) { should eq 'TCP:9092' }
//its(:health_check_interval) { should eq 30 }
//its(:health_check_timeout) { should eq 3 }
//its(:health_check_unhealthy_threshold) { should eq 2 }
//its(:health_check_healthy_threshold) { should eq 2 }
//it { should have_listener(protocol: 'HTTP', port: 9092, instance_protocol: 'HTTP', instance_port: 9092) }
//end

func (t Resource) aws_elb_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
	    "describe elb('%v') do\n" +
	    "  it { should exist }\n" +
		"  its(:load_balancer_name) { should eq '%v' }\n" +
	    t.elb_attrs() +
	    "end\n", t.Name, t.Name)
}

func (t Resource) tags() (tags string) {
	for key, value := range t.Tags {
		tags += fmt.Sprintf("  it { should have_tag('%v').value('%v') }\n", key, value)
	}
	return
}

func (t Resource) subnet_vpc_id_should_be() (ret string) {
	for i := range t.Dependent {
		switch t.Dependent[i].Type {
		case VPC:
			ret += fmt.Sprintf("  its(:vpc_id) { should eq EC2Helper.GetVPCIdFromName('%v') }\n",
				t.Dependent[i].Name)
		}
	}
	return
}

func (t Resource) vpn_gw_id_should_be() (ret string) {
	for i := range t.Dependent {
		switch t.Dependent[i].Type {
		case VPN_GW:
			ret += fmt.Sprintf("  its(:vpn_gateway_id) { should eq EC2Helper.GetVPNGWIdFromName('%v') }\n",
				t.Dependent[i].Name)
		}
	}
	return
}

func (t Resource) vpc_attachments_should_be() (ret string) {
	for i := range t.Dependent {
		switch t.Dependent[i].Type {
		case VPC:
			ret += fmt.Sprintf("  its(:vpc_attachments) { should eq " +
				"[Aws::EC2::Types::VpcAttachment.new(:state => 'attached', " +
				":vpc_id => EC2Helper.GetVPCIdFromName('%v'))] }\n",
				t.Dependent[i].Name)
		}
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

func (t Resource) sg_rules() (ret string)  {
	attrs := gjson.Get(t.Raw, "primary.attributes")
	regexp_pattern := regexp.MustCompile(`^(ingress|egress)\.([0-9]+)\.(to_port|protocol)$`)
	ingress_rules := make(map[string]*SG_rule, 100)
	egress_rules := make(map[string]*SG_rule, 100)

	attrs.ForEach(func(key, value gjson.Result) bool {
		key_string := key.String()
		if strings.EqualFold(key_string, "egress.#") {
			ret += fmt.Sprintf("  its(:outbound_rule_count) { should eq %v }\n", value.String())
			return true
		}
		if strings.EqualFold(key_string, "ingress.#") {
			ret += fmt.Sprintf("  its(:inbound_rule_count) { should eq %v }\n", value.String())
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
				sg_rule.Protocol = strings.ToLower(value.String())
			}

		}
		return true
	})
	for _, value := range ingress_rules {
		ret += value.String()
	}
	for _, value := range egress_rules {
		ret += value.String()
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

func (t Resource) elb_attrs() (eattrs string)  {
	attrs := gjson.Get(t.Raw, "primary.attributes")
	availability_zones_regexp := regexp.MustCompile(`^availability_zones.[0-9]+$`)
	availability_zones := make([]string, 0)
	healthcheck_pattern := regexp.MustCompile(`^health_check\.([0-9]+)\.(.+)$`)
	healthchecks := make(map[string]*ELB_HealthCheck, 10)
	listener_pattern := regexp.MustCompile(`^listener\.([0-9]+)\.(.+)$`)
	listeners := make(map[string]*ELB_Listener, 10)
	attrs.ForEach(func(key, value gjson.Result) bool {
		key_string := key.String()
		if availability_zones_regexp.MatchString(key_string) {
			availability_zones = append(availability_zones, value.String())
		}
		if healthcheck_pattern.MatchString(key_string) {
			pattern_matches := healthcheck_pattern.FindStringSubmatch(key_string)
			healthcheck := get_elb_healthcheck(&healthchecks, pattern_matches[1])
			switch strings.ToLower(pattern_matches[2]) {
			case "healthy_threshold":
				healthcheck.Healthy_Threshold = value.Int()
			case "interval":
				healthcheck.Interval = value.Int()
			case "target":
				healthcheck.Target = value.String()
			case "timeout":
				healthcheck.Timeout = value.Int()
			case "unhealthy_threshold":
				healthcheck.Unhealthy_Threshold = value.Int()
			}
		}
		if listener_pattern.MatchString(key_string) {
			pattern_matches := listener_pattern.FindStringSubmatch(key_string)
			listener := get_elb_listener(&listeners, pattern_matches[1])
			switch strings.ToLower(pattern_matches[2]) {
			case "instance_port":
				listener.Instance_port = value.Int()
			case "instance_protocol":
				listener.Instance_protocol = value.String()
			case "lb_port":
				listener.Lb_port = value.Int()
			case "lb_protocol":
				listener.Lb_protocol = value.String()
			}
		}
		if strings.EqualFold(strings.ToLower(key_string), "zone_id") {
			eattrs += fmt.Sprintf("  its(:canonical_hosted_zone_name_id) { should eq '%v' }\n", value.String())
		}
		return true
	})
	availability_zones_str := ""
	for _, value := range availability_zones {
		availability_zones_str += fmt.Sprintf("'%v',", value)
	}
	eattrs += fmt.Sprintf("  its(:availability_zones) { should == [%v] }\n", availability_zones_str)
	for _, value := range healthchecks {
		eattrs += value.String()
	}
	for _, value := range listeners {
		eattrs += value.String()
	}
	return
}
