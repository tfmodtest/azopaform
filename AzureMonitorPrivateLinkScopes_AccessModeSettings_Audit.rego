package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 condition2
}
condition2(x) if {
r.type == "azurerm_monitor_private_link_scope"
not condition1(x)
}
condition1(x) if {
r.change.after.ingestion_access_mode == "PrivateOnly"
r.change.after.query_access_mode == "PrivateOnly"
}