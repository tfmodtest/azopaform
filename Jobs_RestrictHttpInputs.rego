package main

import rego.v1

r := tfplan.resource_changes[_]

deny if {
 condition21
}
condition21(x) if {
r.type == "azurerm_media_job"
not condition20(x)
}
condition20(x) if {
not condition2(x)
not condition5(x)
count({x|r.change.after.input.#_microsoft-_media-_job_inputs.inputs[x];condition12(x)}) <= 0

count({x|r.change.after.input.#_microsoft-_media-_job_inputs.inputs[x].#_microsoft-_media-_job_inputs.inputs[x];condition19(x)}) <= 0

}
condition2(x) if {
r.change.after.input.#_microsoft-_media-_job_input_http.base_uri
count({x|r.change.after.[parameters('allowed_job_input_http_uri_patterns')];condition1(x)}) == 0
}
condition1(x) if {
regex.match("[current('pattern')]",r.change.after.input.#_microsoft-_media-_job_input_http.base_uri)
}
condition5(x) if {
not r.change.after.input.#_microsoft-_media-_job_input_http.base_uri
count({x|r.change.after.input.#_microsoft-_media-_job_input_clip.files[x];condition4(x)}) > 0
}
condition4(x) if {
count({x|r.change.after.[parameters('allowed_job_input_http_uri_patterns')];condition3(x)}) == 0
}
condition12(x) if {
not condition11(x)
}
condition11(x) if {
not condition7(x)
not condition10(x)
}
condition7(x) if {
r.change.after.input.#_microsoft-_media-_job_inputs.inputs[x].#_microsoft-_media-_job_input_http.base_uri
count({x|r.change.after.[parameters('allowed_job_input_http_uri_patterns')];condition6(x)}) == 0
}
condition6(x) if {
regex.match("[current('pattern')]",r.change.after.input.#_microsoft-_media-_job_inputs.inputs[x].#_microsoft-_media-_job_input_http.base_uri)
}
condition10(x) if {
not r.change.after.input.#_microsoft-_media-_job_inputs.inputs[x].#_microsoft-_media-_job_input_http.base_uri
count({x|r.change.after.input.#_microsoft-_media-_job_inputs.inputs[x].#_microsoft-_media-_job_input_clip.files[x];condition9(x)}) > 0
}
condition9(x) if {
count({x|r.change.after.[parameters('allowed_job_input_http_uri_patterns')];condition8(x)}) == 0
}
condition19(x) if {
not condition18(x)
}
condition18(x) if {
not condition14(x)
not condition17(x)
}
condition14(x) if {
r.change.after.input.#_microsoft-_media-_job_inputs.inputs[x].#_microsoft-_media-_job_inputs.inputs[x].#_microsoft-_media-_job_input_http.base_uri
count({x|r.change.after.[parameters('allowed_job_input_http_uri_patterns')];condition13(x)}) == 0
}
condition13(x) if {
regex.match("[current('pattern')]",r.change.after.input.#_microsoft-_media-_job_inputs.inputs[x].#_microsoft-_media-_job_inputs.inputs[x].#_microsoft-_media-_job_input_http.base_uri)
}
condition17(x) if {
not r.change.after.input.#_microsoft-_media-_job_inputs.inputs[x].#_microsoft-_media-_job_inputs.inputs[x].#_microsoft-_media-_job_input_http.base_uri
count({x|r.change.after.input.#_microsoft-_media-_job_inputs.inputs[x].#_microsoft-_media-_job_inputs.inputs[x].#_microsoft-_media-_job_input_clip.files[x];condition16(x)}) > 0
}
condition16(x) if {
count({x|r.change.after.[parameters('allowed_job_input_http_uri_patterns')];condition15(x)}) == 0
}