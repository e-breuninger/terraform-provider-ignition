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
	"3.0.0": ignitionInterface{
		Parse: parseV30,
		Merge: mergeV30,
	},
	"3.1.0": ignitionInterface{
		Parse: parseV31,
		Merge: mergeV31,
	},
	"3.2.0": ignitionInterface{
		Parse: parseV32,
		Merge: mergeV32,
	},
	"3.3.0": ignitionInterface{
		Parse: parseV33,
		Merge: mergeV33,
	},
	"3.4.0": ignitionInterface{
		Parse: parseV34,
		Merge: mergeV34,
	},
}

func getLibraryForVersion(version string) (ignitionInterface, error) {
	ignition, ok := versionToLibrary[version]
	if !ok {
		return ignitionInterface{}, errors.New("Incompatible version.")
	}
	return ignition, nil
}

func DatasourceConfig() *schema.Resource {
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

	rendered, err := renderConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("rendered", rendered); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(hashcode(rendered))
	return diags
}

// Render a Fedora CoreOS Config or Container Linux Config as Ignition JSON.
func renderConfig(d *schema.ResourceData) (string, error) {
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

	// Butane Config
	ign, err := butaneToIgnition([]byte(content), pretty, strict, snippets)
	return string(ign), err
}

// Translate Fedora CoreOS config to Ignition v3.X.Y
func butaneToIgnition(data []byte, pretty, strict bool, snippets []string) ([]byte, error) {
	ignBytes, report, err := butane.TranslateBytes(data, common.TranslateBytesOptions{
		Pretty: pretty,
	})
	// ErrNoVariant indicates data is a CLC, not an FCC
	if err != nil {
		return nil, err
	}
	if strict && len(report.Entries) > 0 {
		return nil, fmt.Errorf("strict parsing error: %v", report.String())
	}

	// merge FCC snippets into main Ignition config
	return mergeFCCSnippets(ignBytes, pretty, strict, snippets)
}

// Parse Fedora CoreOS Ignition and Butane snippets into Ignition Config.
func mergeFCCSnippets(ignBytes []byte, pretty, strict bool, snippets []string) ([]byte, error) {
	semver, _, _ := ignition_util.GetConfigVersion(ignBytes)
	ignition, err := getLibraryForVersion(semver.String())

	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	ign, _, err := ignition.Parse(ignBytes)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	for _, snippet := range snippets {
		ignextBytes, report, err := butane.TranslateBytes([]byte(snippet), common.TranslateBytesOptions{
			Pretty: pretty,
		})
		if err != nil {
			// For FCC, require snippets be FCCs (don't fall-through to CLC)
			if err == common.ErrNoVariant {
				return nil, fmt.Errorf("Butane snippets require `variant`: %v", err)
			}
			return nil, fmt.Errorf("Butane translate error: %v", err)
		}
		if strict && len(report.Entries) > 0 {
			return nil, fmt.Errorf("strict parsing error: %v", report.String())
		}
		snippetSemver, _, _ := ignition_util.GetConfigVersion(ignextBytes)
		versionIsIncompatible := semver.LessThan(snippetSemver)
		if versionIsIncompatible {
			return nil, fmt.Errorf("Snippet version %s is newer than content version %s and therefore incompatible", snippetSemver.String(), semver.String())
		}

		ignext, _, err := ignition.Parse(ignextBytes)
		if err != nil {
			return nil, fmt.Errorf("snippet parse error: %v", err)
		}
		ign = ignition.Merge(ign, ignext)
	}

	return marshalJSON(ign, pretty)
}

func marshalJSON(v interface{}, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(v, "", "  ")
	}
	return json.Marshal(v)
}
