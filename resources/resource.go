package resources

import (
	"github.com/tidwall/gjson"
	"strings"
	"github.com/sarunask/awspec-generator/loggers"
	"os"
	"io"
	"path/filepath"
	"regexp"
	"math/rand"
	"time"
	"fmt"
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
	//Elastic Load Balancer
	ELB
	//AutoScaling Group
	ASG
	//RDS Instance
	RDS
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
	case ELB:
		return "aws_elb"
	case ASG:
		return "aws_autoscaling_group"
	case RDS:
		return "aws_db_instance"
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
	Tags map[string]*Tag
	//Dependent resources
	Dependent map[int]*Resource
	//Attributes
	Attrs *gjson.Result
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
	case ELB:
		return t.aws_elb_spec()
	case ASG:
		return t.aws_autoscaling_group_spec()
	case RDS:
		return t.aws_rds_instance_spec()
	}
}

// Add resource as dependent to parent
func (t Resource) AddDependency(r *Resource) {
	index := len(t.Dependent)
	t.Dependent[index] = r
}

func (t Resource) Write(folder string) {
	file_name := filepath.Join(folder,
		strings.Replace(t.TerraformName, ".", "-", -1) + "_spec.rb")
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
	res.Tags = make(map[string]*Tag)
	attrs := json.Get("primary.attributes")
	res.Attrs = &attrs
	var alt_name string
	tag_pattern := regexp.MustCompile(`^tag\.([0-9]+)\.(.+)$`)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	(*res.Attrs).ForEach(func(key, value gjson.Result) bool {
		key_string := key.String()
		if tag_pattern.MatchString(key_string) {
			//We have not-random tags from Terraform
			pattern_matches := tag_pattern.FindStringSubmatch(key_string)
			tag := get_tag(&res.Tags, pattern_matches[1])
			switch strings.ToLower(pattern_matches[2]) {
			case "key":
				tag.Name = value.String()
			case "value":
				tag.Value = value.String()
			}
		} else if strings.Index(key_string, "tags.") != -1 {
			tag_name := strings.Replace(key_string, "tags.", "", 1)
			if tag_name != "%" {
				//Other tag format - use random strings
				tag := get_tag(&res.Tags, fmt.Sprintf("%v", r.Int63()))
				tag.Name = tag_name
				tag.Value = value.String()
			}
			return true
		} else if strings.EqualFold(key_string, "name") {
			alt_name = value.String()
		}
		return true
	})
	for _, value := range res.Tags {
		if strings.ToLower(value.Name) == "name" {
			res.Name = value.Value
			break
		}
	}
	if res.Name == "" {
		res.Name = alt_name
	}
	return res
}




