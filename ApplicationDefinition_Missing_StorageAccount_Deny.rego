package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 condition1
}
condition1(x) if {
r.type == "azurerm_managed_application_definition"
not Microsoft.Solutions/applicationDefinitions/storageAccountId
}