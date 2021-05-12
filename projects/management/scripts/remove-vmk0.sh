#!/bin/pwsh
param (
	$HostIP
)

$esxserver = Connect-VIServer -Server $HostIP -User root -Password VMware1!VMware1!
$esxcli = Get-EsxCli -VMhost (Get-VMHost $esxserver) -V2
$esxcli.network.ip.interface.remove.Invoke(@{interfacename="vmk0"})
