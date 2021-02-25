# Provides values the variables used in main.tf

# Use the credentials provided to you by Verizon Media
credentials = {
    api_token = "AULdReDoB3gb0D7LNTx857NQvrcIKyvL"
    ids_client_secret = "CDbbMJw7FFJ11a7433ti1l9XgJHKr2Wk"
    ids_client_id = "31ef8e8f-0120-4112-8554-3eb11e83d58b"
    ids_scope = "ec.rules"

    # for internal testing
    api_address = "http://dev-api.edgecast.com"
    ids_address = "https://id-dev.vdms.io"
}

new_admin_user = {
    customer_account_number = "C1B6"
    first_name = "FirstAdmin1"
    last_name = "LastAdmin1"
    email = "admin.email02192021-4@test.com"
    is_admin = true
}
