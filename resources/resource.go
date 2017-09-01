package resources

import (
	"fmt"
	"github.com/tidwall/gjson"
	"strings"
	"github.com/sarunask/awspec-generator/loggers"
	"os"
	"io"
	"path/filepath"
	"sync"
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
	//Terraform Spec name of resource
	TerraformName string
	//Tags
	Tags map[string]string
	//Dependent resources
	Dependent map[int]*Resource
}

//Tree of resources
type ResourcesTree struct {
	Tree map[string] *Resource
	Lock sync.Mutex
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

// Add resource as dependent to parent
func (t Resource) AddDependency(r *Resource) {
	index := len(t.Dependent)
	t.Dependent[index] = r
}

func (t Resource) subnet_vpc_id_should_be() string {
	var ret string
	for i := range t.Dependent {
		switch t.Dependent[i].Type {
		case VPC:
			ret += fmt.Sprintf("  its(:vpc_id) { should eq EC2Helper.GetVPCIdFromName('%v') }\n",
				t.Dependent[i].Name)
		}
	}
	return ret
}

func (t Resource) aws_vpc_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe vpc(\"%v\") do\n"+
		"  it { should exist }\n" +
		"  it { should be_available }\n" +
		t.tags() +
		"  its(:vpc_id) { should eq EC2Helper.GetVPCIdFromName('%v') }\n" +
		"  its('Assigned IGW count'){ expect(EC2Helper.GetIGWsCountForVPCwithName('%v')).to eq 0 }\n" +
		"end\n", t.Name, t.Name, t.Name)
}

func (t Resource) aws_subnet_spec() string {
	return fmt.Sprintf("require 'awspec'\n" +
		"require 'ec2_helper'\n\n" +
		"describe subnet(\"%v\") do\n"+
		"  it { should exist }\n" +
		"  it { should be_available }\n" +
		t.tags() +
	    t.subnet_vpc_id_should_be() +
		"end\n", t.Name)
}

func (t Resource) aws_sg_spec() string {
	return fmt.Sprintf("require 'awspec'\n\n" +
		"describe security_group(\"%v\") do\n"+
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
	file_name := filepath.Join(folder, t.Name + "_spec.rb")
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

func (t *ResourcesTree) Init() {
	if t.Tree == nil {
		t.Lock.Lock()
		t.Tree = make(map[string]*Resource, 100)
		t.Lock.Unlock()
	}
}

func (t *ResourcesTree) Push(r *Resource) {
	//Init tree if empty
	if t.Tree == nil {
		t.Init()
	}
	//Add resource by TerraformName to tree
	_, ok := t.Tree[r.TerraformName]
	if ok == false {
		if t.Tree == nil {
			loggers.Error.Println("Tree is still nil")
		}
		t.Lock.Lock()
		t.Tree[r.TerraformName] = r
		t.Lock.Unlock()
	}
}

func (t *ResourcesTree) Write(dir string, wg *sync.WaitGroup) {
	for i := range t.Tree {
		wg.Add(1)
		go func() {
			defer wg.Done()
			t.Tree[i].Write(dir)
		}()
	}
}