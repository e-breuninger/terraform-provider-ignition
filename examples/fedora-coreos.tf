# Butane Config for Fedora CoreOS
data "ignition_config" "fedora-coreos-config" {
  content = templatefile("${path.module}/content/fcos.yaml", {
    message = "Hello World!"
  })
  strict       = true
  pretty_print = true

  snippets = [
    file("${path.module}/content/fcos-snippet.yaml"),
  ]
}

# Render as Ignition
resource "local_file" "fedora-coreos" {
  content  = data.ignition_config.fedora-coreos-config.rendered
  filename = "${path.module}/output/fedora-coreos.ign"
}
