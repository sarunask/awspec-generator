package resources

import (
	"fmt"
	"github.com/tidwall/gjson"
	"strings"
)

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
	ret += create_ruby_string_array(availability_zones, "  its(:availability_zones) { should =~ [%v] }\n")
	ret += create_ruby_string_array(elbs, "  its(:load_balancer_names) { should =~ [%v] }\n")
	ret += create_ruby_string_array(termination_policies, "  its(:termination_policies) { should =~ [%v] }\n")
	return
}

