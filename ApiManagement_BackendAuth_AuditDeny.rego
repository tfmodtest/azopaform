package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 xqxZx
}
xqxZx if {
 type == azurerm_api_management_backend
 r.change.after.url
 r.change.after.protocol == http
 GrBVp
}
GrBVp if {
 not JHSImOM
 not ZANkArq
}
JHSImOM if {
 r.change.after.properties.credentials.certificate
 r.change.after.properties.[length(field('Microsoft.ApiManagement.service.backends.credentials.certificate[*]'))] != 0
}
ZANkArq if {
 r.change.after.credentials[0].authorization[0].scheme
 r.change.after.credentials[0].authorization[0].parameter
}
