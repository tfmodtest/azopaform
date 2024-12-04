package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 condition7
}
condition7(x) if {
r.type == "azurerm_federated_identity_credential"
condition6(x)
}
condition6(x) if {
regex.match("*.oic.prod-aks.azure.com",[if(greaterOrEquals(length(split(field('Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/issuer'),'/')),3),split(field('Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/issuer'),'/')[x])
not condition5(x)
}
condition5 if {
not condition4(x)
}
condition4(x) if {
not condition3(x)
not [if(greaterOrEquals(length(split(field('Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/issuer'),'/')),5),split(field('Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/issuer'),'/')[4],'')] in []
}
condition3(x) if {
not condition1(x)
not condition2(x)
}
condition1(x) if {
count({x|r.change.after.[parameters('allowed_locations')]}) != 0

not [if(greaterOrEquals(length(split(if(greaterOrEquals(length(split(field('Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/issuer'),'/')),3),split(field('Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/issuer'),'/')[2],''), '.')),1),split(if(greaterOrEquals(length(split(field('Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/issuer'),'/')),3),split(field('Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/issuer'),'/')[2],''), '.')[0],'')] in []
}
condition2(x) if {
count({x|r.change.after.[parameters('allowed_tenants')]}) != 0

not [if(greaterOrEquals(length(split(field('Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/issuer'),'/')),4),split(field('Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/issuer'),'/')[3],'')] in []
}