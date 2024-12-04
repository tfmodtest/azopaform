package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 condition2
}
condition2(x) if {
r.type == "azurerm_machine_learning_workspace"
not condition1(x)
}
condition1(x) if {
r.change.after.primary_user_assigned_identity
r.change.after.primary_user_assigned_identity != ""
}