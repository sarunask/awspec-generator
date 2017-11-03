package resources

import (
	"fmt"
	"github.com/tidwall/gjson"
	"strings"
)

func (t Resource) aws_rds_instance_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe rds(EC2Helper.GetRDSIdFromName('%v')) do\n"+
		"  it { should exist }\n" +
		"  it { should be_available }\n" +
		t.tags() +
		t.rds_attrs() +
		t.sg_dependencies() +
		"end\n", t.Name)
}

func (t Resource) rds_attrs() (ret string)  {
	//RDS DB instance attributes
	(*t.Attrs).ForEach(func(key, value gjson.Result) bool {
		key_string := key.String()
		value_string := value.String()
		switch strings.ToLower(key_string) {
		case "option_group_name":
			ret += fmt.Sprintf("  it { should have_option_group('%v')}\n",
				value_string)
		case "instance_class":
			ret += fmt.Sprintf("  its(:db_instance_class) { should eq '%v' }\n",
				value_string)
		case "engine":
			ret += fmt.Sprintf("  its(:engine) { should eq '%v' }\n",
				value_string)
		case "engine_version":
			ret += fmt.Sprintf("  its(:engine_version) { should eq '%v' }\n",
				value_string)
		case "db_instance_class":
			ret += fmt.Sprintf("  its(:db_instance_class) { should eq '%v' }\n",
				value_string)
		case "username":
			ret += fmt.Sprintf("  its(:master_username) { should eq '%v' }\n",
				value_string)
		case "name":
			ret += fmt.Sprintf("  its(:db_name) { should eq '%v' }\n",
				value_string)
		case "allocated_storage":
			ret += fmt.Sprintf("  its(:allocated_storage) { should eq %v }\n",
				value_string)
		case "availability_zone":
			ret += fmt.Sprintf("  its(:availability_zone) { should eq '%v' }\n",
				value_string)
		case "backup_retention_period":
			ret += fmt.Sprintf("  its(:backup_retention_period) { should eq %v }\n",
				value_string)
		case "maintenance_window":
			ret += fmt.Sprintf("  its(:preferred_maintenance_window) { should eq '%v' }\n",
				value_string)
		case "backup_window":
			ret += fmt.Sprintf("  its(:preferred_backup_window) { should eq '%v' }\n",
				value_string)
		case "multi_az":
			ret += fmt.Sprintf("  its(:multi_az) { should eq %v }\n",
				value_string)
		case "publicly_accessible":
			ret += fmt.Sprintf("  its(:publicly_accessible) { should eq %v }\n",
				value_string)
		case "auto_minor_version_upgrade":
			ret += fmt.Sprintf("  its(:auto_minor_version_upgrade) { should eq %v }\n",
				value_string)
		case "storage_type":
			ret += fmt.Sprintf("  its(:storage_type) { should eq '%v' }\n",
				value_string)
		case "storage_encrypted":
			ret += fmt.Sprintf("  its(:storage_encrypted) { should eq %v }\n",
				value_string)
		case "kms_key_id":
			ret += fmt.Sprintf("  its(:kms_key_id) { should eq '%v' }\n",
				value_string)
		case "copy_tags_to_snapshot":
			ret += fmt.Sprintf("  its(:copy_tags_to_snapshot) { should eq %v }\n",
				value_string)
		case "monitoring_interval":
			ret += fmt.Sprintf("  its(:monitoring_interval) { should eq %v }\n",
				value_string)
		}
		return true
	})
	return
}

