# Provides values the variables used in main.tf

# Use the credentials provided to you by EdgeCast
credentials = {
    #api_token = "c07f60fb-add0-4526-b96d-7ef2de704a66" #prod
    #api_token = "bf42f045-ff24-43c4-b012-af9379b6faaa" #QA PCC
    api_token = "UHuHfnBLmG6pySwgqLGWLMCSNzeU3ogE" #QA Manual PCC
    ids_client_secret = "<Client Secret>"
    ids_client_id = "<Client ID>"
    ids_scope = "<Scopes>"

    # for internal testing
    api_address = null
    ids_address = null
    api_address_legacy = "https://qa-api.edgecast.com"
    #api_address_legacy = null
}