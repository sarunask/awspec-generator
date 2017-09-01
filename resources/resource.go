package resources

import (
	"github.com/tidwall/gjson"
	"strings"
	"github.com/sarunask/awspec-generator/loggers"
	"os"
	"io"
	"path/filepath"
)

type Type int

const (
	//Unknown type of resource
	Unknown Type = iota
	//VPC
	VPC
	//Subnet
	Subnet
	//Security group
	SG
	//VPN Gateway
	VPN_GW
	//VPN Connection
	VPN_CONNECTION
	//Customer VPN GW
	CUSTOMER_GW
)

// String returns a string representation of the type.
func (t Type) String() string {
	switch t {
	default:
		return ""
	case VPC:
		return "aws_vpc"
	case Subnet:
		return "aws_subnet"
	case SG:
		return "aws_security_group"
	case VPN_GW:
		return "aws_vpn_gateway"
	case VPN_CONNECTION:
		return "aws_vpn_connection"
	case CUSTOMER_GW:
		return "aws_customer_gateway"
	}
}

//Resource of Terraform
type Resource struct {
	//RAW string of Terraform json
	Raw string
	//Type of resource
	Type Type
	//Name of resource, should be either Name tag or name
	Name string
	//Terraform Spec name of resource
	TerraformName string
	//Tags
	Tags map[string]string
	//Dependent resources
	Dependent map[int]*Resource
}

// String returns a string representation of the value.
func (t Resource) String() string {
	switch t.Type {
	default:
		return ""
	case VPC:
		return t.aws_vpc_spec()
	case Subnet:
		return t.aws_subnet_spec()
	case SG:
		return t.aws_sg_spec()
	case VPN_GW:
		return t.aws_vpn_gw_spec()
	case VPN_CONNECTION:
		return t.aws_vpn_connection_spec()
	case CUSTOMER_GW:
		return t.aws_customer_gw_spec()
	}
}

// Add resource as dependent to parent
func (t Resource) AddDependency(r *Resource) {
	index := len(t.Dependent)
	t.Dependent[index] = r
}

func (t Resource) Write(folder string) {
	file_name := filepath.Join(folder, t.Name + "-" + t.Type.String() + "_spec.rb")
	str := t.String()
	file, err := os.OpenFile(file_name, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		loggers.Trace.Println("Error opening file " + file_name)
		loggers.Error.Println(err)
	}
	defer file.Close()
	loggers.Info.Printf("Printing to file: %v content:\n%v\n", file_name, str)
	_, err = io.WriteString(file, str)
	if err != nil {
		loggers.Trace.Println("Error writing to file: ", file_name)
		loggers.Error.Println(err)
	}
	file.Sync()
}

//Parse would take gjson Resource and would parse it into structure Resource
func Parse(json *gjson.Result) Resource {
	var res Resource
	//Get Types we know
	res_type := json.Get("type")
	for i := Type(VPC); Type(i).String() != ""; i++ {
		if 	res_type.String() == Type(i).String() {
			res.Type = Type(i)
			break
		}
	}
	if res.Type == Unknown {
		return res
	}
	//Put RAW
	res.Raw = json.String()
	//Make deps
	if res.Dependent == nil {
		res.Dependent = make(map[int]*Resource)
	}
	//Get Attributes
	res.Tags = make(map[string]string)
	attrs := json.Get("primary.attributes")
	attrs.ForEach(func(key, value gjson.Result) bool {
		key_string := key.String()
		if strings.Index(key_string, "tags.") != -1 {
			tag_name := strings.Replace(key_string, "tags.", "", 1)
			if tag_name != "%" {
				res.Tags[tag_name] = value.String()
			}
		}
		return true
	})
	if name, ok := res.Tags["Name"]; ok != false {
		res.Name = name
	}
	return res
}




