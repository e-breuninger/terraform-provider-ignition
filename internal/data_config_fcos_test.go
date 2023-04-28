package internal

import (
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Fedora CoreOS variant, v1.5.0

const fedoraCoreOSV15Resource = `
data "ignition_config" "fedora-coreos" {
  pretty_print = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.5.0
storage:
  luks:
    - name: data
      device: /dev/vdb
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
}
`

const fedoraCoreOSV15WithSnippets = `
data "ignition_config" "fedora-coreos-snippets" {
  pretty_print = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.5.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.5.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

const fedoraCoreOSV15WithSnippetsPrettyFalse = `
data "ignition_config" "fedora-coreos-snippets" {
  pretty_print = false
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.5.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.5.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

func TestButaneConfig_FCOSv1_5(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSV15Resource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ignition_config.fedora-coreos", "rendered", ignitionV34Expected),
				),
			},
			{
				Config: fedoraCoreOSV15WithSnippets,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ignition_config.fedora-coreos-snippets", "rendered", ignitionV34WithSnippetsExpected),
				),
			},
			{
				Config: fedoraCoreOSV15WithSnippetsPrettyFalse,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ignition_config.fedora-coreos-snippets", "rendered", ignitionV34WithSnippetsPrettyFalseExpected),
				),
			},
		},
	})
}

// Fedora CoreOS variant, v1.4.0

const fedoraCoreOSV14Resource = `
data "ignition_config" "fedora-coreos" {
  pretty_print = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.4.0
storage:
  luks:
    - name: data
      device: /dev/vdb
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
}
`
const ignitionV32Expected = `{
  "ignition": {
    "config": {
      "replace": {
        "verification": {}
      }
    },
    "proxy": {},
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "3.2.0"
  },
  "passwd": {
    "users": [
      {
        "name": "core",
        "sshAuthorizedKeys": [
          "key"
        ]
      }
    ]
  },
  "storage": {
    "luks": [
      {
        "device": "/dev/vdb",
        "keyFile": {
          "verification": {}
        },
        "name": "data"
      }
    ]
  },
  "systemd": {}
}`

const ignitionV33Expected = `{
  "ignition": {
    "config": {
      "replace": {
        "verification": {}
      }
    },
    "proxy": {},
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "3.3.0"
  },
  "kernelArguments": {},
  "passwd": {
    "users": [
      {
        "name": "core",
        "sshAuthorizedKeys": [
          "key"
        ]
      }
    ]
  },
  "storage": {
    "luks": [
      {
        "clevis": {
          "custom": {}
        },
        "device": "/dev/vdb",
        "keyFile": {
          "verification": {}
        },
        "name": "data"
      }
    ]
  },
  "systemd": {}
}`

const ignitionV34Expected = `{
  "ignition": {
    "config": {
      "replace": {
        "verification": {}
      }
    },
    "proxy": {},
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "3.4.0"
  },
  "kernelArguments": {},
  "passwd": {
    "users": [
      {
        "name": "core",
        "sshAuthorizedKeys": [
          "key"
        ]
      }
    ]
  },
  "storage": {
    "luks": [
      {
        "clevis": {
          "custom": {}
        },
        "device": "/dev/vdb",
        "keyFile": {
          "verification": {}
        },
        "name": "data"
      }
    ]
  },
  "systemd": {}
}`

const fedoraCoreOSV14WithSnippets = `
data "ignition_config" "fedora-coreos-snippets" {
  pretty_print = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.4.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.4.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

const ignitionV30WithSnippetsExpected = `{
  "ignition": {
    "config": {
      "replace": {
        "source": null,
        "verification": {}
      }
    },
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "3.0.0"
  },
  "passwd": {
    "users": [
      {
        "name": "core",
        "sshAuthorizedKeys": [
          "key"
        ]
      }
    ]
  },
  "storage": {},
  "systemd": {
    "units": [
      {
        "enabled": true,
        "name": "docker.service"
      }
    ]
  }
}`

const ignitionV31WithSnippetsExpected = `{
  "ignition": {
    "config": {
      "replace": {
        "verification": {}
      }
    },
    "proxy": {},
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "3.1.0"
  },
  "passwd": {
    "users": [
      {
        "name": "core",
        "sshAuthorizedKeys": [
          "key"
        ]
      }
    ]
  },
  "storage": {},
  "systemd": {
    "units": [
      {
        "enabled": true,
        "name": "docker.service"
      }
    ]
  }
}`

const ignitionV32WithSnippetsExpected = `{
  "ignition": {
    "config": {
      "replace": {
        "verification": {}
      }
    },
    "proxy": {},
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "3.2.0"
  },
  "passwd": {
    "users": [
      {
        "name": "core",
        "sshAuthorizedKeys": [
          "key"
        ]
      }
    ]
  },
  "storage": {},
  "systemd": {
    "units": [
      {
        "enabled": true,
        "name": "docker.service"
      }
    ]
  }
}`

const ignitionV33WithSnippetsExpected = `{
  "ignition": {
    "config": {
      "replace": {
        "verification": {}
      }
    },
    "proxy": {},
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "3.3.0"
  },
  "kernelArguments": {},
  "passwd": {
    "users": [
      {
        "name": "core",
        "sshAuthorizedKeys": [
          "key"
        ]
      }
    ]
  },
  "storage": {},
  "systemd": {
    "units": [
      {
        "enabled": true,
        "name": "docker.service"
      }
    ]
  }
}`

const ignitionV34WithSnippetsExpected = `{
  "ignition": {
    "config": {
      "replace": {
        "verification": {}
      }
    },
    "proxy": {},
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "3.4.0"
  },
  "kernelArguments": {},
  "passwd": {
    "users": [
      {
        "name": "core",
        "sshAuthorizedKeys": [
          "key"
        ]
      }
    ]
  },
  "storage": {},
  "systemd": {
    "units": [
      {
        "enabled": true,
        "name": "docker.service"
      }
    ]
  }
}`

const fedoraCoreOSV14WithSnippetsPrettyFalse = `
data "ignition_config" "fedora-coreos-snippets" {
  pretty_print = false
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.4.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.4.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`
const ignitionV30WithSnippetsPrettyFalseExpected = `{"ignition":{"config":{"replace":{"source":null,"verification":{}}},"security":{"tls":{}},"timeouts":{},"version":"3.0.0"},"passwd":{"users":[{"name":"core","sshAuthorizedKeys":["key"]}]},"storage":{},"systemd":{"units":[{"enabled":true,"name":"docker.service"}]}}`

const ignitionV31WithSnippetsPrettyFalseExpected = `{"ignition":{"config":{"replace":{"verification":{}}},"proxy":{},"security":{"tls":{}},"timeouts":{},"version":"3.1.0"},"passwd":{"users":[{"name":"core","sshAuthorizedKeys":["key"]}]},"storage":{},"systemd":{"units":[{"enabled":true,"name":"docker.service"}]}}`

const ignitionV32WithSnippetsPrettyFalseExpected = `{"ignition":{"config":{"replace":{"verification":{}}},"proxy":{},"security":{"tls":{}},"timeouts":{},"version":"3.2.0"},"passwd":{"users":[{"name":"core","sshAuthorizedKeys":["key"]}]},"storage":{},"systemd":{"units":[{"enabled":true,"name":"docker.service"}]}}`

const ignitionV33WithSnippetsPrettyFalseExpected = `{"ignition":{"config":{"replace":{"verification":{}}},"proxy":{},"security":{"tls":{}},"timeouts":{},"version":"3.3.0"},"kernelArguments":{},"passwd":{"users":[{"name":"core","sshAuthorizedKeys":["key"]}]},"storage":{},"systemd":{"units":[{"enabled":true,"name":"docker.service"}]}}`

const ignitionV34WithSnippetsPrettyFalseExpected = `{"ignition":{"config":{"replace":{"verification":{}}},"proxy":{},"security":{"tls":{}},"timeouts":{},"version":"3.4.0"},"kernelArguments":{},"passwd":{"users":[{"name":"core","sshAuthorizedKeys":["key"]}]},"storage":{},"systemd":{"units":[{"enabled":true,"name":"docker.service"}]}}`

func TestButaneConfig_FCOSv1_4(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSV14Resource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ignition_config.fedora-coreos", "rendered", ignitionV33Expected),
				),
			},
			{
				Config: fedoraCoreOSV14WithSnippets,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ignition_config.fedora-coreos-snippets", "rendered", ignitionV33WithSnippetsExpected),
				),
			},
			{
				Config: fedoraCoreOSV14WithSnippetsPrettyFalse,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ignition_config.fedora-coreos-snippets", "rendered", ignitionV33WithSnippetsPrettyFalseExpected),
				),
			},
		},
	})
}

// Fedora CoreOS variant, v1.3.0

const fedoraCoreOSV13Resource = `
data "ignition_config" "fedora-coreos" {
  pretty_print = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.3.0
storage:
  luks:
    - name: data
      device: /dev/vdb
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
}
`

const fedoraCoreOSV13WithSnippets = `
data "ignition_config" "fedora-coreos-snippets" {
  pretty_print = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.3.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.3.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

const fedoraCoreOSV13WithSnippetsPrettyFalse = `
data "ignition_config" "fedora-coreos-snippets" {
  pretty_print = false
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.3.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.3.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

func TestButaneConfig_FCOSv1_3(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSV13Resource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ignition_config.fedora-coreos", "rendered", ignitionV32Expected),
				),
			},
			{
				Config: fedoraCoreOSV13WithSnippets,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ignition_config.fedora-coreos-snippets", "rendered", ignitionV32WithSnippetsExpected),
				),
			},
			{
				Config: fedoraCoreOSV13WithSnippetsPrettyFalse,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ignition_config.fedora-coreos-snippets", "rendered", ignitionV32WithSnippetsPrettyFalseExpected),
				),
			},
		},
	})
}

// Fedora CoreOS variant, v1.2.0

const fedoraCoreOSV12Resource = `
data "ignition_config" "fedora-coreos" {
  pretty_print = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.2.0
storage:
  luks:
    - name: data
      device: /dev/vdb
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
}
`

const fedoraCoreOSV12WithSnippets = `
data "ignition_config" "fedora-coreos-snippets" {
  pretty_print = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.2.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.2.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

const fedoraCoreOSV12WithSnippetsPrettyFalse = `
data "ignition_config" "fedora-coreos-snippets" {
  pretty_print = false
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.2.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.2.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

func TestButaneConfig_FCOSv1_2(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSV12Resource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ignition_config.fedora-coreos", "rendered", ignitionV32Expected),
				),
			},
			{
				Config: fedoraCoreOSV12WithSnippets,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ignition_config.fedora-coreos-snippets", "rendered", ignitionV32WithSnippetsExpected),
				),
			},
			{
				Config: fedoraCoreOSV12WithSnippetsPrettyFalse,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ignition_config.fedora-coreos-snippets", "rendered", ignitionV32WithSnippetsPrettyFalseExpected),
				),
			},
		},
	})
}

// Fedora CoreOS variant, v1.1.0

const fedoraCoreOSV11Resource = `
data "ignition_config" "fedora-coreos" {
  pretty_print = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.1.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
}
`

const ignitionV30Expected = `{
  "ignition": {
    "config": {
      "replace": {
        "source": null,
        "verification": {}
      }
    },
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "3.0.0"
  },
  "passwd": {
    "users": [
      {
        "name": "core",
        "sshAuthorizedKeys": [
          "key"
        ]
      }
    ]
  },
  "storage": {},
  "systemd": {}
}`

const ignitionV31Expected = `{
  "ignition": {
    "config": {
      "replace": {
        "verification": {}
      }
    },
    "proxy": {},
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "3.1.0"
  },
  "passwd": {
    "users": [
      {
        "name": "core",
        "sshAuthorizedKeys": [
          "key"
        ]
      }
    ]
  },
  "storage": {},
  "systemd": {}
}`

const fedoraCoreOSV11WithSnippets = `
data "ignition_config" "fedora-coreos-snippets" {
  pretty_print = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.1.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.1.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

const fedoraCoreOSV11WithSnippetsPrettyFalse = `
data "ignition_config" "fedora-coreos-snippets" {
  pretty_print = false
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.1.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.1.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

func TestButaneConfig_FCOSv1_1(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSV11Resource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ignition_config.fedora-coreos", "rendered", ignitionV31Expected),
				),
			},
			{
				Config: fedoraCoreOSV11WithSnippets,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ignition_config.fedora-coreos-snippets", "rendered", ignitionV31WithSnippetsExpected),
				),
			},
			{
				Config: fedoraCoreOSV11WithSnippetsPrettyFalse,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ignition_config.fedora-coreos-snippets", "rendered", ignitionV31WithSnippetsPrettyFalseExpected),
				),
			},
		},
	})
}

// Fedora CoreOS variant, v1.0.0

const fedoraCoreOSV10Resource = `
data "ignition_config" "fedora-coreos" {
  pretty_print = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.0.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
}
`

const fedoraCoreOSV10WithSnippets = `
data "ignition_config" "fedora-coreos-snippets" {
  pretty_print = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.0.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.0.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

const fedoraCoreOSV10WithSnippetsPrettyFalse = `
data "ignition_config" "fedora-coreos-snippets" {
  pretty_print = false
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.0.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.0.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

func TestButaneConfig_FCOSv1_0(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSV10Resource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ignition_config.fedora-coreos", "rendered", ignitionV30Expected),
				),
			},
			{
				Config: fedoraCoreOSV10WithSnippets,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ignition_config.fedora-coreos-snippets", "rendered", ignitionV30WithSnippetsExpected),
				),
			},
			{
				Config: fedoraCoreOSV10WithSnippetsPrettyFalse,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ignition_config.fedora-coreos-snippets", "rendered", ignitionV30WithSnippetsPrettyFalseExpected),
				),
			},
		},
	})
}

const fedoraCoreOSMixSnippetBehind = `
data "ignition_config" "fedora-coreos-mix-versions" {
  pretty_print = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.4.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.2.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

const ignitionV33MixExpected = `{
  "ignition": {
    "config": {
      "replace": {
        "verification": {}
      }
    },
    "proxy": {},
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "3.3.0"
  },
  "kernelArguments": {},
  "passwd": {
    "users": [
      {
        "name": "core",
        "sshAuthorizedKeys": [
          "key"
        ]
      }
    ]
  },
  "storage": {},
  "systemd": {
    "units": [
      {
        "enabled": true,
        "name": "docker.service"
      }
    ]
  }
}`

func TestFedoraCoreOSMix_SnippetBehind(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSMixSnippetBehind,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ignition_config.fedora-coreos-mix-versions", "rendered", ignitionV33MixExpected),
				),
			},
		},
	})
}

const fedoraCoreOSMixSnippetAhead = `
data "ignition_config" "fedora-coreos-mix-versions" {
  pretty_print = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.2.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.4.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

func TestFedoraCoreOSMixVersions_SnippetAhead(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config:      fedoraCoreOSMixSnippetAhead,
				ExpectError: regexp.MustCompile(`snippet version 3\.3\.0 is newer than content version 3\.2\.0 and therefore incompatible`),
			},
		},
	})
}

const invalidResource = `
data "ignition_config" "invalid" {
  content = "foo"
  strict = true
  some_invalid_field = "strict-mode-will-reject"
}
`

func TestInvalidResource(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config:      invalidResource,
				ExpectError: regexp.MustCompile(`An argument named "some_invalid_field" is not expected here`),
			},
		},
	})
}
