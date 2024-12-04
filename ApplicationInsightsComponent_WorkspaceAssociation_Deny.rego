package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 condition2
}
condition2(x) if {
r.type == "azurerm_application_insights"
not condition1(x)
}
condition1(x) if {
r.change.after.workspace_id != ""
r.change.after.workspace_id
}