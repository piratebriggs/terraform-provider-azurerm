package mssql

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmMsSqlServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmMsSqlServerRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateMsSqlServerName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"administrator_login": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"azuread_administrator": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"login_username": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"object_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"connection_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"identity": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"minimal_tls_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_network_access_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"fully_qualified_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func dataSourceArmMsSqlServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServersClient
	connectionClient := meta.(*clients.Client).MSSQL.ServerConnectionPoliciesClient
	adminClient := meta.(*clients.Client).MSSQL.ServerAzureADAdministratorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("SQL Server %q (Resource Group %q) was not found", name, resGroup)
		}

		return fmt.Errorf("reading SQL Server %q: %v", name, err)
	}

	if id := resp.ID; id != nil {
		d.SetId(*resp.ID)
	}
	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if err := d.Set("identity", flattenAzureRmSqlServerIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if props := resp.ServerProperties; props != nil {
		d.Set("version", props.Version)
		d.Set("administrator_login", props.AdministratorLogin)
		d.Set("fully_qualified_domain_name", props.FullyQualifiedDomainName)
		d.Set("minimal_tls_version", props.MinimalTLSVersion)
		d.Set("public_network_access_enabled", props.PublicNetworkAccess == sql.ServerPublicNetworkAccessEnabled)
	}

	adminResp, err := adminClient.Get(ctx, resGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(adminResp.Response) {
			return fmt.Errorf("reading SQL Server %s AAD admin: %v", name, err)
		}
	} else {
		if err := d.Set("azuread_administrator", flatternAzureRmMsSqlServerAdministrator(adminResp)); err != nil {
			return fmt.Errorf("setting `azuread_administrator`: %+v", err)
		}
	}

	connection, err := connectionClient.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("reading SQL Server %s Blob Connection Policy: %v ", name, err)
	}

	if props := connection.ServerConnectionPolicyProperties; props != nil {
		d.Set("connection_policy", string(props.ConnectionType))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
