package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 condition4
}
condition4(x) if {
not condition1(x)
not condition3(x)
}
condition1(x) if {
r.type == "azurerm_virtual_machine_scale_set_extension"
r.change.after.extensions.publisher == "Microsoft.EnterpriseCloud.Monitoring"
r.change.after.extensions.type == "MicrosoftMonitoringAgent"
}
condition3(x) if {
r.type == "azurerm_linux_virtual_machine_scale_set"
count({x|r.change.after.extension_profile.extensions[x];condition2(x)}) > 0
}
condition2(x) if {
r.change.after.extension_profile.extensions[x].type == "MicrosoftMonitoringAgent"
}