package manifold

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/manifoldco/kubernetes-credentials/helpers/client"
)

func dataSourceManifoldProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceManifoldProjectRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The project label of which you want to retrieve the data.",
			},

			"resource": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A specific resource you want to load for this project.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The resource label you want to get the resource for",
						},

						"credential": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The name for this credential which will be used to reference later on. Defaults to the `key` value.",
									},

									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The key to fetch from the resource credentials.",
									},

									"default": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The default value for this key if it's not set.",
									},
								},
							},
						},
					},
				},
			},

			"credentials": {
				Type:        schema.TypeMap,
				Computed:    true,
				Sensitive:   true,
				Description: "The credentials linked to this project.",
			},
		},
	}
}

func dataSourceManifoldProjectRead(d *schema.ResourceData, meta interface{}) error {
	cl := meta.(*client.Client)
	ctx := context.Background()

	projectLabel, projectID, err := getProjectInformation(cl, d.Get("project").(string), true)
	if err != nil {
		return err
	}

	resourceList := d.Get("resource").(*schema.Set).List()
	filteredResources := resourceSpecsFromList(resourceList)

	cv, err := cl.GetResourcesCredentialValues(ctx, projectLabel, filteredResources)
	if err != nil {
		return err
	}

	credMap, err := client.FlattenResourcesCredentialValues(cv)
	if err != nil {
		return err
	}

	d.SetId(projectID.String())
	d.Set("project", *projectLabel)
	d.Set("credentials", credMap)
	return nil
}
