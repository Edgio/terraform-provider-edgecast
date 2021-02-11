# Provides values the variables used in main.tf

# Use the credentials provided to you by Verizon Media
credentials = {
    api_token = "<API Token>"
    ids_client_secret = "<Client Secret>"
    ids_client_id = "<Client ID>"
    ids_scope = "<Scopes>"
}

test_customer_info = {
    account_number = "<Account Number>"
    customeruserid = "<Customer User ID>"
    portaltypeid = 1
}

# Valid values are "production" and "staging"
rulesEngineEnvironment = "staging"

# Below is an example of a Rules Engine policy defined as a string in a .tfvars file
# Note that main.tf in this directory reads directly from the file named "httplarge-policy.json"

/* # httplarge_policy = <<POLICYCREATE
    {
    'name':'test policy YYYYMMDD',
    'description':'This is a test policy.',
    'rules': [
        {
            'name':'test rule',
            'description': 'This is a test rule.',
            'matches': [{
                'type': 'match.origin.customer-origin.literal',
                'value': '/000000/Origin-X/',
                'features': [{
                    'type': 'feature.caching.compress-file-types',
                    'media-types': ['text/plain text/html', 'text/css application/x-javascript', 'text/javascript']
                }]
                }, {
                    'type': 'match.request.request-header.wildcard',
                    'name': 'User-Agent',
                    'result': 'nomatch',
                    'value': '*MSIE\\ 5*Mac* *MSIE\\ 4* *Mozilla/4* *compatible;*',
                    'ignore-case': 'True',
                    'features': [{
                        'type': 'feature.caching.compress-file-types',
                        'media-types': ['text/plain text/html', 'text/css application/x-javascript', 'text/javascript']
                    }]
                }]
        }
    ]
}
# POLICYCREATE */