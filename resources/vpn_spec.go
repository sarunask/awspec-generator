package resources

import "fmt"

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
