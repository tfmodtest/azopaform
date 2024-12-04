package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 condition1
}
condition1(x) if {
r.type == "azurerm_dashboard_grafana"
r.change.after.public_network_access != "Disabled"
}