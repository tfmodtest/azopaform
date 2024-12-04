package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 condition1
}
condition1(x) if {
r.type == "azurerm_virtual_machine_extension"
r.change.after.publisher == "Microsoft.EnterpriseCloud.Monitoring"
r.change.after.type == "MicrosoftMonitoringAgent"
}