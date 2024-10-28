package pkg

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/require"
)

const count_json = `{
  "properties": {
    "displayName": "App Service Environment should have TLS 1.0 and 1.1 disabled",
    "policyType": "BuiltIn",
    "mode": "Indexed",
    "description": "TLS 1.0 and 1.1 are out-of-date protocols that do not support modern cryptographic algorithms. Disabling inbound TLS 1.0 and 1.1 traffic helps secure apps in an App Service Environment.",
    "metadata": {
      "version": "2.0.1",
      "category": "App Service"
    },
    "version": "2.0.1",
    "parameters": {
      "effect": {
        "type": "string",
        "defaultValue": "Audit",
        "allowedValues": [
          "Audit",
          "Deny",
          "Disabled"
        ],
        "metadata": {
          "displayName": "Effect",
          "description": "Enable or disable the execution of the policy"
        }
      }
    },
    "policyRule": {
      "if": {
        "allOf": [
          {
            "field": "type",
            "equals": "Microsoft.Web/hostingEnvironments"
          },
          {
            "field": "kind",
            "like": "ASE*"
          },
          {
            "count": {
              "field": "Microsoft.Web/HostingEnvironments/clusterSettings[*]",
              "where": {
                "allOf": [
                  {
                    "field": "Microsoft.Web/HostingEnvironments/clusterSettings[*].name",
                    "equals": "DisableTls1.0"
                  },
                  {
                    "field": "Microsoft.Web/HostingEnvironments/clusterSettings[*].value",
                    "equals": "1"
                  }
                ]
              }
            },
            "less": 1
          }
        ]
      },
      "then": {
        "effect": "[parameters('effect')]"
      }
    }
  },
  "id": "/providers/Microsoft.Authorization/policyDefinitions/d6545c6b-dd9d-4265-91e6-0b451e2f1c50",
  "name": "d6545c6b-dd9d-4265-91e6-0b451e2f1c50"
}`

const deny_json = `{
  "properties": {
    "displayName": "App Service apps should use a SKU that supports private link",
    "policyType": "BuiltIn",
    "mode": "Indexed",
    "description": "With supported SKUs, Azure Private Link lets you connect your virtual network to Azure services without a public IP address at the source or destination. The Private Link platform handles the connectivity between the consumer and services over the Azure backbone network. By mapping private endpoints to apps, you can reduce data leakage risks. Learn more about private links at: https://aka.ms/private-link.",
    "metadata": {
      "version": "4.1.0",
      "category": "App Service"
    },
    "version": "4.1.0",
    "parameters": {
      "effect": {
        "type": "String",
        "metadata": {
          "displayName": "Effect",
          "description": "Enable or disable the execution of the policy"
        },
        "allowedValues": [
          "Audit",
          "Deny",
          "Disabled"
        ],
        "defaultValue": "Audit"
      }
    },
    "policyRule": {
      "if": {
        "allOf": [
          {
            "field": "type",
            "equals": "Microsoft.Web/serverFarms"
          },
          {
            "allOf": [
              {
                "field": "Microsoft.Web/serverFarms/sku.tier",
                "notIn": [
                  "Basic",
                  "Standard",
                  "ElasticPremium",
                  "Premium",
                  "PremiumV2",
                  "Premium0V3",
                  "PremiumV3",
                  "PremiumMV3",
                  "Isolated",
                  "IsolatedV2",
                  "WorkflowStandard"
                ]
              },
              {
                "field": "Microsoft.Web/serverFarms/sku.name",
                "notIn": [
                  "B1",
                  "B2",
                  "B3",
                  "S1",
                  "S2",
                  "S3",
                  "EP1",
                  "EP2",
                  "EP3",
                  "P1",
                  "P2",
                  "P3",
                  "P1V2",
                  "P2V2",
                  "P3V2",
                  "P0V3",
                  "P1V3",
                  "P2V3",
                  "P3V3",
                  "P1MV3",
                  "P2MV3",
                  "P3MV3",
                  "P4MV3",
                  "P5MV3",
                  "I1",
                  "I2",
                  "I3",
                  "I1V2",
                  "I2V2",
                  "I3V2",
                  "I4V2",
                  "I5V2",
                  "I6V2",
                  "WS1",
                  "WS2",
                  "WS3"
                ]
              }
            ]
          }
        ]
      },
      "then": {
        "effect": "[parameters('effect')]"
      }
    }
  },
  "id": "/providers/Microsoft.Authorization/policyDefinitions/546fe8d2-368d-4029-a418-6af48a7f61e5",
  "name": "546fe8d2-368d-4029-a418-6af48a7f61e5"
}`

const nested_json = `{
    "properties": {
        "displayName": "[AzCliTools][SFI-2.5.1] Disable all high risk inbound ports for Network Security Group",
        "description": "This policy disables all high risk inbound ports on the network security group.",
        "mode": "All",
        "metadata": {
            "version": "1.0.0",
            "category": "SFI"
        },
        "policyRule": {
            "if": {
                "allOf": [
                    {
                        "anyOf": [
                            {
                                "field": "type",
                                "equals": "Microsoft.Network/networkSecurityGroups"
                            },
                            {
                                "field": "type",
                                "equals": "Microsoft.Network/networkSecurityGroups/securityRules"
                            }
                        ]
                    },
                    {
                        "anyOf": [
                            {
                                "allOf": [
                                    {
                                        "field": "type",
                                        "equals": "Microsoft.Network/networkSecurityGroups/securityRules"
                                    },
                                    {
                                        "field": "Microsoft.Network/networkSecurityGroups/securityRules/access",
                                        "equals": "Allow"
                                    },
                                    {
                                        "field": "Microsoft.Network/networkSecurityGroups/securityRules/direction",
                                        "equals": "Inbound"
                                    },
                                    {
                                        "field": "Microsoft.Network/networkSecurityGroups/securityRules/sourceAddressPrefix",
                                        "in": [
                                            "*",
                                            "Internet"
                                        ]
                                    },
                                    {
                                        "anyOf": [
                                            {
                                                "field": "Microsoft.Network/networkSecurityGroups/securityRules/destinationPortRange",
                                                "in": [
                                                    "*",
                                                    "22",
                                                    "3389"
                                                ]
                                            },
                                            {
                                                "not": {
                                                    "field": "Microsoft.Network/networkSecurityGroups/securityRules/destinationPortRanges[*]",
                                                    "notIn": [
                                                        "*",
                                                        "22",
                                                        "3389"
                                                    ]
                                                }
                                            }
                                        ]
                                    }
                                ]
                            },
                            {
                                "count": {
                                    "field": "Microsoft.Network/networkSecurityGroups/securityRules[*]",
                                    "where": {
                                        "allOf": [
                                            {
                                                "field": "Microsoft.Network/networkSecurityGroups/securityRules[*].access",
                                                "equals": "Allow"
                                            },
                                            {
                                                "field": "Microsoft.Network/networkSecurityGroups/securityRules[*].direction",
                                                "equals": "Inbound"
                                            },
                                            {
                                                "field": "Microsoft.Network/networkSecurityGroups/securityRules[*].sourceAddressPrefix",
                                                "in": [
                                                    "*",
                                                    "Internet"
                                                ]
                                            },
                                            {
                                                "field": "Microsoft.Network/networkSecurityGroups/securityRules[*].destinationPortRange",
                                                "in": [
                                                    "*",
                                                    "22",
                                                    "3389"
                                                ]
                                            }
                                        ]
                                    }
                                },
                                "greater": 0
                            }
                        ]
                    }
                ]
            },
            "then": {
                "effect": "deny"
            }
        }
    }
}`

const nested_json2 = `{
  "properties": {
    "displayName": "Configure disaster recovery on virtual machines by enabling replication via Azure Site Recovery",
    "policyType": "BuiltIn",
    "mode": "Indexed",
    "description": "Virtual machines without disaster recovery configurations are vulnerable to outages and other disruptions. If the virtual machine does not already have disaster recovery configured, this would initiate the same by enabling replication using preset configurations to facilitate business continuity.  You can optionally include/exclude virtual machines containing a specified tag to control the scope of assignment. To learn more about disaster recovery, visit https://aka.ms/asr-doc.",
    "metadata": {
      "version": "2.0.0",
      "category": "Compute"
    },
    "version": "2.0.0",
    "parameters": {
      "sourceRegion": {
        "type": "String",
        "metadata": {
          "displayName": "Source Region",
          "description": "Region in which the source virtual machine is deployed",
          "strongType": "location",
          "serviceName": "ASR"
        }
      },
      "targetRegion": {
        "type": "String",
        "metadata": {
          "displayName": "Target Region",
          "description": "Region to be used to deploy the virtual machine in case of a disaster",
          "strongType": "location",
          "serviceName": "ASR"
        }
      },
      "targetResourceGroupId": {
        "type": "String",
        "metadata": {
          "displayName": "Target Resource Group",
          "description": "Resource group to be used to create the virtual machine in the target region",
          "strongType": "Microsoft.Resources/resourceGroups",
          "assignPermissions": true,
          "serviceName": "ASR"
        }
      },
      "vaultResourceGroupId": {
        "type": "String",
        "metadata": {
          "displayName": "Vault Resource Group",
          "description": "The resource group containing the recovery services vault used for disaster recovery configurations",
          "strongType": "Microsoft.Resources/resourceGroups",
          "assignPermissions": true,
          "serviceName": "ASR"
        }
      },
      "vaultId": {
        "type": "String",
        "metadata": {
          "displayName": "Recovery Services Vault",
          "description": "Recovery services vault to be used for disaster recovery configurations",
          "strongType": "Microsoft.RecoveryServices/vaults",
          "serviceName": "ASR"
        }
      },
      "recoveryNetworkId": {
        "type": "String",
        "metadata": {
          "displayName": "Recovery Virtual Network",
          "description": "Id of an existing virtual network in the target region or name of the virtual network to be created in target region",
          "strongType": "Microsoft.Network/virtualNetworks",
          "serviceName": "ASR"
        },
        "defaultValue": ""
      },
      "targetZone": {
        "type": "String",
        "metadata": {
          "displayName": "Target Availability Zone",
          "description": "Availability zone in the designated target region to be used by virtual machines during disaster",
          "strongType": "zone",
          "serviceName": "ASR"
        },
        "defaultValue": ""
      },
      "cacheStorageAccountId": {
        "type": "String",
        "metadata": {
          "displayName": "Cache storage account",
          "description": "Existing cache storage account ID or prefix for the cache storage account name to be created in source region.",
          "strongType": "Microsoft.Storage/storageAccounts",
          "serviceName": "ASR"
        },
        "defaultValue": ""
      },
      "tagName": {
        "type": "String",
        "metadata": {
          "displayName": "Tag Name",
          "description": "Name of the tag to use for including or excluding VMs in the scope of this policy. This should be used along with the tag value parameter.",
          "serviceName": "ASR"
        },
        "defaultValue": ""
      },
      "tagValue": {
        "type": "Array",
        "metadata": {
          "displayName": "Tag Values",
          "description": "Values of the tag to use for including or excluding VMs in the scope of this policy (in case of multiple values, use a comma-separated list). This should be used along with the tag name parameter.",
          "serviceName": "ASR"
        },
        "defaultValue": []
      },
      "tagType": {
        "type": "String",
        "metadata": {
          "displayName": "Tag Type",
          "description": "Tag type can be either Inclusion Tag or Exclusion Tag. Inclusion tag type will make sure VMs with tag name and tag value are included in replication, Exclusion tag type will make sure VMs with tag name and tag value are excluded from replication.",
          "serviceName": "ASR"
        },
        "allowedValues": [
          "Inclusion",
          "Exclusion",
          ""
        ],
        "defaultValue": ""
      },
      "effect": {
        "type": "String",
        "metadata": {
          "displayName": "Effect",
          "description": "Enable or disable the execution of the policy"
        },
        "allowedValues": [
          "DeployIfNotExists",
          "Disabled"
        ],
        "defaultValue": "Disabled"
      }
    },
    "policyRule": {
      "if": {
        "allOf": [
          {
            "field": "type",
            "equals": "Microsoft.Compute/virtualMachines"
          },
          {
            "field": "location",
            "equals": "[parameters('sourceRegion')]"
          },
          {
            "anyOf": [
              {
                "allOf": [
                  {
                    "value": "[parameters('tagType')]",
                    "equals": "Inclusion"
                  },
                  {
                    "field": "[concat('tags[', parameters('tagName'), ']')]",
                    "in": "[parameters('tagValue')]"
                  }
                ]
              },
              {
                "allOf": [
                  {
                    "value": "[parameters('tagType')]",
                    "equals": "Exclusion"
                  },
                  {
                    "field": "[concat('tags[', parameters('tagName'), ']')]",
                    "notIn": "[parameters('tagValue')]"
                  }
                ]
              },
              {
                "anyOf": [
                  {
                    "value": "[empty(parameters('tagName'))]",
                    "equals": "true"
                  },
                  {
                    "value": "[empty(parameters('tagValue'))]",
                    "equals": "true"
                  },
                  {
                    "value": "[empty(parameters('tagType'))]",
                    "equals": "true"
                  }
                ]
              }
            ]
          }
        ]
      },
      "then": {
        "effect": "[parameters('effect')]",
        "details": {
          "type": "Microsoft.Resources/links",
          "evaluationDelay": "PT15M",
          "roleDefinitionIds": [
            "/providers/Microsoft.Authorization/roleDefinitions/8e3af657-a8ff-443c-a75c-2fe8c4bcb635"
          ],
          "existenceCondition": {
            "allOf": [
              {
                "field": "name",
                "like": "ASR-Policy-Protect-*"
              },
              {
                "field": "Microsoft.Resources/links/targetId",
                "contains": "/replicationProtectedItems/"
              }
            ]
          },
          "deployment": {
            "properties": {
              "mode": "incremental",
              "template": {
                "$schema": "http://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
                "contentVersion": "1.0.0.0",
                "parameters": {
                  "apiVersion": {
                    "type": "String"
                  },
                  "avSetId": {
                    "type": "String"
                  },
                  "dataDiskIds": {
                    "type": "object"
                  },
                  "osDiskId": {
                    "type": "String"
                  },
                  "ppgId": {
                    "type": "String"
                  },
                  "recoveryNetworkId": {
                    "type": "String"
                  },
                  "recoverySubscriptionId": {
                    "type": "String"
                  },
                  "sourceRegion": {
                    "type": "String"
                  },
                  "sourceResourceGroupName": {
                    "type": "String"
                  },
                  "targetRegion": {
                    "type": "String"
                  },
                  "targetResourceGroupName": {
                    "type": "String"
                  },
                  "targetZone": {
                    "type": "String"
                  },
                  "vaultName": {
                    "type": "String"
                  },
                  "vaultResourceGroupName": {
                    "type": "String"
                  },
                  "vmId": {
                    "type": "String"
                  },
                  "vmZones": {
                    "type": "Object"
                  },
                  "cacheStorageAccountId": {
                    "type": "String"
                  }
                },
                "variables": {
                  "avSetApiVersion": "2019-03-01",
                  "deploymentApiVersion": "2017-05-10",
                  "vmApiVersion": "2019-07-01",
                  "ppgApiVersion": "2019-12-01",
                  "storageAccountApiVersion": "2018-07-01",
                  "portalLinkPrefix": "https://portal.azure.com/#@microsoft.onmicrosoft.com/resource",
                  "schemaLink": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
                  "defaultAvSet": "defaultAvSet-asr",
                  "defaultPPG": "defaultPPG-asr",
                  "eligibilityResultsDefault": "default",
                  "protectedItemSuffix": "-policy",
                  "recoveryAvSetPrefix": "RecoveryAvSet-",
                  "recoveryPPGPrefix": "RecoveryPPG-",
                  "storagePrefix": "Storage-",
                  "avSetType": "Microsoft.Compute/availabilitySets",
                  "deploymentType": "Microsoft.Resources/deployments",
                  "networkType": "Microsoft.Network/virtualNetworks",
                  "ppgType": "Microsoft.Compute/proximityPlacementGroups",
                  "replicationEligibilityResultsType": "Microsoft.RecoveryServices/replicationEligibilityResults",
                  "storageType": "Microsoft.Storage/storageAccounts",
                  "vaultType": "Microsoft.RecoveryServices/vaults",
                  "avSetTemplateName": "[concat(variables('recoveryAvSetPrefix'), last(split(parameters('vmId'), '/')))]",
                  "avSetTemplateName64": "[if(greater(length(variables('avSetTemplateName')), 64), substring(variables('avSetTemplateName'), 0, 64), variables('avSetTemplateName'))]",
                  "ppgTemplateName": "[concat(variables('recoveryPPGPrefix'), last(split(parameters('vmId'), '/')))]",
                  "ppgTemplateName64": "[if(greater(length(variables('ppgTemplateName')), 64), substring(variables('ppgTemplateName'), 0, 64), variables('ppgTemplateName'))]",
                  "storageAccountTemplateName": "[concat(variables('storagePrefix'), last(split(parameters('vmId'), '/')))]",
                  "storageAccountTemplateName64": "[concat(variables('storagePrefix'), uniqueString(variables('storageAccountTemplateName')))]",
                  "replicationProtectedIntentTemplateName": "[concat('ASR-', parameters('sourceResourceGroupName'), '-', last(split(parameters('vmId'), '/')))]",
                  "replicationProtectedIntentTemplateName64": "[if(greater(length(variables('replicationProtectedIntentTemplateName')), 64), substring(variables('replicationProtectedIntentTemplateName'), 0, 64), variables('replicationProtectedIntentTemplateName'))]",
                  "vmDataDiskIds": "[array(parameters('dataDiskIds').rawValue)]",
                  "vmDiskCount": "[add(length(variables('vmDataDiskIds')), int(1))]",
                  "diskIds": "[concat(array(parameters('osDiskId')), array(parameters('dataDiskIds').rawValue))]",
                  "vaultId": "[resourceId(parameters('vaultResourceGroupName'), variables('vaultType'), parameters('vaultName'))]",
                  "eligibilityResultsId": "[extensionResourceId(parameters('vmId'), variables('replicationEligibilityResultsType'), variables('eligibilityResultsDefault'))]",
                  "protectedIntentName": "[concat(parameters('vaultName'), '/', guid(resourceGroup().id, last(split(parameters('vmId'), '/'))), variables('protectedItemSuffix'))]",
                  "recoveryAvSetName": "[if(empty(parameters('avSetId')), variables('defaultAvSet'), concat(last(split(parameters('avSetId'), '/')), '-asr'))]",
                  "recoveryAvSetId": "[if(empty(parameters('avSetId')), '', resourceId(parameters('targetResourceGroupName'), variables('avSetType'), variables('recoveryAvSetName')))]",
                  "recoveryAvType": "[if(not(empty(parameters('avSetId'))), 'AvailabilitySet', if(not(empty(parameters('targetZone'))), 'AvailabilityZone', 'Single'))]",
                  "recoveryAvZone": "[parameters('targetZone')]",
                  "recoveryPPGName": "[if(empty(parameters('ppgId')), variables('defaultPPG'), concat(last(split(parameters('ppgId'), '/')), '-asr'))]",
                  "recoveryPPGId": "[if(empty(parameters('ppgId')), '', resourceId(parameters('targetResourceGroupName'), variables('ppgType'), variables('recoveryPPGName')))]",
                  "targetResourceGroupId": "[concat('/subscriptions/', parameters('recoverySubscriptionId'), '/resourceGroups/', parameters('targetResourceGroupName'))]",
                  "storageAccountSKUName": "Standard_LRS",
                  "storageAccountKind": "Storage",
                  "cacheStorageAccountArmId": "[if(empty(parameters('cacheStorageAccountId')),'',if(contains(parameters('cacheStorageAccountId'),'/'),parameters('cacheStorageAccountId'),resourceId(parameters('vaultResourceGroupName'), variables('storageType'), parameters('cacheStorageAccountId'))))]"
                },
                "resources": [
                  {
                    "condition": "[not(empty(parameters('ppgId')))]",
                    "apiVersion": "[variables('deploymentApiVersion')]",
                    "name": "[variables('ppgTemplateName64')]",
                    "type": "Microsoft.Resources/deployments",
                    "resourceGroup": "[parameters('targetResourceGroupName')]",
                    "properties": {
                      "mode": "Incremental",
                      "template": {
                        "$schema": "[variables('schemaLink')]",
                        "contentVersion": "1.0.0.0",
                        "parameters": {},
                        "variables": {},
                        "resources": [
                          {
                            "condition": "[not(empty(parameters('ppgId')))]",
                            "type": "[variables('ppgType')]",
                            "name": "[variables('recoveryPPGName')]",
                            "apiVersion": "[variables('ppgApiVersion')]",
                            "location": "[parameters('targetRegion')]",
                            "properties": {
                              "proximityPlacementGroupType": "[if(empty(parameters('ppgId')), 'Standard', reference(parameters('ppgId'), variables('ppgApiVersion')).proximityPlacementGroupType)]"
                            }
                          }
                        ]
                      },
                      "parameters": {}
                    }
                  },
                  {
                    "condition": "[not(empty(parameters('avSetId')))]",
                    "apiVersion": "[variables('deploymentApiVersion')]",
                    "name": "[variables('avSetTemplateName64')]",
                    "type": "Microsoft.Resources/deployments",
                    "resourceGroup": "[parameters('targetResourceGroupName')]",
                    "properties": {
                      "mode": "Incremental",
                      "template": {
                        "$schema": "[variables('schemaLink')]",
                        "contentVersion": "1.0.0.0",
                        "parameters": {},
                        "variables": {},
                        "resources": [
                          {
                            "condition": "[not(empty(parameters('avSetId')))]",
                            "type": "[variables('avSetType')]",
                            "sku": {
                              "name": "[if(empty(parameters('avSetId')), 'Aligned', reference(parameters('avSetId'), variables('avSetApiVersion'), 'Full').sku.name)]"
                            },
                            "name": "[variables('recoveryAvSetName')]",
                            "apiVersion": "[variables('avSetApiVersion')]",
                            "location": "[parameters('targetRegion')]",
                            "tags": {},
                            "properties": {
                              "platformUpdateDomainCount": "[if(empty(parameters('avSetId')), '5', reference(parameters('avSetId'), variables('avSetApiVersion')).platformUpdateDomainCount)]",
                              "platformFaultDomainCount": "[if(empty(parameters('avSetId')), '2', reference(parameters('avSetId'), variables('avSetApiVersion')).platformFaultDomainCount)]",
                              "proximityPlacementGroup": "[if(empty(parameters('ppgId')), json('null'), json(concat('{', '\"id\"', ':', '\"', variables('recoveryPPGId'), '\"', '}')))]"
                            }
                          }
                        ]
                      },
                      "parameters": {}
                    },
                    "dependsOn": [
                      "[variables('ppgTemplateName64')]"
                    ]
                  },
                  {
                    "condition": "[and(not(empty(parameters('cacheStorageAccountId'))), not(contains(parameters('cacheStorageAccountId'), '/')))]",
                    "apiVersion": "[variables('deploymentApiVersion')]",
                    "name": "[variables('storageAccountTemplateName64')]",
                    "type": "Microsoft.Resources/deployments",
                    "resourceGroup": "[parameters('vaultResourceGroupName')]",
                    "properties": {
                      "mode": "Incremental",
                      "template": {
                        "$schema": "[variables('schemaLink')]",
                        "contentVersion": "1.0.0.0",
                        "parameters": {},
                        "variables": {},
                        "resources": [
                          {
                            "condition": "[and(not(empty(parameters('cacheStorageAccountId'))), not(contains(parameters('cacheStorageAccountId'), '/')))]",
                            "type": "[variables('storageType')]",
                            "name": "[parameters('cacheStorageAccountId')]",
                            "apiVersion": "[variables('storageAccountApiVersion')]",
                            "location": "[parameters('sourceRegion')]",
                            "sku": {
                              "name": "[variables('storageAccountSKUName')]"
                            },
                            "kind": "[variables('storageAccountKind')]",
                            "properties": {
                              "supportsHttpsTrafficOnly": true
                            }
                          }
                        ]
                      },
                      "parameters": {}
                    }
                  },
                  {
                    "apiVersion": "[variables('deploymentApiVersion')]",
                    "name": "[variables('replicationProtectedIntentTemplateName64')]",
                    "type": "Microsoft.Resources/deployments",
                    "resourceGroup": "[parameters('vaultResourceGroupName')]",
                    "properties": {
                      "mode": "Incremental",
                      "template": {
                        "$schema": "[variables('schemaLink')]",
                        "contentVersion": "1.0.0.0",
                        "parameters": {},
                        "variables": {},
                        "resources": [
                          {
                            "condition": "[lessOrEquals(length(reference(variables('eligibilityResultsId'), '2018-07-10').errors), int('0'))]",
                            "type": "Microsoft.RecoveryServices/vaults/replicationProtectionIntents",
                            "name": "[variables('protectedIntentName')]",
                            "apiVersion": "[parameters('apiVersion')]",
                            "properties": {
                              "providerSpecificDetails": {
                                "instanceType": "A2A",
                                "fabricObjectId": "[parameters('vmId')]",
                                "primaryLocation": "[parameters('sourceRegion')]",
                                "recoveryLocation": "[parameters('targetRegion')]",
                                "recoverySubscriptionId": "[parameters('recoverySubscriptionId')]",
                                "recoveryAvailabilityType": "[variables('recoveryAvType')]",
                                "recoveryAvailabilityZone": "[variables('recoveryAvZone')]",
                                "recoveryResourceGroupId": "[variables('targetResourceGroupId')]",
                                "recoveryAvailabilitySetCustomInput": "[if(empty(parameters('avSetId')), json('null'), json(concat('{', '\"resourceType\"', ':', '\"Existing\",', '\"recoveryAvailabilitySetId\"', ':', '\"', variables('recoveryAvSetId'), '\"', '}')))]",
                                "recoveryProximityPlacementGroupCustomInput": "[if(empty(parameters('ppgId')), json('null'), json(concat('{', '\"resourceType\"', ':', '\"Existing\",', '\"recoveryProximityPlacementGroupId\"', ':', '\"', variables('recoveryPPGId'), '\"', '}')))]",
                                "recoveryVirtualNetworkCustomInput": "[if(contains(parameters('recoveryNetworkId'), '/'), json(concat('{', '\"resourceType\"', ':', '\"Existing\",', '\"recoveryVirtualNetworkId\"', ':', '\"', parameters('recoveryNetworkId'), '\"', '}')),if(empty(parameters('recoveryNetworkId')), json('null'), json(concat('{', '\"resourceType\"', ':', '\"New\",', '\"recoveryVirtualNetworkName\"', ':', '\"', parameters('recoveryNetworkId'), '\"', '}'))))]",
                                "primaryStagingStorageAccountCustomInput": "[if(empty(variables('cacheStorageAccountArmId')),json('null'),json(concat('{', '\"resourceType\"', ':', '\"Existing\",', '\"azureStorageAccountId\"', ':', '\"', variables('cacheStorageAccountArmId'), '\"', '}')))]",
                                "vmDisks": [],
                                "copy": [
                                  {
                                    "name": "vmManagedDisks",
                                    "count": "[variables('vmDiskCount')]",
                                    "input": {
                                      "diskId": "[if(equals(copyIndex('vmManagedDisks'), int(0)), reference(parameters('vmId'), variables('vmApiVersion')).storageProfile.osDisk.managedDisk.Id, variables('vmDataDiskIds')[sub(copyIndex('vmManagedDisks'), int(1))])]",
                                      "recoveryResourceGroupCustomInput": {
                                        "resourceType": "Existing",
                                        "recoveryResourceGroupId": "[variables('targetResourceGroupId')]"
                                      }
                                    }
                                  }
                                ]
                              }
                            }
                          }
                        ],
                        "outputs": {
                          "vmName": {
                            "value": "[last(split(parameters('vmId'), '/'))]",
                            "type": "string"
                          },
                          "availabilitySetUrl": {
                            "value": "[if(empty(parameters('avSetId')), '', concat(variables('portalLinkPrefix'), variables('recoveryAvSetId')))]",
                            "type": "string"
                          },
                          "proximityPlacementGroupUrl": {
                            "value": "[if(empty(parameters('ppgId')), '', concat(variables('portalLinkPrefix'), variables('recoveryPPGId')))]",
                            "type": "string"
                          },
                          "replicationEligibilityResults": {
                            "value": "[reference(variables('eligibilityResultsId'), parameters('apiVersion'))]",
                            "type": "Object"
                          }
                        }
                      },
                      "parameters": {}
                    },
                    "dependsOn": [
                      "[variables('ppgTemplateName64')]",
                      "[variables('avSetTemplateName64')]",
                      "[variables('storageAccountTemplateName64')]"
                    ]
                  }
                ],
                "outputs": {}
              },
              "parameters": {
                "apiVersion": {
                  "value": "2018-07-10"
                },
                "avSetId": {
                  "value": "[field('Microsoft.Compute/virtualMachines/availabilitySet.id')]"
                },
                "dataDiskIds": {
                  "value": {
                    "rawValue": "[field('Microsoft.Compute/virtualMachines/storageProfile.dataDisks[*].managedDisk.id')]",
                    "emptyArray": []
                  }
                },
                "osDiskId": {
                  "value": "[field('Microsoft.Compute/virtualMachines/storageProfile.osDisk.managedDisk.id')]"
                },
                "ppgId": {
                  "value": "[field('Microsoft.Compute/virtualMachines/proximityPlacementGroup.id')]"
                },
                "recoveryNetworkId": {
                  "value": "[parameters('recoveryNetworkId')]"
                },
                "recoverySubscriptionId": {
                  "value": "[subscription().subscriptionId]"
                },
                "sourceRegion": {
                  "value": "[parameters('sourceRegion')]"
                },
                "sourceResourceGroupName": {
                  "value": "[resourcegroup().Name]"
                },
                "targetRegion": {
                  "value": "[parameters('targetRegion')]"
                },
                "targetResourceGroupName": {
                  "value": "[last(split(parameters('targetResourceGroupId'), '/'))]"
                },
                "targetZone": {
                  "value": "[parameters('targetZone')]"
                },
                "vaultName": {
                  "value": "[last(split(parameters('vaultId'), '/'))]"
                },
                "vaultResourceGroupName": {
                  "value": "[last(split(parameters('vaultResourceGroupId'), '/'))]"
                },
                "vmId": {
                  "value": "[field('id')]"
                },
                "vmZones": {
                  "value": {
                    "rawValue": "[field('Microsoft.Compute/virtualMachines/zones')]",
                    "emptyArray": []
                  }
                },
                "cacheStorageAccountId": {
                  "value": "[parameters('cacheStorageAccountId')]"
                }
              }
            }
          }
        }
      }
    }
  },
  "id": "/providers/Microsoft.Authorization/policyDefinitions/ac34a73f-9fa5-4067-9247-a3ecae514468",
  "name": "ac34a73f-9fa5-4067-9247-a3ecae514468"
}`

func TestBasicTestAzurePolicyToRego(t *testing.T) {
	rulesJson, err := os.ReadFile("rules.json")
	require.NoError(t, err)
	outputJson, err := os.ReadFile("output.json")
	require.NoError(t, err)

	expectedDenyRego := `package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 aaaaa
}
aaaaa if {
type == "azurerm_service_plan"
aaaaa
}
aaaaa if {
not r.change.after.sku[0].tier in ["Basic","Standard","ElasticPremium","Premium","PremiumV2","Premium0V3","PremiumV3","PremiumMV3","Isolated","IsolatedV2","WorkflowStandard"]
not r.change.after.sku_name in ["B1","B2","B3","S1","S2","S3","EP1","EP2","EP3","P1","P2","P3","P1V2","P2V2","P3V2","P0V3","P1V3","P2V3","P3V3","P1MV3","P2MV3","P3MV3","P4MV3","P5MV3","I1","I2","I3","I1V2","I2V2","I3V2","I4V2","I5V2","I6V2","WS1","WS2","WS3"]
}`

	expectedCountRego := `package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 aaaaa
}
aaaaa if {
type == "azurerm_app_service_environment_v3"
regex.match("ASE*",kind)
count({x|r.change.after.cluster_setting[x];aaaaaaaaa(x)}) < 1
}
aaaaaaaaa(x) if {
aaaaa(x)
}
aaaaa(x) if {
r.change.after.cluster_setting[x].name == "DisableTls1.0"
r.change.after.cluster_setting[x].value == "1"
}`

	expectedNestedRego := `package main

import rego.v1

r := tfplan.resource_changes[_]

deny if {
 aaaaa
}
condition1 if {
 not aaaaaaa
 not aaaaaaa
}
aaaaaaa if {
 type != azurerm_network_security_group
 type != azurerm_network_security_rule
}
aaaaaaa if {
 count({x|r.change.after.security_rule[x];aaaaaaaaa(x)}) <= 0
 not aaaaaaa
}
condition4(x) if {
 condition6(x)
}
condition6(x) if {
 r.change.after.security_rule[x].access == Allow
 r.change.after.security_rule[x].direction == Inbound
 some r.change.after.security_rule[x].sourceAddressPrefix in ["*","Internet"]
 some r.change.after.security_rule[x].destinationPortRange in ["*","22","3389"]
}
condition5 if {
 type == azurerm_network_security_rule
 r.change.after.access == Allow
 r.change.after.direction == Inbound
 some r.change.after.source_address_prefix in ["*","Internet"]
 not condition7
}
condition7 if {
 not r.change.after.destination_port_range in ["*","22","3389"]
 condition8
}
condition8 if {
 not r.change.after.destination_port_ranges[0] in ["*","22","3389"]
}
`
	expectedNestedRego2 := ``

	cases := []struct {
		desc         string
		inputDirPath string
		mockFs       map[string]string
		expected     map[string]string
	}{
		{
			desc:         "deny.json",
			inputDirPath: "",
			mockFs: map[string]string{
				"deny.json": deny_json,
			},
			expected: map[string]string{
				"deny.rego": expectedDenyRego,
			},
		},
		{
			desc:         "dirPath",
			inputDirPath: "/config",
			mockFs: map[string]string{
				"/config/deny.json": deny_json,
			},
			expected: map[string]string{
				"deny.rego": expectedDenyRego,
			},
		},
		{
			desc:         "multiple json files in dirPath",
			inputDirPath: "/config",
			mockFs: map[string]string{
				"/config/deny1.json": deny_json,
				"/config/deny2.json": deny_json,
			},
			expected: map[string]string{
				"deny1.rego": expectedDenyRego,
				"deny2.rego": expectedDenyRego,
			},
		},
		{
			desc:         "json files in grandson's folders",
			inputDirPath: "/config",
			mockFs: map[string]string{
				"/config/deny1/deny1.json": deny_json,
				"/config/deny2/deny2.json": deny_json,
			},
			expected: map[string]string{
				"deny1.rego": expectedDenyRego,
				"deny2.rego": expectedDenyRego,
			},
		},
		{
			desc:         "policy contains count operator",
			inputDirPath: "",
			mockFs: map[string]string{
				"count.json": count_json,
			},
			expected: map[string]string{
				"count.rego": expectedCountRego,
			},
		},
		{
			desc:         "policy contains nested operations",
			inputDirPath: "",
			mockFs: map[string]string{
				"nested.json": nested_json,
			},
			expected: map[string]string{
				"nested.rego": expectedNestedRego,
			},
		},
		{
			desc:         "policy contains nested operations 2",
			inputDirPath: "",
			mockFs: map[string]string{
				"nested2.json": nested_json2,
			},
			expected: map[string]string{
				"nested2.rego": expectedNestedRego2,
			},
		},
	}

	//for i := 0; i < 10; i++ {
	for _, c := range cases {
		t.Run(fmt.Sprintf("%s", c.desc), func(t *testing.T) {
			files := map[string]string{
				"rules.json":  string(rulesJson),
				"output.json": string(outputJson),
			}
			for n, f := range c.mockFs {
				files[n] = f
			}
			mockFs := fakeFs(files)
			stub := gostub.Stub(&RandIntRange, func(min int, max int) int {
				return 0
			}).Stub(&Fs, mockFs)
			defer stub.Reset()
			policyPath := ""
			if len(c.mockFs) == 1 {
				for n, _ := range c.mockFs {
					policyPath = n
				}
			}
			require.NoError(t, AzurePolicyToRego(policyPath, c.inputDirPath, NewContext()))
			for n, expected := range c.expected {
				content, err := afero.ReadFile(mockFs, n)
				require.NoError(t, err)
				assert.Equal(t, expected, string(content))
			}
		})
	}
	//}
}

func fakeFs(files map[string]string) afero.Fs {
	fs := afero.NewMemMapFs()
	for n, content := range files {
		_ = afero.WriteFile(fs, n, []byte(content), 0644)
	}
	return fs
}

func TestMapEffectToAction(t *testing.T) {
	tests := []struct {
		name          string
		thenBody      *ThenBody
		defaultEffect string
		expected      string
		expectError   bool
	}{
		{
			name: "Effect is deny",
			thenBody: &ThenBody{
				Effect: "deny",
			},
			defaultEffect: "",
			expected:      "deny",
		},
		{
			name: "Effect is [parameters('effect')] and defaultEffect is deny",
			thenBody: &ThenBody{
				Effect: "[parameters('effect')]",
			},
			defaultEffect: "deny",
			expected:      "deny",
		},
		{
			name: "Effect is [parameters('effect')] and defaultEffect is audit",
			thenBody: &ThenBody{
				Effect: "[parameters('effect')]",
			},
			defaultEffect: "audit",
			expected:      "warn",
		},
		{
			name: "Effect is [parameters('effect')] and defaultEffect is disabled",
			thenBody: &ThenBody{
				Effect: "[parameters('effect')]",
			},
			defaultEffect: "disabled",
			expected:      "disabled",
		},
		{
			name: "Effect is empty and defaultEffect is deny",
			thenBody: &ThenBody{
				Effect: "",
			},
			defaultEffect: "deny",
			expected:      "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.thenBody.MapEffectToAction(tt.defaultEffect)
			if tt.expectError {
				assert.NotNil(t, err)
				return
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}
