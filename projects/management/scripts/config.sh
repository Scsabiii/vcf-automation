#!/bin/pwsh

param (
	$LocalIP,
	$IP,
	$Gateway,
	$Netmask
)

$ErrorActionPreference = 'Stop'

$esxiserver = Connect-VIServer -Server $LocalIP -User root -Password VMware1!VMware1!
$esxcli = Get-EsxCli -VMhost (Get-VMHost $esxiserver) -V2
try {
    $esxcli.system.maintenanceMode.set.Invoke(@{enable=$true})
} catch {
    echo $_.Exception.Message
}

#--------------------------------------------------------------
# data
#--------------------------------------------------------------
$externalip = $IP
$externalgw = $Gateway
$externalnetmask = $Netmask
$mgmtnetwork = @{name='management-vcf01'; vlan='{{ management_network.vlan_id }}'; interface='{{ management_network.esxi_interface }}'}
$vmnetwork = @{name='VM Network'; vlan='{{ management_network.vlan_id }}'}

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

function Tag-Interface-Management {
    param (
        $InterfaceName
    )
    try {
        $esxcli.network.ip.interface.tag.add.Invoke(@{
            interfacename=$InterfaceName;
            tagname="Management"
        })
    }
    catch {
        if ($_.Exception.Message -notmatch 'Vmknic is already tagged with Management') {
            throw 
        }
    }
}

function Untag-Interface-Management {
    param (
        $InterfaceName
    )
    try {
        $esxcli.network.ip.interface.ipv4.get.Invoke(@{interfacename=$InterfaceName})
    }
    catch {
        if ($_.Exception.Message -match 'Not found') {
            return
        }
    }
    try {
        $esxcli.network.ip.interface.tag.remove.Invoke(@{
            interfacename=$InterfaceName;
            tagname="Management"
        })
    }
    catch {
        if ($_.Exception.Message -notmatch 'Vmknic is not tagged with Management') {
            throw 
        }
    }
}

#--------------------------------------------------------------
# configuration
#--------------------------------------------------------------
# step 1: create tagged port in every network on vSwitch0
Add-PortGroup -PortGroupName $mgmtnetwork['name']
@($mgmtnetwork, $vmnetwork).ForEach({
    $esxcli.network.vswitch.standard.portgroup.set.Invoke(
        @{portgroupname=$_['name']; vlanid=$_['vlan']
    })
})

# step 2: create management network interface
Add-Network-Interface -PortGroupName $mgmtnetwork['name'] -InterfaceName $mgmtnetwork['interface']

# step 3: configure management network interface
$esxcli.network.ip.interface.ipv4.set.Invoke(@{
    interfacename=$mgmtnetwork['interface']
    ipv4=$externalip
    gateway=$externalgw
    netmask=$externalnetmask
    type='static'
})

# step 4: route default network to management network
$esxcli.network.ip.route.ipv4.add.Invoke(@{
    gateway=$externalgw
    network='default'
})

# step 5: tag management interface as Management service
Untag-Interface-Management -InterfaceName 'vmk0'
Tag-Interface-Management -InterfaceName $mgmtnetwork['interface']

# step 6: set domain name and fqdn
$esxcli.system.hostname.set.Invoke(@{domain="vcf01.qa-de-1.cloud.sap"})

# step 7: set search domains:
$esxcli.network.ip.dns.search.remove.Invoke(@{domain="openstack.qa-de-1.cloud.sap"})
$esxcli.network.ip.dns.search.add.Invoke(@{domain="vcf01.qa-de-1.cloud.sap"})

# step 8: set dns server
$esxcli.network.ip.dns.server.remove.Invoke(@{all=$true})
$esxcli.network.ip.dns.server.add.Invoke(@{server="147.204.9.200"})
$esxcli.network.ip.dns.server.add.Invoke(@{server="147.204.9.201"})

# step 9: set ntp server
 $esxcli.system.ntp.set.Invoke(@{
    enabled=$true
     loglevel="warning"
     server=@("147.204.9.202","147.204.9.203","147.204.9.204")
 })

# last step: disable maintenance mode
$esxcli.system.maintenanceMode.set.Invoke(@{enable=$false})
