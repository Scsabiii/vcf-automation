#!/bin/pwsh

param (
	$LocalIP,
	$IP,
	$Gateway,
	$Netmask
)

$ErrorActionPreference = 'Stop'

Set-PowerCLIConfiguration -InvalidCertificateAction Ignore -Confirm:$false
Set-PowerCLIConfiguration -Scope User -ParticipateInCEIP $false -Confirm:$false

$esxiserver = Connect-VIServer -Server $LocalIP -User root -Password pss4devus
$esxcli = Get-EsxCli -VMhost (Get-VMHost $esxiserver) -V2

#--------------------------------------------------------------
# data
#--------------------------------------------------------------
$externalip = $IP
$externalgw = $Gateway
$externalnetmask = $Netmask
$networks = @(
{% for network in private_networks -%}
    @{name='{{ network.name }}'; vlan='{{ network.vlan_id }}'; interface='{{ network.esxi_interface }}'}{{ "," if not loop.last }}
{% endfor -%}
)
$mgmtnetwork = @{name='management-vcf01'; vlan='{{ management_network.vlan_id }}'; interface='{{ management_network.esxi_interface }}'}
$vmnetwork = @{name='VM Network'; vlan='{{ management_network.vlan_id }}'}

# @{name='vmotion'; vlan=1000; interface='vmk10'},
# @{name='edgetep'; vlan=1001; interface='vmk11'},
# @{name='hosttep'; vlan=1002; interface='vmk12'},
# @{name='nfs'; vlan=1003; interface='vmk13'},
# @{name='vsan'; vlan=1004; interface='vmk14'},
# @{name='vsanwitness'; vlan=1005; interface='vmk15'}
# $mgmtnetwork = @{name='management-vcf01'; vlan=1007; interface='vmk20'}
# $vmnetwork = @{name='VM Network'; vlan=1007}

#--------------------------------------------------------------
# functions
#--------------------------------------------------------------
# ignore already exists error when creating new object

function Add-PortGroup {
    param (
        $PortGroupName
    )
    try {
        $esxcli.network.vswitch.standard.portgroup.add.Invoke(@{
            portgroupname=$PortGroupName; vswitchname='vSwitch0'
        })
    } catch {
        $msg = "A portgroup with the name $PortGroupName already exists"
        if ($_.Exception.Message -ne $msg) {
            throw
        }
    }
}

function Add-NetStack {
    param (
        $NetStackName
    )
    try {
        $esxcli.network.ip.netstack.add.Invoke(@{netstack=$NetStackName})
    } catch {
        $msg = "Netstack instance '$NetStackName' is already found in kernel"
        if ($_.Exception.Message -notmatch $msg) {
            throw
        }
    }
}

function Add-Network-Interface {
    param (
        $InterfaceName,
        $PortGroupName,
        $NetStackName
    )
    try {
        if ($NetStackName -eq "") {
            $esxcli.network.ip.interface.add.Invoke(@{
                portgroupname=$PortGroupName;
                interfacename=$InterfaceName;
                netstack=$NetStackName
            })
        } else {
            $esxcli.network.ip.interface.add.Invoke(@{
                portgroupname=$PortGroupName;
                interfacename=$InterfaceName;
            })
        }
    }
    catch {
        if ($_.Exception.Message -notmatch 'Already exists') {
            throw
        }
    }
}

#--------------------------------------------------------------
# configuration
#--------------------------------------------------------------
# step 1: create tagged port in every network on vSwitch0
($networks + $mgmtnetwork).ForEach({
    Add-PortGroup -PortGroupName $_['name']
})
($networks + $mgmtnetwork +$vmnetwork).ForEach({
    $esxcli.network.vswitch.standard.portgroup.set.Invoke(
        @{portgroupname=$_['name']; vlanid=$_['vlan']
    })
})

# step 2: create network interface for each port created in step 1
($networks + $mgmtnetwork).ForEach({
    if ($_['name'] -eq 'vmotion') {
        Add-NetStack -NetStackName 'vmotion'
        Add-Network-Interface -PortGroupName $_['name'] -InterfaceName $_['interface'] -NetStackName 'vmotion'
    } else {
        Add-Network-Interface -PortGroupName $_['name'] -InterfaceName $_['interface']
    }
})

# step 3: configure management network interface
$esxcli.network.ip.interface.ipv4.set.Invoke(@{
    interfacename=$mgmtnetwork['interface']
    ipv4=$externalip
    gateway=$externalgw
    netmask=$externalnetmask
    type='static'
})

# step 4  route default network to management network
$esxcli.network.ip.route.ipv4.add.Invoke(@{
    gateway=$externalgw
    network='default'
})
