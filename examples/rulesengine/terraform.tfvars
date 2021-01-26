#Please update data in <> in order to run your terraform
partner_info = {
    #for pointing to staging environment, leave null to default to production
    api_address = ""
    # You must provide either an API Token or IDS credentials, but not both
    api_token = null
    ids_client_secret = ""
    ids_client_id = ""
    ids_scope = ""
}

customer_info = {
    account_number = ""
    customeruserid = ""
    portaltypeid = 0
}

httplarge_policy = <<POLICYCREATE
    {
    '@type': 'policy',
    'name':'test policy1182021-40',
    'description':'This is a test policy of PolicyCreate.',
    'state':'draft',
    'rules': [
        {
            '@type':'rule',
            'name':'test rule1',
            'description': 'This is a test rule1.',
            'matches': [{
                'type': 'match.origin.customer-origin.literal',
                'value': '/8017DA1/Origin-NGUX/',
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
POLICYCREATE
