# azopaform

A Go-based command-line tool that translates Azure Policy definitions from JSON format to Rego language for use with Open Policy Agent (OPA).

## Overview

This tool bridges the gap between Azure's native policy language and the broader cloud-native policy ecosystem. It converts Azure Policy JSON definitions into equivalent Rego policies that can be used with OPA, Gatekeeper, and other policy engines.

## Features

- **Single file or batch processing**: Process individual policy files or entire directories
- **Configurable output**: Customize package names, utility files, and rule naming
- **Comprehensive condition support**: Handles all Azure Policy condition types including:
  - Field comparisons (`equals`, `notEquals`, `in`, `notIn`, etc.)
  - Logical operators (`allOf`, `anyOf`, `not`)
  - Advanced operations (`contains`, `like`, `exists`, `greater`/`less` than)
- **Parameter processing**: Converts policy parameters with types, defaults, and metadata
- **Effect translation**: Maps Azure Policy effects to Rego equivalents
- **Formatted output**: Generates properly formatted, readable Rego code

## Installation

```bash
go install github.com/tfmodtest/azopaform@latest
```

Or clone and build locally:

```bash
git clone https://github.com/tfmodtest/azopaform.git
cd azopaform
go build -o azopaform .
```

## Usage

### Process a single policy file

```bash
./azopaform -path /path/to/policy.json
```

### Process all JSON files in a directory

```bash
./azopaform -dir /path/to/policies/
```

### Advanced options

```bash
./azopaform \
  -path /path/to/policy.json \
  -package "my.policies" \
  -util-file-name "helpers.rego" \
  -generate-rule-name=true
```

## Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-path` | Path to a single policy definition file | |
| `-dir` | Directory containing policy definition files | |
| `-package` | Package name for generated Rego files | `main` |
| `-util-file-name` | Name of the utility Rego file | `util.rego` |
| `-util-library-package-name` | Use external utility package instead of generating util file | |
| `-generate-rule-name` | Use policy display name as rule name | `true` |

## Example

### Input: Azure Policy JSON

```json
{
  "properties": {
    "displayName": "Require HTTPS for storage accounts",
    "policyRule": {
      "if": {
        "allOf": [
          {
            "field": "type",
            "equals": "Microsoft.Storage/storageAccounts"
          },
          {
            "field": "Microsoft.Storage/storageAccounts/supportsHttpsTrafficOnly",
            "notEquals": true
          }
        ]
      },
      "then": {
        "effect": "deny"
      }
    }
  }
}
```

### Output: Generated Rego

```rego
package main

import rego.v1

deny contains msg if {
    input.type == "Microsoft.Storage/storageAccounts"
    input["Microsoft.Storage/storageAccounts/supportsHttpsTrafficOnly"] != true
    msg := "HTTPS is required for storage accounts"
}
```

## Supported Azure Policy Conditions

The tool supports translation of all Azure Policy condition types:

- **Field conditions**: `field`, `value`, `source`
- **Logical operators**: `allOf`, `anyOf`, `not`
- **Comparison operators**: `equals`, `notEquals`, `in`, `notIn`, `contains`, `notContains`
- **Pattern matching**: `like`, `notLike`, `match`, `notMatch`
- **Existence checks**: `exists`, `containsKey`, `notContainsKey`
- **Numeric comparisons**: `greater`, `greaterOrEquals`, `less`, `lessOrEquals`
- **Array operations**: `count`

## Use Cases

- **Kubernetes Policy Migration**: Convert Azure policies for use with Gatekeeper
- **Multi-cloud Governance**: Standardize policies across different cloud platforms
- **CI/CD Integration**: Use OPA in build pipelines with existing Azure policy logic
- **Infrastructure as Code**: Apply consistent policies in Terraform, Pulumi, etc.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Related Projects

- [Open Policy Agent](https://www.openpolicyagent.org/) - Policy engine for cloud native environments
- [Gatekeeper](https://github.com/open-policy-agent/gatekeeper) - OPA for Kubernetes
- [Azure Policy](https://docs.microsoft.com/en-us/azure/governance/policy/) - Azure's native policy service
