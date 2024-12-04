package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 condition1
}
condition1(x) if {
r.type == "azurerm_media_services_account"
r.change.after.encryption[x].type != "CustomerKey"
}