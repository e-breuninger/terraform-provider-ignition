package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	butane "github.com/coreos/butane/config"
	"github.com/coreos/butane/config/common"
	"github.com/coreos/go-semver/semver"
	ignition_util "github.com/coreos/ignition/v2/config/util"
	ignition_v3_0 "github.com/coreos/ignition/v2/config/v3_0"
	types_v3_0 "github.com/coreos/ignition/v2/config/v3_0/types"
	ignition_v3_1 "github.com/coreos/ignition/v2/config/v3_1"
	types_v3_1 "github.com/coreos/ignition/v2/config/v3_1/types"
	ignition_v3_2 "github.com/coreos/ignition/v2/config/v3_2"
	types_v3_2 "github.com/coreos/ignition/v2/config/v3_2/types"
	ignition_v3_3 "github.com/coreos/ignition/v2/config/v3_3"
	types_v3_3 "github.com/coreos/ignition/v2/config/v3_3/types"
	ignition_v3_4 "github.com/coreos/ignition/v2/config/v3_4"
	types_v3_4 "github.com/coreos/ignition/v2/config/v3_4/types"
	"github.com/coreos/vcontext/report"
)

func parseV30(rawConfig []byte) (interface{}, report.Report, error) {
	return ignition_v3_0.ParseCompatibleVersion(rawConfig)
}

func mergeV30(parent interface{}, child interface{}) interface{} {
	return ignition_v3_0.Merge(parent.(types_v3_0.Config), child.(types_v3_0.Config))
}

func parseV31(rawConfig []byte) (interface{}, report.Report, error) {
	return ignition_v3_1.ParseCompatibleVersion(rawConfig)
}

func mergeV31(parent interface{}, child interface{}) interface{} {
	return ignition_v3_1.Merge(parent.(types_v3_1.Config), child.(types_v3_1.Config))
}

func parseV32(rawConfig []byte) (interface{}, report.Report, error) {
	return ignition_v3_2.ParseCompatibleVersion(rawConfig)
}

func mergeV32(parent interface{}, child interface{}) interface{} {
	return ignition_v3_2.Merge(parent.(types_v3_2.Config), child.(types_v3_2.Config))
}

func parseV33(rawConfig []byte) (interface{}, report.Report, error) {
	return ignition_v3_3.ParseCompatibleVersion(rawConfig)
}

func mergeV33(parent interface{}, child interface{}) interface{} {
	return ignition_v3_3.Merge(parent.(types_v3_3.Config), child.(types_v3_3.Config))
}

func parseV34(rawConfig []byte) (interface{}, report.Report, error) {
	return ignition_v3_4.ParseCompatibleVersion(rawConfig)
}

func mergeV34(parent, child interface{}) interface{} {
	return ignition_v3_4.Merge(parent.(types_v3_4.Config), child.(types_v3_4.Config))
}

type ignitionInterface struct {
	Parse func(rawConfig []byte) (interface{}, report.Report, error)
	Merge func(parent interface{}, child interface{}) interface{}
}

var versionToLibrary = map[string]ignitionInterface{
	"3.0.0": {
		Parse: parseV30,
		Merge: mergeV30,
	},
	"3.1.0": {
		Parse: parseV31,
		Merge: mergeV31,
	},
	"3.2.0": {
		Parse: parseV32,
		Merge: mergeV32,
	},
	"3.3.0": {
		Parse: parseV33,
		Merge: mergeV33,
	},
	"3.4.0": {
		Parse: parseV34,
		Merge: mergeV34,
	},
}

func getLibraryForVersion(version string) (ignitionInterface, error) {
	ignition, ok := versionToLibrary[version]
	if !ok {
		return ignitionInterface{}, errors.New("incompatible version")
	}
	return ignition, nil
}

func datasourceConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceConfigRead,

		Schema: map[string]*schema.Schema{
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"snippets": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},
			"pretty_print": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"strict": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"rendered": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "rendered ignition configuration",
			},
		},
	}
}

func datasourceConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	renderedBytes, err := renderConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}
	rendered := string(renderedBytes)

	if err := d.Set("rendered", rendered); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(hashcode(rendered))
	return diags
}

type getConfigVersion func(ignition []byte) (semver.Version, error)

func renderConfig(d *schema.ResourceData) ([]byte, error) {
	// unchecked assertions seem to be the norm in Terraform :S
	content := d.Get("content").(string)
	pretty := d.Get("pretty_print").(bool)
	strict := d.Get("strict").(bool)
	snippetsIface := d.Get("snippets").([]interface{})

	snippets := make([]string, len(snippetsIface))
	for i, v := range snippetsIface {
		if v == nil {
			continue
		}
		snippets[i] = v.(string)
	}

	// transpile content
	ignitionConfig, contentVersion, ignition, err := transpileButane(
		content,
		strict,
		func(ignitionBytes []byte) (semver.Version, error) {
			version, _, err := ignition_util.GetConfigVersion(ignitionBytes)
			return version, err
		},
	)
	if err != nil {
		return nil, fmt.Errorf("content parse error: %v", err)
	}

	// transpile snippets and merge them with content
	for _, snippet := range snippets {
		snippetIgnitionConfig, _, _, err := transpileButane(snippet, strict, ensureMaxVersion(contentVersion))
		if err != nil {
			return nil, fmt.Errorf("snippet parse error: %v", err)
		}
		ignitionConfig = ignition.Merge(ignitionConfig, snippetIgnitionConfig)
	}

	// marshal json
	if pretty {
		return json.MarshalIndent(ignitionConfig, "", "  ")
	}
	return json.Marshal(ignitionConfig)
}

// Transpile Butane into a Ignition configuration object determined by the Ignitition version given.
// Returns the Ignition configuration object, the Ignition version used and the matching Ignition library
func transpileButane(
	butaneConfig string,
	strict bool,
	getIgnitionVersion getConfigVersion,
) (interface{}, semver.Version, ignitionInterface, error) {
	ignitionBytes, report, err := butane.TranslateBytes([]byte(butaneConfig), common.TranslateBytesOptions{})
	if err != nil {
		return nil, semver.Version{}, ignitionInterface{}, err
	}
	if strict && len(report.Entries) > 0 {
		return nil, semver.Version{}, ignitionInterface{}, fmt.Errorf("strict parsing error: %v", report.String())
	}
	version, err := getIgnitionVersion(ignitionBytes)
	if err != nil {
		return nil, semver.Version{}, ignitionInterface{}, err
	}
	ignition, err := getLibraryForVersion(version.String())
	if err != nil {
		return nil, semver.Version{}, ignitionInterface{}, err
	}
	ignitionConfig, _, err := ignition.Parse(ignitionBytes)

	return ignitionConfig, version, ignition, err
}

// prepare function to validate snippets against
func ensureMaxVersion(maxVersion semver.Version) getConfigVersion {
	return func(ignitionBytes []byte) (semver.Version, error) {
		version, _, err := ignition_util.GetConfigVersion(ignitionBytes)
		if err != nil {
			return semver.Version{}, err
		}
		if maxVersion.LessThan(version) {
			return semver.Version{}, fmt.Errorf(
				"version %s is newer than max version %s and therefore incompatible",
				version.String(),
				maxVersion.String(),
			)
		}
		return maxVersion, nil
	}
}
