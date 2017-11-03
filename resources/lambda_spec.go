package resources

import (
	"fmt"
	"regexp"
	"github.com/tidwall/gjson"
)

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
