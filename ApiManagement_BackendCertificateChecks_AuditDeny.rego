package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 BNnqS
}
BNnqS if {
 type == azurerm_api_management_backend
 not wHAsjdf
}
wHAsjdf if {
 r.change.after.properties.tls.validateCertificateChain != false
 r.change.after.properties.tls.validateCertificateName != false
}
