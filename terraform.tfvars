# Provides values the variables used in main.tf

# Use the credentials provided to you by EdgeCast
credentials = {
    #api_token = "c07f60fb-add0-4526-b96d-7ef2de704a66" #prod
    api_token = "bf42f045-ff24-43c4-b012-af9379b6faaa" #QA PCC
    #api_token = "UHuHfnBLmG6pySwgqLGWLMCSNzeU3ogE" #QA Manual PCC
    #api_token = "8WyEvKm52I0iRetecA64gZcsN90wwxsV" # Frank 4FDBB
    ids_client_secret = "movIX3KoFzA5qjscGyPxesbBVjvvkp51"
    ids_client_id = "1cc1a3dd-12da-474e-9b46-178595e1480d"
    ids_scope = "ec.rules"

    # for internal testing
    api_address = null
    ids_address = null
    api_address_legacy = "https://qa-api.edgecast.com"
    #api_address_legacy = null
}