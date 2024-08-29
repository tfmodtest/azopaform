package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 BmHZB
}
BmHZB if {
 type == azurerm_api_management
 GMP
}
GMP if {
 not some r.change.after.properties.sku.name in [[parameters('listOfAllowedSKUs')]]
}
