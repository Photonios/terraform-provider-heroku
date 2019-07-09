package heroku

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceHerokuAddon() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHerokuAddonRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"app": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"plan": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"provider_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"config_vars": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceHerokuAddonRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Config).Api

	name := d.Get("name").(string)

	addon, addonConfig, err := resourceHerokuAddonRetrieve(name, client)
	if err != nil {
		return err
	}

	configVars := make(map[string]string)
	for _, configVar := range addonConfig {
		configVars[configVar.Name] = *configVar.Value
	}

	d.SetId(addon.ID)
	d.Set("name", addon.Name)
	d.Set("app", addon.App.Name)
	d.Set("plan", addon.Plan.Name)
	d.Set("provider_id", addon.ProviderID)
	d.Set("config_vars", configVars)

	return nil
}
