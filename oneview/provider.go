// (C) Copyright 2016 Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package oneview

import (
	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var (
	ovMutexKV          = mutexkv.NewMutexKV()
	serverHardwareURIs = make(map[string]bool)
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"ov_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEVIEW_OV_DOMAIN", ""),
			},
			"ov_username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEVIEW_OV_USER", ""),
			},
			"ov_password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEVIEW_OV_PASSWORD", nil),
			},
			"ov_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEVIEW_OV_ENDPOINT", nil),
			},
			"ov_sslverify": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEVIEW_OV_SSLVERIFY", true),
			},
			"ov_apiversion": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEVIEW_OV_API_VERSION", 200),
			},
			"ov_ifmatch": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEVIEW_OV_IF_MATCH", "*"),
			},
			"icsp_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEVIEW_ICSP_DOMAIN", ""),
			},
			"icsp_username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEVIEW_ICSP_USER", ""),
			},
			"icsp_password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEVIEW_ICSP_PASSWORD", ""),
			},
			"icsp_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEVIEW_ICSP_ENDPOINT", ""),
			},
			"icsp_sslverify": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEVIEW_ICSP_SSLVERIFY", true),
			},
			"icsp_apiversion": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEVIEW_ICSP_API_VERSION", 200),
			},
			"i3s_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEVIEW_I3S_ENDPOINT", ""),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"oneview_fc_network":              dataSourceFCNetwork(),
			"oneview_enclosure":               dataSourceEnclosure(),
			"oneview_ethernet_network":        dataSourceEthernetNetwork(),
			"oneview_interconnect_type":       dataSourceInterconnectType(),
			"oneview_interconnect":            dataSourceInterconnects(),
			"oneview_logical_interconnect":    dataSourceLogicalInterconnect(),
			"oneview_scope":                   dataSourceScope(),
			"oneview_server_hardware":         dataSourceServerHardware(),
			"oneview_server_hardware_type":    dataSourceServerHardwareType(),
			"oneview_logical_enclosure":       dataSourceLogicalEnclosure(),
			"oneview_enclosure_group":         dataSourceEnclosureGroup(),
			"oneview_server_profile":          dataSourceServerProfile(),
			"oneview_server_profile_template": dataSourceServerProfileTemplate(),
			"oneview_storage_system":          dataSourceStorageSystem(),
			"oneview_uplink_set":              dataSourceUplinkSet(),
			"oneview_network_set":             dataSourceNetworkSet(),
			"oneview_storage_volume":          dataSourceStorageVolume(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"oneview_server_profile":             resourceServerProfile(),
			"oneview_enclosure":                  resourceEnclosure(),
			"oneview_enclosure_group":            resourceEnclosureGroup(),
			"oneview_ethernet_network":           resourceEthernetNetwork(),
			"oneview_network_set":                resourceNetworkSet(),
			"oneview_fcoe_network":               resourceFCoENetwork(),
			"oneview_fc_network":                 resourceFCNetwork(),
			"oneview_scope":                      resourceScope(),
			"oneview_server_profile_template":    resourceServerProfileTemplate(),
			"oneview_storage_system":             resourceStorageSystem(),
			"oneview_logical_interconnect_group": resourceLogicalInterconnectGroup(),
			"oneview_logical_interconnect":       resourceLogicalInterconnect(),
			"oneview_logical_switch_group":       resourceLogicalSwitchGroup(),
			"oneview_uplink_set":                 resourceUplinkSet(),
			"oneview_i3s_plan":                   resourceI3SPlan(),
			"oneview_logical_enclosure":          resourceLogicalEnclosure(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		OVDomain:     d.Get("ov_domain").(string),
		OVUsername:   d.Get("ov_username").(string),
		OVPassword:   d.Get("ov_password").(string),
		OVEndpoint:   d.Get("ov_endpoint").(string),
		OVSSLVerify:  d.Get("ov_sslverify").(bool),
		OVAPIVersion: d.Get("ov_apiversion").(int),
		OVIfMatch:    d.Get("ov_ifmatch").(string),
	}

	if err := config.loadAndValidate(); err != nil {
		return nil, err
	}

	if val, ok := d.GetOk("i3s_endpoint"); ok {
		config.I3SEndpoint = val.(string)
		if err := config.loadAndValidateI3S(); err != nil {
			return nil, err
		}
	}

	return &config, nil
}
