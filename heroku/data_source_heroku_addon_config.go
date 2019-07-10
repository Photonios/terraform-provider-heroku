package heroku

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	heroku "github.com/heroku/heroku-go/v3"
)

const (
	HEROKU_ADDON_STATE_PROVISIONING = "provisioning"
	HEROKU_ADDON_STATE_PROVISIONED  = "provisioned"
)

func dataSourceHerokuAddonConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHerokuAddonConfigRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"config": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceHerokuAddonConfigRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Config).Api

	name := d.Get("name").(string)

	addon, err := resourceHerokuAddonRetrieve(name, client)
	if err != nil {
		return err
	}

	if err := dataSourceHerokuAddonWaitForProvisioned(addon, client); err != nil {
		return err
	}

	config, err := dataSourceHerokuAddonConfigRetrieve(name, client)
	if err != nil {
		return err
	}

	d.SetId(addon.ID)
	d.Set("config", config)

	return nil
}

// Waits for the AddOn to be in the "provisioned" state if it is
// not already. This is needed because AddOn configuration is not
// available until the AddOn is fully provisioned.
func dataSourceHerokuAddonWaitForProvisioned(addon *heroku.AddOn, client *heroku.Service) error {
	if addon.State == HEROKU_ADDON_STATE_PROVISIONED {
		log.Printf("[DEBUG] Addon (%s) is provisioned", addon.Name)
		return nil
	}

	refreshFunc := func() (interface{}, string, error) {
		refreshedAddon, err := resourceHerokuAddonRetrieve(addon.ID, client)
		if err != nil {
			return nil, "", err
		}

		return (*heroku.AddOn)(refreshedAddon), refreshedAddon.State, nil
	}

	log.Printf("[DEBUG] Waiting for Addon (%s) to be provisioned", addon.Name)
	stateConf := &resource.StateChangeConf{
		Pending: []string{HEROKU_ADDON_STATE_PROVISIONING},
		Target:  []string{HEROKU_ADDON_STATE_PROVISIONED},
		Refresh: refreshFunc,
		Timeout: 20 * time.Minute,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Addon (%s) to be provisioned: %s",
			addon.Name, err)
	}

	return nil
}

func dataSourceHerokuAddonConfigRetrieve(name string, client *heroku.Service) (map[string]string, error) {
	log.Printf("[DEBUG] Retrieving AddOn (%s) config", name)

	listRange := &heroku.ListRange{Descending: true}
	addonConfig, err := client.AddOnConfigList(context.TODO(), name, listRange)

	if err != nil {
		return nil, fmt.Errorf("Error retrieving AddOn (%s) config: %s",
			name, err)
	}

	config := make(map[string]string)
	for _, configVar := range addonConfig {
		config[configVar.Name] = *configVar.Value
	}

	return config, nil
}
