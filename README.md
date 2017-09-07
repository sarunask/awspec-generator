# awspec-generator

Version: 0.1.0

Aim of this project to build a tool, which would be able to take Terraform state JSON and output 
AWSpec infrastructure tests.
For now I have implemented such test resources:
* VPC
* Subnet
* Security Group
* VPN Gateway
* VPN Connection
* Customer Gateway
* ELB
* ASG
* RDS Instance

## Install program
In order to install program, you would need to Install GO language and setup it as described 
[here](https://golang.org/doc/install).
After you install GO, you could do from CLI:
```bash
go get -u github.com/sarunask/awspec-generator
```
If you added GOPATH/bin directory to you PATH, you would be able to run:
```bash
awspec-generator --help
```

## Get Terraform state
In order to get Terraform state output you would have to execute such command:
```bash
terraform state pull > status.json
```  

## Run program
After you got Terraform state, you are ready to create AWSpec SPEC files:
```bash
awspec-generator -json_file status.json
```
If generation was successful, you would see **./spec** directory with several files. 
Now you can install [AWSpec](https://github.com/k1LoW/awspec), move generated .rb files and **./addons/ec2_helper.rb** 
into you own **spec** directory and run tests.

**NOTE: I strongly suggest to use [RVM](https://rvm.io/) to install Ruby 2.4 and create gemset for your AWSpec 
installation.**
