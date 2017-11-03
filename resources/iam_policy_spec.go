package resources

import "fmt"

func (t Resource) aws_iam_policy_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe iam_policy('%v') do\n"+
		"  it { should exist }\n" +
		"  it { should be_attachable }\n" +
		"end\n", t.Name)
}

