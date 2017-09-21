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
		t.asg_dependencies() +
		"end\n", t.Name)
}

func (t Resource) aws_launch_configuration_spec() string {
	name_prefix := get_attribute_by_name(t.Attrs, "name_prefix")

	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe launch_configuration(EC2Helper.GetLaunchConfigIdFromName('%v')) do\n"+
		"  it { should exist }\n" +
		"  its(:launch_configuration_name) { should start_with('%v') }\n" +
		t.lc_attrs() +
		"end\n", name_prefix.String(), name_prefix.String())
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
		"end\n", get_attribute_by_name(t.Attrs, "zone_id").String())
}

func (t Resource) aws_cloudwatch_metric_alarm_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe cloudwatch_alarm(%v) do\n"+
		"  it { should exist }\n" +
		"  its(:metric_name) { should eq %v }\n" +
		"  its(:alarm_description) { should eq %v }\n" +
		"  its(:namespace) { should eq %v }\n" +
		"  its(:actions_enabled) { should eq %v }\n" +
		"  its(:comparison_operator) { should eq %v }\n" +
		"  its(:threshold) { should eq %v }\n" +
		"  its(:evaluation_periods) { should eq %v }\n" +
		"  its(:unit) { should eq %v }\n" +
		"  its(:period) { should eq %v }\n" +
		"  its(:statistic) { should eq %v }\n" +
		"  its(:extended_statistic) { should eq %v }\n" +
		"  its(:evaluate_low_sample_count_percentile) { should eq %v }\n" +
		"  its(:treat_missing_data) { should eq %v }\n" +
		t.alarm_attrs() +
		"end\n",
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"alarm_name").String()),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"metric_name").String()),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"alarm_description").String()),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"namespace").String()),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"actions_enabled").String()),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"comparison_operator").String()),
		get_attribute_by_name(t.Attrs,"threshold").Uint(),
		get_attribute_by_name(t.Attrs,"evaluation_periods").Uint(),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"unit").String()),
		get_attribute_by_name(t.Attrs,"period").Uint(),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"statistic").String()),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"extended_statistic").String()),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"evaluate_low_sample_count_percentiles").String()),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"treat_missing_data").String()),
			)
}

func (t Resource) aws_lambda_function_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe lambda(%v) do\n"+
		"  it { should exist }\n" +
		"  its(:runtime) { should eq %v }\n" +
		"  its(:handler) { should eq %v }\n" +
		"  its(:description) { should eq %v }\n" +
		"  its(:timeout) { should eq %v }\n" +
		"  its(:memory_size) { should eq %v }\n" +
		t.lambda_attrs() +
		"end\n",
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"function_name").String()),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"runtime").String()),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"handler").String()),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"description").String()),
		get_attribute_by_name(t.Attrs,"timeout").Uint(),
		get_attribute_by_name(t.Attrs,"memory_size").Uint(),

	)
}

func (t Resource) aws_s3_bucket_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe s3_bucket(%v) do\n"+
		"  it { should exist }\n" +
		"  its(:name) { should start_with(%v) }\n" +
		t.tags() +
		"end\n",
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"bucket").String()),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"bucket_prefix").String()),
	)
}

func (t Resource) get_attribute(attr_name string) (ret string) {
	(*t.Attrs).ForEach(func(key, value gjson.Result) bool {
		if strings.EqualFold(key.String(), attr_name) {
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

func (t Resource) asg_dependencies() (ret string) {
	for i := range t.Dependent {
		switch t.Dependent[i].Type {
		case LAUNCH_CONFIG:
			name_prefix := get_attribute_by_name(t.Dependent[i].Attrs, "name_prefix")
			ret += fmt.Sprintf(
				"  it { should have_launch_configuration(EC2Helper.GetLaunchConfigIdFromName('%v')) }\n",
				name_prefix)
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
	eattrs += create_ruby_string_array(availability_zones, "  its(:availability_zones) { should == [%v] }\n")
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
	availability_zones := get_list_items_by_pattern(t.Attrs, `^availability_zones\.[0-9]+$`)
	elbs := get_list_items_by_pattern(t.Attrs, `^load_balancers\.[0-9]+$`)
	termination_policies := get_list_items_by_pattern(t.Attrs, `^termination_policies\.[0-9]+$`)

	(*t.Attrs).ForEach(func(key, value gjson.Result) bool {
		key_string := key.String()
		value_string := value.String()
		switch strings.ToLower(key_string) {
		case "name_prefix":
			ret += fmt.Sprintf("  its(:auto_scaling_group_name) { should start_with('%v')}\n",
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
	ret += create_ruby_string_array(availability_zones, "  its(:availability_zones) { should == [%v] }\n")
	ret += create_ruby_string_array(elbs, "  its(:load_balancer_names) { should == [%v] }\n")
	ret += create_ruby_string_array(termination_policies, "  its(:termination_policies) { should == [%v] }\n")
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

func (t Resource) lc_attrs() (ret string)  {
	//LaunchConfiguration attributes
	ami_id := get_attribute_by_name(t.Attrs, "image_id")
	if ami_id.String() != "" {
		ret += fmt.Sprintf("  its(:image_id) { should eq '%v' }\n", ami_id.String())
	}
	return
}

func (t Resource) alarm_attrs() (ret string) {
	ok_actions := get_list_items_by_pattern(t.Attrs, `^ok_actions\.[0-9]+$`)
	alarm_actions := get_list_items_by_pattern(t.Attrs, `^alarm_actions\.[0-9]+$`)
	insufficient_data_actions := get_list_items_by_pattern(t.Attrs, `^insufficient_data_actions\.[0-9]+$`)
	ret += create_ruby_string_array(ok_actions, "  its(:ok_actions) { should == [%v]}\n")
	ret += create_ruby_string_array(alarm_actions, "  its(:alarm_actions) { should == [%v]}\n")
	ret += create_ruby_string_array(insufficient_data_actions, "  its(:insufficient_data_actions) { should == [%v]}\n")
	return
}

func (t Resource) lambda_attrs() (ret string) {
	env_vars_pattern := regexp.MustCompile(`^environment\.[0-9]+\.variables\.([^%]+)$`)
	env_vars := make(map[string]string)
	security_groups := get_list_items_by_pattern(t.Attrs, `^vpc_config\.[0-9]+\.security_group_ids\.([^#]+)$`)
	subnet_ids := get_list_items_by_pattern(t.Attrs, `^vpc_config\.[0-9]+\.subnet_ids\.([^#]+)$`)
	vpc_ids := get_list_items_by_pattern(t.Attrs, `^vpc_config\.[0-9]+\.vpc_id$`)
	(*t.Attrs).ForEach(func(key, value gjson.Result) bool {
		key_string := key.String()
		value_string := value.String()
		if env_vars_pattern.MatchString(key_string) {
			pattern_matches := env_vars_pattern.FindStringSubmatch(key_string)
			env_vars[pattern_matches[1]] = value_string
		}
		return true
	})
	ret = fmt.Sprintf(
		"  its(:environment) { should eq Aws::Lambda::Types::EnvironmentResponse.new({variables: %v }) }\n",
		create_ruby_hash(&env_vars))
	ret += fmt.Sprintf(
		"  its(:vpc_config) { should eq Aws::Lambda::Types::VpcConfigResponse.new({'subnet_ids'=>%v, " +
			"'security_group_ids'=>%v, 'vpc_id'=>'%v'}) }\n",
		create_ruby_array(subnet_ids),
		create_ruby_array(security_groups),
		(*vpc_ids)[0],
			)
	return
}
