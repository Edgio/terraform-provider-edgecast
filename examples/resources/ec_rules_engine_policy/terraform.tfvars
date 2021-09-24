# Provides values the variables used in main.tf

# Use the credentials provided to you by Edgecast
credentials = {
    api_token = "<API Token>"
    ids_client_secret = "<Client Secret>"
    ids_client_id = "<Client ID>"
    ids_scope = "<Scopes>"

    # for internal testing
    api_address = null
    api_address_legacy = null
    ids_address = null
}

test_customer_info = {
    account_number = "C1B6"
    customeruserid = "133172"
    portaltypeid = 1
}

# Valid values are "production" and "staging"
rulesEngineEnvironment = "staging"
policy = null