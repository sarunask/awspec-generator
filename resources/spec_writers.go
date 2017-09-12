package resources

import (
	"fmt"
	"github.com/tidwall/gjson"
	"strings"
	"regexp"
)


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

func (t Resource) aws_elb_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
	    "describe elb('%v') do\n" +
	    "  it { should exist }\n" +
		"  its(:load_balancer_name) { should eq '%v' }\n" +
	    t.elb_attrs() +
	    "end\n", t.Name, t.Name)
}

func (t Resource) aws_autoscaling_group_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe autoscaling_group(EC2Helper.GetASGIdFromName('%v')) do\n"+
		"  it { should exist }\n" +
		t.tags() +
		t.asg_attrs() +
		"end\n", t.Name)
}

func (t Resource) aws_rds_instance_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe rds(EC2Helper.GetRDSIdFromName('%v')) do\n"+
		"  it { should exist }\n" +
		"  it { should be_available }\n" +
		t.tags() +
		t.rds_attrs() +
		t.sg_dependencies() +
		"end\n", t.Name)
}

func (t Resource) aws_iam_policy_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe iam_policy('%v') do\n"+
		"  it { should exist }\n" +
		"  it { should be_attachable }\n" +
		"end\n", t.Name)
}

func (t Resource) aws_iam_role_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe iam_role('%v') do\n"+
		"  it { should exist }\n" +
		"end\n", t.Name)
}

func (t Resource) aws_ec2_instance_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe ec2(EC2Helper.GetEC2IdFromName('%v','%v')) do\n"+
		"  it { should exist }\n" +
		"  it { should be_running }\n" +
	    t.tags() +
	    t.ec2_attrs() +
		t.sg_dependencies() +
		"end\n", t.Name, t.FindTagValue("service"))
}

func (t Resource) aws_route53_zone_record_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe route53_hosted_zone('%v') do\n"+
		"  it { should exist }\n" +
		t.route53_attrs() +
		"end\n", t.get_route53_zone_id())
}

func (t Resource) get_route53_zone_id() (ret string) {
	(*t.Attrs).ForEach(func(key, value gjson.Result) bool {
		if strings.EqualFold(key.String(), "zone_id") {
			ret = value.String()
			return false
		}
		return true
	})
	return
}

func (t Resource) tags() (ret string) {
	for _, value := range t.Tags {
		ret += fmt.Sprintf("  it { should have_tag('%v').value('%v') }\n",
			value.Name, value.Value)
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

func (t Resource) sg_dependencies() (ret string) {
	for i := range t.Dependent {
		switch t.Dependent[i].Type {
		case SG:
			ret += fmt.Sprintf("  it { should have_security_group('%v') }\n",
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
	regexp_pattern := regexp.MustCompile(`^(ingress|egress)\.([0-9]+)\.(to_port|protocol)$`)
	ingress_rules := make(map[string]*SG_rule, 100)
	egress_rules := make(map[string]*SG_rule, 100)

	(*t.Attrs).ForEach(func(key, value gjson.Result) bool {
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

func append_array_when_regexp_match(reg *regexp.Regexp, arr *[]string, key string, val string) {
	//This function would append to provided array of string (arr) val if key match Regexp reg
	if reg.MatchString(key) {
		*arr = append( *arr, val)
	}
}

func (t Resource) elb_attrs() (eattrs string)  {
	availability_zones_regexp := regexp.MustCompile(`^availability_zones.[0-9]+$`)
	availability_zones := make([]string, 0)
	healthcheck_pattern := regexp.MustCompile(`^health_check\.([0-9]+)\.(.+)$`)
	healthchecks := make(map[string]*ELB_HealthCheck, 10)
	listener_pattern := regexp.MustCompile(`^listener\.([0-9]+)\.(.+)$`)
	listeners := make(map[string]*ELB_Listener, 10)
	(*t.Attrs).ForEach(func(key, value gjson.Result) bool {
		key_string := key.String()
		value_string := value.String()
		value_int64 := value.Int()
		append_array_when_regexp_match(availability_zones_regexp, &availability_zones, key_string, value_string)
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
	availability_zones_str := ""
	for _, value := range availability_zones {
		availability_zones_str += fmt.Sprintf("'%v',", value)
	}
	eattrs += fmt.Sprintf("  its(:availability_zones) { should == [%v] }\n", availability_zones_str)
	for _, val := range healthchecks {
		eattrs += val.String()
	}
	for _, val := range listeners {
		eattrs += val.String()
	}
	return
}

func (t Resource) asg_attrs() (ret string)  {
	//AutoScallingGroups attributes
	availability_zones_regexp := regexp.MustCompile(`^availability_zones.[0-9]+$`)
	availability_zones := make([]string, 0)
	elbs_regexp := regexp.MustCompile(`^load_balancers.[0-9]+$`)
	elbs := make([]string, 0)
	termination_policies_regexp := regexp.MustCompile(`^termination_policies.[0-9]+$`)
	termination_policies := make([]string, 0)

	(*t.Attrs).ForEach(func(key, value gjson.Result) bool {
		key_string := key.String()
		value_string := value.String()
		append_array_when_regexp_match(availability_zones_regexp, &availability_zones, key_string, value_string)
		append_array_when_regexp_match(elbs_regexp, &elbs, key_string, value_string)
		append_array_when_regexp_match(termination_policies_regexp, &termination_policies, key_string, value_string)
		switch strings.ToLower(key_string) {
		case "name_prefix":
			ret += fmt.Sprintf("  its(:auto_scaling_group_name) { should start_with('%v')}\n",
				value_string)
		case "launch_configuration":
			ret += fmt.Sprintf("  it { should have_launch_configuration('%v') }\n",
				value_string)
		case "max_size":
			ret += fmt.Sprintf("  its(:max_size) { should == %v }\n",
				value_string)
		case "min_size":
			ret += fmt.Sprintf("  its(:min_size) { should == %v }\n",
				value_string)
		case "desired_capacity":
			ret += fmt.Sprintf("  its(:desired_capacity) { should == %v }\n",
				value_string)
		case "default_cooldown":
			ret += fmt.Sprintf("  its(:default_cooldown) { should == %v }\n",
				value_string)
		case "health_check_type":
			ret += fmt.Sprintf("  its(:health_check_type) { should == '%v' }\n",
				value_string)
		case "health_check_grace_period":
			ret += fmt.Sprintf("  its(:health_check_grace_period) { should == %v }\n",
				value_string)
		case "protect_from_scale_in":
			ret += fmt.Sprintf("  its(:new_instances_protected_from_scale_in) { should == %v }\n",
				value_string)
		case "placement_group":
			plgrp_str := "nil"
			if value_string != "" {
				plgrp_str = fmt.Sprintf("'%v'", value_string)
			}
			ret += fmt.Sprintf("  its(:placement_group) { should == %v }\n",
				plgrp_str)
		}

		return true
	})
	availability_zones_str := ""
	for _, value := range availability_zones {
		availability_zones_str += fmt.Sprintf("'%v',", value)
	}
	ret += fmt.Sprintf("  its(:availability_zones) { should == [%v] }\n",
		availability_zones_str)
	elbs_str := ""
	for _, value := range elbs {
		elbs_str += fmt.Sprintf("'%v',", value)
	}
	ret += fmt.Sprintf("  its(:load_balancer_names) { should == [%v] }\n",
		elbs_str)
	termination_policies_str := ""
	for _, value := range termination_policies {
		termination_policies_str += fmt.Sprintf("'%v',", value)
	}
	ret += fmt.Sprintf("  its(:termination_policies) { should == [%v] }\n",
		termination_policies_str)

	return
}

func (t Resource) rds_attrs() (ret string)  {
	//RDS DB instance attributes
	(*t.Attrs).ForEach(func(key, value gjson.Result) bool {
		key_string := key.String()
		value_string := value.String()
		switch strings.ToLower(key_string) {
		case "option_group_name":
			ret += fmt.Sprintf("  it { should have_option_group('%v')}\n",
				value_string)
		case "instance_class":
			ret += fmt.Sprintf("  its(:db_instance_class) { should eq '%v' }\n",
				value_string)
		case "engine":
			ret += fmt.Sprintf("  its(:engine) { should eq '%v' }\n",
				value_string)
		case "engine_version":
			ret += fmt.Sprintf("  its(:engine_version) { should eq '%v' }\n",
				value_string)
		case "db_instance_class":
			ret += fmt.Sprintf("  its(:db_instance_class) { should eq '%v' }\n",
				value_string)
		case "username":
			ret += fmt.Sprintf("  its(:master_username) { should eq '%v' }\n",
				value_string)
		case "name":
			ret += fmt.Sprintf("  its(:db_name) { should eq '%v' }\n",
				value_string)
		case "allocated_storage":
			ret += fmt.Sprintf("  its(:allocated_storage) { should eq %v }\n",
				value_string)
		case "availability_zone":
			ret += fmt.Sprintf("  its(:availability_zone) { should eq '%v' }\n",
				value_string)
		case "backup_retention_period":
			ret += fmt.Sprintf("  its(:backup_retention_period) { should eq %v }\n",
				value_string)
		case "maintenance_window":
			ret += fmt.Sprintf("  its(:preferred_maintenance_window) { should eq '%v' }\n",
				value_string)
		case "backup_window":
			ret += fmt.Sprintf("  its(:preferred_backup_window) { should eq '%v' }\n",
				value_string)
		case "multi_az":
			ret += fmt.Sprintf("  its(:multi_az) { should eq %v }\n",
				value_string)
		case "publicly_accessible":
			ret += fmt.Sprintf("  its(:publicly_accessible) { should eq %v }\n",
				value_string)
		case "auto_minor_version_upgrade":
			ret += fmt.Sprintf("  its(:auto_minor_version_upgrade) { should eq %v }\n",
				value_string)
		case "storage_type":
			ret += fmt.Sprintf("  its(:storage_type) { should eq '%v' }\n",
				value_string)
		case "storage_encrypted":
			ret += fmt.Sprintf("  its(:storage_encrypted) { should eq %v }\n",
				value_string)
		case "kms_key_id":
			ret += fmt.Sprintf("  its(:kms_key_id) { should eq '%v' }\n",
				value_string)
		case "copy_tags_to_snapshot":
			ret += fmt.Sprintf("  its(:copy_tags_to_snapshot) { should eq %v }\n",
				value_string)
		case "monitoring_interval":
			ret += fmt.Sprintf("  its(:monitoring_interval) { should eq %v }\n",
				value_string)
		}
		return true
	})
	return
}

func (t Resource) ec2_attrs() (ret string)  {
	//EC2 instance attributes
	(*t.Attrs).ForEach(func(key, value gjson.Result) bool {
		key_string := key.String()
		value_string := value.String()
		switch strings.ToLower(key_string) {
		case "ami":
			ret += fmt.Sprintf("  its(:image_id) { should eq '%v' }\n",
				value_string)
		case "instance_type":
			ret += fmt.Sprintf("  its(:instance_type) { should eq '%v' }\n",
				value_string)
		case "key_name":
			format :=  "  its(:key_name) { should eq '%v' }\n"
			if value_string == "" {
				format =  "  its(:key_name) { should eq nil%v }\n"
			}
			ret += fmt.Sprintf(format, value_string)
		case "ebs_optimized":
			ret += fmt.Sprintf("  its(:ebs_optimized) { should eq %v }\n",
				value_string)
		}
		return true
	})
	return
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
		records_str += fmt.Sprintf("%v(", record_type)
		for _, value := range records {
			records_str += fmt.Sprintf("'%v' ", value)
		}
		records_str += ")"
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
