package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 condition10
}
condition10(x) if {
r.type == "azurerm_federated_identity_credential"
not condition9(x)
}
condition9(x) if {
[parameters('allowFederatedCredentials')] != false
not condition1(x)
not condition2(x)
not condition3(x)
not condition4(x)
not condition8(x)
}
condition1(x) if {
regex.match("*.oic.prod-aks.azure.com",[if(greaterOrEquals(length(split(field('Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/issuer'),'/')),3),split(field('Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/issuer'),'/')[x])
[parameters('allowAKS')] == false
}
condition2(x) if {
r.change.after.issuer == "https://token.actions.githubusercontent.com"
[parameters('allowGitHub')] == false
}
condition3(x) if {
r.change.after.issuer == "https://cognito-identity.amazonaws.com"
[parameters('allowAWS')] == false
}
condition4(x) if {
r.change.after.issuer == "https://accounts.google.com"
[parameters('allowGCS')] == false
}
condition8(x) if {
not condition6(x)
not condition7(x)
}
condition6 if {
not condition5(x)
}
condition5(x) if {
not regex.match("*.oic.prod-aks.azure.com",[if(greaterOrEquals(length(split(field('Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/issuer'),'/')),3),split(field('Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/issuer'),'/')[x])
r.change.after.issuer != "https://token.actions.githubusercontent.com"
r.change.after.issuer != "https://cognito-identity.amazonaws.com"
r.change.after.issuer != "https://accounts.google.com"
}
condition7 if {

some r.change.after.issuer in []
}