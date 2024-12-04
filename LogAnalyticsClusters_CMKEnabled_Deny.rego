package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 condition5
}
condition5(x) if {
r.type == "azurerm_log_analytics_cluster"
not condition4(x)
}
condition4(x) if {
condition1(x)
condition2(x)
condition3(x)
}
condition1(x) if {
r.change.after.key_vault_properties.key_vault_uri != ""
r.change.after.key_vault_properties.key_vault_uri
}
condition2(x) if {
r.change.after.key_vault_properties.key_name != ""
r.change.after.key_vault_properties.key_name
}
condition3 if {

r.change.after.key_vault_properties.key_version
}