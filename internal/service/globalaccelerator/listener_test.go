package globalaccelerator_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/globalaccelerator"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfglobalaccelerator "github.com/hashicorp/terraform-provider-aws/internal/service/globalaccelerator"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccGlobalAcceleratorListener_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_globalaccelerator_listener.example"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, globalaccelerator.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccListenerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "client_affinity", "NONE"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "port_range.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "port_range.*", map[string]string{
						"from_port": "80",
						"to_port":   "81",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGlobalAcceleratorListener_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_globalaccelerator_listener.example"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, globalaccelerator.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccListenerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerExists(ctx, resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfglobalaccelerator.ResourceListener(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccGlobalAcceleratorListener_update(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_globalaccelerator_listener.example"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, globalaccelerator.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckListenerDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccListenerConfig_basic(rName),
			},
			{
				Config: testAccListenerConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "client_affinity", "SOURCE_IP"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "UDP"),
					resource.TestCheckResourceAttr(resourceName, "port_range.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "port_range.*", map[string]string{
						"from_port": "443",
						"to_port":   "444",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckListenerExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).GlobalAcceleratorConn()

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Global Accelerator Listener ID is set")
		}

		_, err := tfglobalaccelerator.FindListenerByARN(ctx, conn, rs.Primary.ID)

		return err
	}
}

func testAccCheckListenerDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).GlobalAcceleratorConn()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_globalaccelerator_listener" {
				continue
			}

			_, err := tfglobalaccelerator.FindListenerByARN(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("Global Accelerator Listener %s still exists", rs.Primary.ID)
		}
		return nil
	}
}

func testAccListenerConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_globalaccelerator_accelerator" "example" {
  name            = %[1]q
  ip_address_type = "IPV4"
  enabled         = false
}

resource "aws_globalaccelerator_listener" "example" {
  accelerator_arn = aws_globalaccelerator_accelerator.example.id
  protocol        = "TCP"

  port_range {
    from_port = 80
    to_port   = 81
  }
}
`, rName)
}

func testAccListenerConfig_update(rName string) string {
	return fmt.Sprintf(`
resource "aws_globalaccelerator_accelerator" "example" {
  name            = %[1]q
  ip_address_type = "IPV4"
  enabled         = false
}

resource "aws_globalaccelerator_listener" "example" {
  accelerator_arn = aws_globalaccelerator_accelerator.example.id
  client_affinity = "SOURCE_IP"
  protocol        = "UDP"

  port_range {
    from_port = 443
    to_port   = 444
  }
}
`, rName)
}
