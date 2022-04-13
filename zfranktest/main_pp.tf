terraform {
  required_providers {
    ec = {
      source  = "github.com/terraform-providers/ec"
      version = "0.4.7"
    }
  }
}
 
provider "ec" {
    api_token = "xRz[REDACTED]8Vl"
    ids_client_secret = "iicuMk4Ahv3P2fzlVuoqUFXMtBkjAz7M"
    ids_client_id = "a4f72e37-881f-4eb0-b360-82cdb627530d"
    ids_scope = "ec.rules"
}
 