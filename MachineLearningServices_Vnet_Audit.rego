package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 condition2
}
condition2 if {
r.type == "azurerm_machine_learning_compute_cluster"
some r.change.after.compute_type in ["AmlCompute","ComputeInstance"]
not condition1
}
condition1 if {
r.change.after.subnet.id
[empty(field('Microsoft.MachineLearningServices/workspaces/computes/subnet.id'))] != true
}