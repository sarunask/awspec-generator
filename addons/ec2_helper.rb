require 'json'

class EC2Helper
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
end
