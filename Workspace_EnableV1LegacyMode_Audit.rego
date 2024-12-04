package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 condition1
}
condition1 if {
r.type == "azurerm_machine_learning_workspace"
r.change.after.v1_legacy_mode != "true"
}