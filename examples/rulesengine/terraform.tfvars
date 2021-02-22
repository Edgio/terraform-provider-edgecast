# Provides values the variables used in main.tf

# Use the credentials provided to you by Verizon Media
credentials = {
    api_token = "<API Token>"
    ids_client_secret = "CDbbMJw7FFJ11a7433ti1l9XgJHKr2Wk"
    ids_client_id = "31ef8e8f-0120-4112-8554-3eb11e83d58b"
    ids_scope = "ec.rules"
    api_address = "http://dev-api.edgecast.com"
    ids_address = "https://id-dev.vdms.io"
}

test_customer_info = {
    account_number = "C1B6"
    customeruserid = "133172"
    portaltypeid = 1
}

# Valid values are "production" and "staging"
rulesEngineEnvironment = "staging"
policy = null