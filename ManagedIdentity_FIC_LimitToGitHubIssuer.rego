package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 condition6
}
condition6(x) if {
r.type == "azurerm_federated_identity_credential"
condition5(x)
}
condition5(x) if {
r.change.after.issuer == "https://token.actions.githubusercontent.com"
not condition4(x)
}
condition4 if {
not condition3(x)
}
condition3(x) if {
not condition2(x)
not [if(greaterOrEquals(length(split(field('Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/subject'),':')),2),split(field('Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/subject'),':')[1],'')] in []
}
condition2(x) if {
not condition1(x)
}
condition1(x) if {
count({x|r.change.after.[parameters('allowed_repo_owners')]}) != 0

not [if(greaterOrEquals(length(split(if(greaterOrEquals(length(split(field('Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/subject'),':')),2),split(field('Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/issuer'),':')[1],''), '/')),2),split(if(greaterOrEquals(length(split(field('Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/subject'),':')),2),split(field('Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/subject'),':')[1],''), '/')[0],'')] in []
}