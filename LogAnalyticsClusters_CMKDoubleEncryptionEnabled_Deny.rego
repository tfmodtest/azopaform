package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 condition2
}
condition2(x) if {
r.type == "azurerm_log_analytics_cluster"
not condition1(x)
}
condition1 if {

r.change.after.is_double_encryption_enabled == "true"
}