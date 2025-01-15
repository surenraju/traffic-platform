package main

deny[msg] {
    resource := input.resource_changes[_]
    actions := resource.change.actions
    "delete" == actions[_]
    msg = sprintf("Delete operation detected for resource: %s", [resource.address])
}