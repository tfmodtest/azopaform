package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 condition2
}
condition2 if {
r.type == "azurerm_machine_learning_compute_cluster"
r.change.after.compute_type == "ComputeInstance"
not condition1
}
condition1 if {
r.change.after.idle_time_before_shutdown
[empty(field('Microsoft.MachineLearningServices/workspaces/computes/idleTimeBeforeShutdown'))] != true
}