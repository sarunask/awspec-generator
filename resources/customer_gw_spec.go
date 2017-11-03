package resources

import "fmt"

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

