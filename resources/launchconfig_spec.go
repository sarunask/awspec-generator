package resources

import "fmt"

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

func (t Resource) lc_attrs() (ret string)  {
	//LaunchConfiguration attributes
	ami_id := get_attribute_by_name(t.Attrs, "image_id")
	if ami_id.String() != "" {
		ret += fmt.Sprintf("  its(:image_id) { should eq '%v' }\n", ami_id.String())
	}
	return
}
