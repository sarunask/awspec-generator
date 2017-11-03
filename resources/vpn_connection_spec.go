package resources

import "fmt"

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
