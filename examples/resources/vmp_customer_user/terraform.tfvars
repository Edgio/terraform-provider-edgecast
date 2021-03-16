# Provides values the variables used in main.tf

# Use the credentials provided to you by Verizon Media
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

new_admin_user = {
    customer_account_number = "<customer account>"
    first_name = "<Admin First Name>"
    last_name = "<Admin Last Name>"
    email = "<Admin Email>"
    is_admin = true
}