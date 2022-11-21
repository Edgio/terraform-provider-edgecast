resource "edgecast_originv3_group" "group_httplarge_1" {

	name = "TestTFOriginGroup"
	platform = "http-large"
	host_header = "edgecastcdn.net"
	shield_pops = ["AMB", "BLR"]
	network_type_id = 2
	strict_pci_certified = "false"

	tls_settings {
        sni_hostname = "origin.example.com"
        allow_self_signed = "false"
        public_keys_to_verify = ["ff8b4a82b08ea5f7be124e6b4363c00d7462655f","c571398b01fce46a8a177abdd6174dfee6137358"]
    }
}

resource "edgecast_originv3_group" "example" {
}