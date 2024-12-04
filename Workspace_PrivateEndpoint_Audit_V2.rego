package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 condition2
}
condition2(x) if {
r.type == "azurerm_machine_learning_workspace"
count({x|r.change.after.private_endpoint_connections[x];condition1(x)}) < 1
}
condition1(x) if {
r.change.after.private_endpoint_connections[x].private_link_service_connection_state.status == "Approved"
}