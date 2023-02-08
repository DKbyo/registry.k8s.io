/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"testing"
)

func TestGenerateRangesGo(t *testing.T) {
	// raw data to generate from
	const rawData = `{
  "syncToken": "1649878400",
  "createDate": "2022-04-13-19-33-20",
  "prefixes": [
    {
      "ip_prefix": "3.5.140.0/22",
      "region": "ap-northeast-2",
      "service": "AMAZON",
      "network_border_group": "ap-northeast-2"
    },
    {
      "ip_prefix": "52.95.174.0/24",
      "region": "me-south-1",
      "service": "AMAZON",
      "network_border_group": "me-south-1"
    },
    {
      "ip_prefix": "15.185.0.0/16",
      "region": "me-south-1",
      "service": "AMAZON",
      "network_border_group": "me-south-1"
    },
    {
      "ip_prefix": "69.107.7.136/29",
      "region": "me-south-1",
      "service": "AMAZON",
      "network_border_group": "me-south-1"
    }
  ],
  "ipv6_prefixes": [
    {
      "ipv6_prefix": "2a05:d07a:a000::/40",
      "region": "eu-south-1",
      "service": "AMAZON",
      "network_border_group": "eu-south-1"
    },
    {
      "ipv6_prefix": "2a05:d03a:a000:200::/56",
      "region": "eu-south-1",
      "service": "AMAZON",
      "network_border_group": "eu-south-1"
    },
    {
      "ipv6_prefix": "2a05:d03a:a000:400::/56",
      "region": "eu-south-1",
      "service": "AMAZON",
      "network_border_group": "eu-south-1"
    },
    {
      "ipv6_prefix": "2a05:d03a:a000::/56",
      "region": "eu-south-1",
      "service": "AMAZON",
      "network_border_group": "eu-south-1"
    }
  ]
}
`
	rtp, err := parseAWS(rawData)
	if err != nil {
		t.Fatalf("unexpected error parsing test data: %v", err)
	}

	// expected generated result
	const goldenText = `/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// File generated by ranges2go DO NOT EDIT

package cloudcidrs

import (
	"net/netip"
)

// AWS cloud
const AWS = "AWS"

// regionToRanges contains a preparsed map of cloud IPInfo to netip.Prefix
var regionToRanges = map[IPInfo][]netip.Prefix{
	{Cloud: AWS, Region: "ap-northeast-2"}: {
		netip.PrefixFrom(netip.AddrFrom4([4]byte{3, 5, 140, 0}), 22),
	},
	{Cloud: AWS, Region: "eu-south-1"}: {
		netip.PrefixFrom(netip.AddrFrom16([16]byte{42, 5, 208, 58, 160, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0}), 56),
		netip.PrefixFrom(netip.AddrFrom16([16]byte{42, 5, 208, 58, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), 56),
		netip.PrefixFrom(netip.AddrFrom16([16]byte{42, 5, 208, 122, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), 40),
	},
	{Cloud: AWS, Region: "me-south-1"}: {
		netip.PrefixFrom(netip.AddrFrom4([4]byte{52, 95, 174, 0}), 24),
		netip.PrefixFrom(netip.AddrFrom4([4]byte{69, 107, 7, 136}), 29),
	},
}
`
	// generate and compare
	w := &bytes.Buffer{}
	if err := generateRangesGo(w, map[string]regionsToPrefixes{"AWS": rtp}); err != nil {
		t.Fatalf("unexpected error generating: %v", err)
	}
	result := w.String()
	if result != goldenText {
		t.Error("result does not equal expected golden text")
		t.Error("golden text:")
		t.Error(goldenText)
		t.Error("result:")
		t.Error(result)
		t.Fail()
	}
}
