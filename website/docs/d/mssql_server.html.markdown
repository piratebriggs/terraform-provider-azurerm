---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_server"
description: |-
  Gets information about an existing Azure SQL Database Server.
---

# Data Source: azurerm_mssql_server

Use this data source to access information about an existing Azure SQL Database Server.

## Example Usage

```hcl
data "azurerm_mssql_server" "example" {
  name                = "mssql-server"
  resource_group_name = "database-rg"
}
```

## Argument Reference

* `name` - Specifies the name of the SQL Server.

* `resource_group_name` - Specifies the name of the Resource Group where the SQL Server exists.

## Attribute Reference

* `location` - The supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `version` - The version of the server. Possible values are: `2.0` for v11 server and `12.0` for v12 server.

* `administrator_login` - The administrator login name for the server.

* `azuread_administrator` - An `azuread_administrator` block as defined below.

* `connection_policy` - The connection policy of the server.

* `identity` - An `identity` block as defined below.

* `minimal_tls_version` - The minimal TLS version for all SQL Database and SQL Data Warehouse databases associated with the server.

* `public_network_access_enabled` - Whether or not public network access is allowed for this server.

* `tags` - A mapping of tags assigned to the resource.

---

`identity` exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this SQL Server.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this SQL Server.

-> You can access the Principal ID via `${data.azurerm_mssql_server.example.identity.0.principal_id}` and the Tenant ID via `${data.azurerm_mssql_server.example.identity.0.tenant_id}`

---

`azuread_administrator` exports the following:

* `login_username` - The login username of the Azure AD Administrator of this SQL Server.

* `object_id` - The object id of the Azure AD Administrator of this SQL Server.

* `tenant_id` - The tenant id of the Azure AD Administrator of this SQL Server.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the SQL database.
