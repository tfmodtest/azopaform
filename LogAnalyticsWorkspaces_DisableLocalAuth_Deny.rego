package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 condition1
}
condition1(x) if {
r.type == "azurerm_log_analytics_workspace"
r.change.after.features.disable_local_auth != "true"
}