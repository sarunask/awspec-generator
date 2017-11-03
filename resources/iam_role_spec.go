package resources

import "fmt"

func (t Resource) aws_iam_role_spec() string {
	name := get_attribute_by_name(t.Attrs,"name_prefix").String()
	if len(name) == 0 {
		name = t.Name
	}
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe iam_role(EC2Helper::GetIAMRoleWhichBeginsWith('%v')) do\n"+
		"  it { should exist }\n" +
		"end\n", name)
}
