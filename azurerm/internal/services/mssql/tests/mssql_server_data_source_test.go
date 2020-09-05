package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMMsSqlServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mssql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMsSqlServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMMsSqlServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mssql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMsSqlServer_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "missadministrator"),
					resource.TestCheckResourceAttr(data.ResourceName, "minimal_tls_version", "1.2"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_network_access_enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.ENV", "Staging"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "12.0"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMMsSqlServer_basic(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlServer_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_mssql_server" "test" {
  name                = azurerm_mssql_server.test.name
  resource_group_name = azurerm_mssql_server.test.resource_group_name
}

`, template)
}

func testAccDataSourceAzureRMMsSqlServer_complete(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlServer_complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_mssql_server" "test" {
  name                = azurerm_mssql_server.test.name
  resource_group_name = azurerm_mssql_server.test.resource_group_name
}

`, template)
}
