package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/keycloak/terraform-provider-keycloak/keycloak"
	"regexp"
	"testing"
)

// Tests for attaching SAML mappers to SAML client scopes are omitted
// because the keycloak_saml_client_scope resource does not exist yet.

func TestAccKeycloakSamlUserAttributeProtocolMapper_basicClient(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	resourceName := "keycloak_saml_user_attribute_protocol_mapper.saml_user_attribute_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakSamlUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlUserAttributeProtocolMapper_basic_client(clientId, mapperName),
				Check:  testKeycloakSamlUserAttributeProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakSamlUserAttributeProtocolMapper_import(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	clientResourceName := "keycloak_saml_user_attribute_protocol_mapper.saml_user_attribute_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakSamlUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlUserAttributeProtocolMapper_basic_client(clientId, mapperName),
				Check:  testKeycloakSamlUserAttributeProtocolMapperExists(clientResourceName),
			},
			{
				ResourceName:      clientResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getGenericProtocolMapperIdForClient(clientResourceName),
			},
		},
	})
}

func TestAccKeycloakSamlUserAttributeProtocolMapper_update(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	userAttribute := acctest.RandomWithPrefix("tf-acc")
	updatedUserAttribute := acctest.RandomWithPrefix("tf-acc")
	resourceName := "keycloak_saml_user_attribute_protocol_mapper.saml_user_attribute_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakSamlUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlUserAttributeProtocolMapper_userAttribute(clientId, mapperName, userAttribute),
				Check:  testKeycloakSamlUserAttributeProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakSamlUserAttributeProtocolMapper_userAttribute(clientId, mapperName, updatedUserAttribute),
				Check:  testKeycloakSamlUserAttributeProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakSamlUserAttributeProtocolMapper_createAfterManualDestroy(t *testing.T) {
	t.Parallel()
	var mapper = &keycloak.SamlUserAttributeProtocolMapper{}

	clientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	resourceName := "keycloak_saml_user_attribute_protocol_mapper.saml_user_attribute_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakSamlUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlUserAttributeProtocolMapper_basic_client(clientId, mapperName),
				Check:  testKeycloakSamlUserAttributeProtocolMapperFetch(resourceName, mapper),
			},
			{
				PreConfig: func() {
					err := keycloakClient.DeleteSamlUserAttributeProtocolMapper(testCtx, mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)
					if err != nil {
						t.Error(err)
					}
				},
				Config: testKeycloakSamlUserAttributeProtocolMapper_basic_client(clientId, mapperName),
				Check:  testKeycloakSamlUserAttributeProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakSamlUserAttributeProtocolMapper_validateClaimValueType(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")
	invalidSamlNameFormat := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakSamlUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlUserAttributeProtocolMapper_attributes(clientId,
					mapperName,
					"foo",
					"bar",
					invalidSamlNameFormat,
					false),
				ExpectError: regexp.MustCompile("expected saml_attribute_name_format to be one of .+ got " + invalidSamlNameFormat),
			},
		},
	})
}

func TestAccKeycloakSamlUserAttributeProtocolMapper_updateClientIdForceNew(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	updatedClientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	userAttribute := acctest.RandomWithPrefix("tf-acc")
	resourceName := "keycloak_saml_user_attribute_protocol_mapper.saml_user_attribute_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakSamlUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlUserAttributeProtocolMapper_userAttribute(clientId, mapperName, userAttribute),
				Check:  testKeycloakSamlUserAttributeProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakSamlUserAttributeProtocolMapper_userAttribute(updatedClientId, mapperName, userAttribute),
				Check:  testKeycloakSamlUserAttributeProtocolMapperExists(resourceName),
			},
		},
	})
}
func TestAccKeycloakSamlUserAttributeProtocolMapperHasChangingAttributes(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	resourceName := "keycloak_saml_user_attribute_protocol_mapper.saml_user_attribute_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakSamlUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlUserAttributeProtocolMapper_attributes(clientId,
					mapperName,
					"foo",
					"bar",
					"Basic",
					false),
				Check: testKeycloakSamlUserAttributeProtocolMapperHasGivenAttributes(resourceName,
					mapperName,
					"foo",
					"bar",
					"Basic",
					false),
			},
			{
				Config: testKeycloakSamlUserAttributeProtocolMapper_attributes(clientId,
					mapperName,
					"email",
					"email",
					"Unspecified",
					true),
				Check: testKeycloakSamlUserAttributeProtocolMapperHasGivenAttributes(resourceName,
					mapperName,
					"email",
					"email",
					"Unspecified",
					true),
			},
		},
	})
}

func testAccKeycloakSamlUserAttributeProtocolMapperDestroy() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for resourceName, rs := range state.RootModule().Resources {
			if rs.Type != "keycloak_saml_user_attribute_protocol_mapper" {
				continue
			}

			mapper, _ := getSamlUserAttributeMapperUsingState(state, resourceName)

			if mapper != nil {
				return fmt.Errorf("saml user attribute protocol mapper with id %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testKeycloakSamlUserAttributeProtocolMapperExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, err := getSamlUserAttributeMapperUsingState(state, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}
func testKeycloakSamlUserAttributeProtocolMapperHasGivenAttributes(resourceName string, mapperName string, userAttribute string, samlAttributeName string, samlAttributeNameFormat string, aggregateAttributes bool) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		mapper, err := getSamlUserAttributeMapperUsingState(state, resourceName)
		if err != nil {
			return err
		}
		if mapper.Name != mapperName {
			return fmt.Errorf("mapper name is %s, expected %s", mapper.Name, mapperName)
		} else if mapper.UserAttribute != userAttribute {
			return fmt.Errorf("userAttribute is %s, expected %s", mapper.UserAttribute, userAttribute)
		} else if mapper.SamlAttributeName != samlAttributeName {
			return fmt.Errorf("samlAttribute name is %s, expected %s", mapper.SamlAttributeName, samlAttributeName)
		} else if mapper.SamlAttributeNameFormat != samlAttributeNameFormat {
			return fmt.Errorf("samlAttributeName format is %s, expected %s", mapper.SamlAttributeNameFormat, samlAttributeNameFormat)
		} else if mapper.AggregateAttributeValues != aggregateAttributes {
			return fmt.Errorf("aggregateAttributes is %t, expected %t", mapper.AggregateAttributeValues, aggregateAttributes)
		}
		return nil
	}
}

func testKeycloakSamlUserAttributeProtocolMapperFetch(resourceName string, mapper *keycloak.SamlUserAttributeProtocolMapper) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		fetchedMapper, err := getSamlUserAttributeMapperUsingState(state, resourceName)
		if err != nil {
			return err
		}

		mapper.Id = fetchedMapper.Id
		mapper.ClientId = fetchedMapper.ClientId
		mapper.ClientScopeId = fetchedMapper.ClientScopeId
		mapper.RealmId = fetchedMapper.RealmId

		return nil
	}
}

func getSamlUserAttributeMapperUsingState(state *terraform.State, resourceName string) (*keycloak.SamlUserAttributeProtocolMapper, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found in TF state: %s ", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]
	clientId := rs.Primary.Attributes["client_id"]
	clientScopeId := rs.Primary.Attributes["client_scope_id"]

	return keycloakClient.GetSamlUserAttributeProtocolMapper(testCtx, realm, clientId, clientScopeId, id)
}

func testKeycloakSamlUserAttributeProtocolMapper_basic_client(clientId, mapperName string) string {
	return testKeycloakSamlUserAttributeProtocolMapper_attributes(clientId, mapperName, "foo", "bar", "Unspecified", false)
}

func testKeycloakSamlUserAttributeProtocolMapper_userAttribute(clientId, mapperName, userAttribute string) string {
	return testKeycloakSamlUserAttributeProtocolMapper_attributes(clientId, mapperName, userAttribute, "bar", "Unspecified", false)
}
func testKeycloakSamlUserAttributeProtocolMapper_attributes(clientId, mapperName string, userAttribute string, samlAttributeName string, samlAttributeNameFormat string, aggregateAttributes bool) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
	realm_id  = data.keycloak_realm.realm.id
	client_id = "%s"
}

resource "keycloak_saml_user_attribute_protocol_mapper" "saml_user_attribute_mapper" {
	name                       = "%s"
	realm_id                   = data.keycloak_realm.realm.id
	client_id                  = "${keycloak_saml_client.saml_client.id}"

	user_attribute             = "%s"
	saml_attribute_name        = "%s"
	saml_attribute_name_format = "%s"
	aggregate_attributes	   = %t
}`, testAccRealm.Realm, clientId,
		mapperName,
		userAttribute,
		samlAttributeName,
		samlAttributeNameFormat,
		aggregateAttributes)
}
