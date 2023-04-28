# terraform-provider-ignition

`terraform-provider-ignition` allows Terraform to validate a [Butane configuration](https://coreos.github.io/butane/specs/) and transpile it to an [Ignition configuration](https://coreos.github.io/ignition/).

The Butane configuration is transpiled to the corresponding ignition version according to the [Butane specification](https://coreos.github.io/butane/specs/#butane-specifications-and-ignition-specifications).  
The Butane version is taken from the `content` attribute.

The Ignition versions `3.0.0`, `3.1.0`, `3.2.0`, `3.3.0` and `3.4.0` are supported.

## Usage

Configure the config transpiler provider (e.g. `providers.tf`).

```tf
provider "ignition" {}

terraform {
  required_providers {
    ignition = {
      source  = "e-breuninger/ignition"
      version = "1.0.0"
    }
  }
}
```

Define a Butane config for Fedora CoreOS or Flatcar Linux:

```yaml
variant: fcos
version: 1.5.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - ssh-key foo
```

```yaml
variant: flatcar
version: 1.1.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - ssh-key foo
```

Define a `ignition_config` data source with strict validation.

```tf
data "ignition_config" "worker" {
  content      = file("worker.yaml")
  strict       = true
  pretty_print = false

  snippets = [
    file("units.yaml"),
    file("storage.yaml"),
  ]
}
```

Optionally, template the `content`.

```tf
data "ignition_config" "worker" {
  content = templatefile("worker.yaml", {
    ssh_authorized_key = "ssh-ed25519 AAAA...",
  })
  strict       = true
}
```

Render the `ignition_config` as Ignition for use by machine instances.

```tf
resource "aws_instance" "worker" {
  user_data = data.ignition_config.worker.rendered
}
```

Run `terraform init` to ensure plugin version requirements are met.

```
$ terraform init
```

## Requirements

* Terraform v0.13+ [installed](https://www.terraform.io/downloads.html)

## Development

### Binary

To develop the provider plugin locally, build an executable with Go v1.18+.

```
make
```

## Credits

This provider is a fork of the terraform provider [poseidon/ct](https://github.com/poseidon/terraform-provider-ct)
