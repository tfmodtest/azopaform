package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 condition2
}
condition2 if {
r.type == "azurerm_machine_learning_workspace"
not condition1
}
condition1 if {

r.change.after.encryption.status == "enabled"
}