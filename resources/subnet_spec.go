package resources

import "fmt"

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
