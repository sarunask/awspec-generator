package resources

import (
	"fmt"
	"github.com/tidwall/gjson"
	"strings"
	"github.com/sarunask/awspec-generator/loggers"
	"os"
	"io"
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
	//Tags
	Tags map[string]string
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
	}
}

func (t Resource) aws_vpc_spec() string {
	return fmt.Sprintf("describe vpc(\"%v\") do\n"+
		"  it { should exist }\n" +
		"  it { should be_available }\n" +
		t.tags() +
		"end\n", t.Name)
}

func (t Resource) aws_subnet_spec() string {
	return fmt.Sprintf("describe subnet(\"%v\") do\n"+
		"  it { should exist }\n" +
		"  it { should be_available }\n" +
		t.tags() +
		"end\n", t.Name)
}

func (t Resource) aws_sg_spec() string {
	return fmt.Sprintf("describe security_group(\"%v\") do\n"+
		"  it { should exist }\n" +
		"  it { should be_available }\n" +
		t.tags() +
		"end\n", t.Name)
}

func (t Resource) tags() string {
	var str string
	for key, value := range t.Tags {
		str += fmt.Sprintf("  it { should have_tag('%v').value('%v') }\n", key, value)
	}
	return str
}

func (t Resource) Write(folder string) {
	//string(filepath.Separator)
	//file_name := filepath.Join(folder, t.Name + "_spec.rb")
	file_name := t.Name + "_spec.rb"
	file, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		loggers.Error.Println("Error opening file " + file_name)
	}
	defer file.Close()
	_, err = io.WriteString(file, t.String())
	if err != nil {
		loggers.Error.Println("Error writing to file: ", file_name)
	}
	file.Sync()
}

//Parse would take gjson Resource and would parse it into structure Resource
func Parse(json *gjson.Result) Resource {
	var res Resource
	//Get Types we know
	res_type := json.Get("type")
	switch res_type.String() {
	default:
		res.Type = Unknown
		return res
	case VPC.String():
		res.Type = VPC
	case Subnet.String():
		res.Type = Subnet
	case SG.String():
		res.Type = Subnet
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