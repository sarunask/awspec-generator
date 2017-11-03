package resources

import (
	"fmt"
	"github.com/tidwall/gjson"
	"strings"
)

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
