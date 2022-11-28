resource "edgecast_originv3_httplarge" "group_1" {

	name = "TestTFOriginGroup"
	host_header = "edgecastcdn.net"
	shield_pops = ["AMB", "BLR"]
	network_type_id = 2
	strict_pci_certified = false

	tls_settings {
        sni_hostname = "origin.example.com"
        allow_self_signed = false
        public_keys_to_verify = ["ff8b4a82b08ea5f7be124e6b4363c00d7462655f","c571398b01fce46a8a177abdd6174dfee6137358"]
    }
	
	origin {
    	name = "marketing-origin-entry-a"
    	host = "https://cdn-la.example.com"
    	port = 443
		is_primary = true
		storage_type_id = 1
		protocol_type_id = 2
	}
	origin {
    	name = "marketing-origin-entry-b"
    	host = "https://cdn-lb.example.com"
    	port = 443
		is_primary = true
		storage_type_id = 1
		protocol_type_id = 2
	}
}