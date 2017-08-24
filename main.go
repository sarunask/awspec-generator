package main

import (
	"flag"
	"os"
	"log"
	"io/ioutil"
	"github.com/tidwall/gjson"
	"github.com/sarunask/awspec-generator/resources"
)

const (
	SPEC_DIR = "spec"
)

func exit_with_message(msg string, code int) {
	log.Println(msg)
	os.Exit(code)
}

func parse_resource(resource *gjson.Result, c chan resources.Resource) {
	resource.ForEach(func(key, value gjson.Result) bool {
		resource := resources.Parse(&value)
		if resource.Type != resources.Unknown {
			c <- resource
		}
		return true
	})
}

func write_resource(c chan resources.Resource) {
	for res := range c {
		res.Write(SPEC_DIR)
	}
}

func read_terraform_status(status_file string) {
	//Func will read file with name status_file
	//It would search for any resources in modules array
	//Would send resource for further parsing to parse_resource func
	status_json_bytes, err := ioutil.ReadFile(status_file)
	if err != nil {
		log.Printf("Error reading file %v\n", status_file)
		return
	}
	status_json_string := string(status_json_bytes)
	ress := gjson.Get(status_json_string, "modules.#.resources")
	if !ress.Exists() {
		log.Println("Resources are not present in JSON.")
		return
	}
	var ch_resources chan resources.Resource
	ch_resources = make(chan resources.Resource, 100)
	defer close(ch_resources)
	go write_resource(ch_resources)
	ress.ForEach(func(key, value gjson.Result) bool {
		if value.String() != "{}" {
			parse_resource(&value, ch_resources)
		}
		return true
	})
}

func create_spec_dir() {
	spec_dir := "spec"
	if stat, err := os.Stat(spec_dir); err != nil {
		if os.IsNotExist(err) {
			//Doesn't exists - create it
			os.Mkdir(spec_dir, 0755)
		} else {
			//Other error
			exit_with_message("Error getting stats for spec directory", 2)
		}
	} else if ! stat.IsDir() {
		exit_with_message("spec file exists and is not a directory", 3)
	}
}

func main() {
	var json_file string
	flag.StringVar(&json_file, "json_file", "",
		"Path to Terraform JSON status file to parse")
	flag.Parse()
	if json_file == "" {
		exit_with_message("Usage: awspec-generator -json_file json_file_path\nSee more with -h", 1)
	}
	create_spec_dir()
	read_terraform_status(json_file)
}
