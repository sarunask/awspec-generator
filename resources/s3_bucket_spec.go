package resources

import "fmt"

func (t Resource) aws_s3_bucket_spec() string {
	name := get_attribute_by_name(t.Attrs,"bucket_prefix").String()
	if len(name) == 0 {
		name = get_attribute_by_name(t.Attrs,"bucket").String()
	}
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe s3_bucket(EC2Helper::GetS3BucketIdFromName(%v)) do\n"+
		"  it { should exist }\n" +
		"  its(:name) { should start_with(%v) }\n" +
		t.tags() +
		"end\n",
		create_ruby_string(
			name),
		create_ruby_string(
			get_attribute_by_name(t.Attrs,"bucket_prefix").String()),
	)
}
