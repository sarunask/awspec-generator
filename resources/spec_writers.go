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

func (t Resource) aws_vpc_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe vpc(\"%v\") do\n"+
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
		"describe subnet(\"%v\") do\n"+
		"  it { should exist }\n" +
		"  it { should be_available }\n" +
		t.tags() +
		t.subnet_vpc_id_should_be() +
		"end\n", t.Name)
}

func (t Resource) aws_sg_spec() string {
	return fmt.Sprintf("require 'awspec'\n\n" +
		"describe security_group(\"%v\") do\n"+
		"  it { should exist }\n" +
	    "  its('group_name') { should eq '%v'}\n" +
		t.tags() +
	    t.sg_rules() +
		"end\n", t.Name, t.Name)
}

func (t Resource) aws_vpn_gw_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe vpn_gateway(\"%v\") do\n"+
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
		"describe customer_gateway(\"%v\") do\n"+
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
		"describe vpn_connection(\"%v\") do\n"+
		"  it { should exist }\n" +
		"  it { should be_available }\n" +
		"  its(:type) {should eq 'ipsec.1'}\n" +
		t.tags() +
		t.vpn_gw_id_should_be() +
		"end\n", t.Name)
}

func (t Resource) tags() (ret string) {
	for key, value := range t.Tags {
		ret += fmt.Sprintf("  it { should have_tag('%v').value('%v') }\n", key, value)
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