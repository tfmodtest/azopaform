package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 condition3
}
condition3(x) if {
r.type == "azurerm_maps_account"
not condition2(x)
}
condition2(x) if {
Microsoft.Maps/accounts/cors.corsRules[x].allowedOrigins
count({x|Microsoft.Maps/accounts/cors.corsRules[x].allowedOrigins[x];condition1(x)}) <= 0

}
condition1(x) if {
some Microsoft.Maps/accounts/cors.corsRules[x].allowedOrigins[x] in ["*",""]
}