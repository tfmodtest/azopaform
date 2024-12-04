package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 condition1
}
condition1(x) if {
r.type == "azurerm_monitor_scheduled_query_rules_alert"
r.change.after.check_workspace_alerts_storage_configured != "true"
}