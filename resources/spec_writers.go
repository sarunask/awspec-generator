package resources

import (
	"fmt"
	"github.com/tidwall/gjson"
	"strings"
)



func (t Resource) get_attribute(attr_name string) (ret string) {
	(*t.Attrs).ForEach(func(key, value gjson.Result) bool {
		if strings.EqualFold(key.String(), attr_name) {
			ret = value.String()
			return false
		}
		return true
	})
	return
}

func (t Resource) tags() (ret string) {
	for _, value := range t.Tags {
		ret += fmt.Sprintf("  it { should have_tag('%v').value('%v') }\n",
			value.Name, value.Value)
	}
	return
}


func (t Resource) asg_dependencies() (ret string) {
	for i := range t.Dependent {
		switch t.Dependent[i].Type {
		case LAUNCH_CONFIG:
			name_prefix := get_attribute_by_name(t.Dependent[i].Attrs, "name_prefix")
			ret += fmt.Sprintf(
				"  it { should have_launch_configuration(EC2Helper.GetLaunchConfigIdFromName('%v')) }\n",
				name_prefix)
		}
	}
	return
}

func (t Resource) sg_dependencies() (ret string) {
	for i := range t.Dependent {
		switch t.Dependent[i].Type {
		case SG:
			ret += fmt.Sprintf("  it { should have_security_group('%v') }\n",
				t.Dependent[i].Name)
		}
	}
	return
}
