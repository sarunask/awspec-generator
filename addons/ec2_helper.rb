require 'json'

class EC2Helper
    def self.GetEC2IdFromName(name)
        instances = Array.new
        # Filter the ec2 instances for name and state pending or running
        ec2 = Aws::EC2::Resource.new()
        ec2.instances({filters: [
            {name: 'tag:Name', values: [name]},
            {name: 'instance-state-name', values: [ 'pending', 'running']}
        ]}).each do |i|
            instances.push(i.id)
        end

        # If we found a single instance return it, otherwise throw an error.
        if instances.count == 1 then
            return instances[0]
        elsif instances.count == 0 then
            STDERR.puts 'Error: ' + name + ' Instance not found'
        else
            STDERR.puts 'Error: ' + name + ' more than one running instance exists with that Name'
        end
    end
    def self.GetVPCIdFromName(name)
        vpcs = Array.new
        # Filter the ec2 instances for name and state pending or running
        ec2 = Aws::EC2::Client.new()
        begin
            resp = ec2.describe_vpcs({filters: [
                {name: 'tag:Name', values: [name]}
            ]})
            resp.vpcs.each do |i|
                vpcs.push(i[:vpc_id])
            end
        rescue IPAddr::InvalidAddressError
            cmd = "aws ec2 describe-vpcs --filter 'Name=tag:Name,Values=[#{name}]'"
            resp = JSON.parse(%x[ #{cmd} ])
            resp['Vpcs'].each do |i|
                vpcs.push(i['VpcId'])
            end
        end
        # If we found a single instance return it, otherwise throw an error.
        if vpcs.count == 1 then
            return vpcs[0]
        elsif vpcs.count == 0 then
            STDERR.puts 'Error: ' + name + ' VPC not found'
        else
            STDERR.puts 'Error: ' + name + ' more than one VPC exists with that Name'
        end
    end
    def self.GetIGWsCountForVPCwithName(name)
        igws = Array.new
        # Filter the ec2 instances for name and state pending or running
        ec2 = Aws::EC2::Client.new()
        vpc_id = self.GetVPCIdFromName(name)
        begin
            resp = ec2.describe_internet_gateways({filters: [
                {name: 'attachment.vpc-id', values: [vpc_id]}
            ]})
            resp.internet_gateways.each do |i|
                igws.push(i[:internet_gateway_id])
            end
        rescue IPAddr::InvalidAddressError
            cmd = "aws ec2 describe-internet-gateways --filters 'Name=attachment.vpc-id,Values=#{name}'"
            resp = JSON.parse(%x[ #{cmd} ])
            resp['InternetGateways'].each do |i|
                vpcs.push(i['InternetGatewayId'])
            end
        end
        # If we found a single instance return it, otherwise throw an error.
        return igws.count
    end
    def self.GetVPNGWIdFromName(name)
        vpngws = Array.new
        # Filter the ec2 instances for name and state pending or running
        ec2 = Aws::EC2::Client.new()
        begin
            resp = ec2.describe_vpn_gateways({filters: [
                {name: 'tag:Name', values: [name]}
            ]})
            resp.vpn_gateways.each do |i|
                vpngws.push(i[:vpn_gateway_id])
            end
        rescue IPAddr::InvalidAddressError
            cmd = "aws ec2 describe-vpn-gateways --filters 'Name=tag:Name,Values=#{name}'"
            resp = JSON.parse(%x[ #{cmd} ])
            resp['VpnGateways'].each do |i|
                vpngws.push(i['VpnGatewayId'])
            end
        end
        # If we found a single vpn_gw_id return it, otherwise throw an error.
        if vpngws.count == 1 then
            return vpngws[0]
        elsif vpngws.count == 0 then
            STDERR.puts 'Error: ' + name + ' VPN GW not found'
        else
            STDERR.puts 'Error: ' + name + ' more than one VPN GW exists with that Name'
        end
    end
    def self.GetASGIdFromName(name)
        asgs = Array.new
        # Filter the ec2 instances for name and state pending or running
        autoscale = Aws::AutoScaling::Client.new()
        begin
            resp = autoscale.describe_auto_scaling_groups()
            resp.auto_scaling_groups.each do |i|
                i.tags.each do |tag|
                    if tag.key == "Name" and tag.value == name then
                        asgs.push(i['auto_scaling_group_name'])
                    end
                end
            end
        rescue IPAddr::InvalidAddressError
            cmd = "aws autoscaling describe-auto-scaling-groups"
            resp = JSON.parse(%x[ #{cmd} ])
            resp['AutoScalingGroups'].each do |i|
                i.Tags.each do |tag|
                    if tag.Key == "Name" and tag.Value == name then
                        asgs.push(i['AutoScalingGroupName'])
                    end
                end
            end
        end
        # If we found a single vpn_gw_id return it, otherwise throw an error.
        if asgs.count == 1 then
            return asgs[0]
        elsif asgs.count == 0 then
            STDERR.puts 'Error: ' + name + ' AutoScalingGroup not found'
        else
            STDERR.puts 'Error: ' + name + ' more than one AutoScalingGroup exists with that Name'
        end
    end
    def self.GetRDSIdFromName(name)
        rds = Array.new
        # Filter the ec2 instances for name and state pending or running
        rds_client = Aws::RDS::Client.new()
        begin
            resp = rds_client.describe_db_instances()
            resp.db_instances.each do |i|
                if i.db_instance_arn.include? name
                    rds.push(i.db_instance_identifier)
                end
            end
        rescue IPAddr::InvalidAddressError
            cmd = "aws rds describe-db-instances"
            resp = JSON.parse(%x[ #{cmd} ])
            resp['AutoScalingGroups'].each do |i|
                if i.DBInstanceArn.include? name
                    rds.push(i.DBInstanceIdentifier)
                end
            end
        end
        # If we found a single vpn_gw_id return it, otherwise throw an error.
        if rds.count == 1 then
            return rds[0]
        elsif rds.count == 0 then
            STDERR.puts 'Error: ' + name + ' RDS DB instance not found'
        else
            STDERR.puts 'Error: ' + name + ' more than one RDS DB instance exists with that Name'
        end
    end
end