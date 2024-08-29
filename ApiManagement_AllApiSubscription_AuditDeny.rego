package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 LgcQD
}
LgcQD if {
 type == azurerm_api_management_subscription
 regex.match("*/apis",Microsoft.ApiManagement/service/subscriptions/scope)
 r.change.after.state == active
}
