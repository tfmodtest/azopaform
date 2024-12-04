package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 condition1
}
condition1(x) if {
r.type == "azurerm_application_insights"
r.change.after.disable_local_auth != "true"
}