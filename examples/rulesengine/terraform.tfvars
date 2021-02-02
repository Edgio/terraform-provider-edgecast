#Please update data in <> in order to run your terraform
provider_config = {
    #for pointing to staging environment, leave null to default to production
    api_address = ""
    # You must provide either an API Token or IDS credentials, but not both
    api_token = ""
    ids_client_secret = ""
    ids_client_id = ""
    ids_scope = ""
}

test_customer_info = {
    # optional params - for internal testing
    account_number = ""
    customeruserid = ""
    portaltypeid = ""
}

#deploy_to: Production or Staging
httplarge_policy = {
    deploy_to = "staging" 
    policy = <<POLICYCREATE
{
    '@type': 'policy',
    'name':'test policy-01272021-58',
    'description':'This is a test policy of PolicyCreate.',
    'state':'locked',
    'platform':'3',
    'rules': [
        {
            '@type':'rule',
            'name':'rule1',
            'description': 'This is a test rule26.',
            'matches': [{
                'type': 'match.always',
                'features': [{
                    'type': 'feature.comment',
                    'value': 'test'
                }]
            }]
        }
    ]
}
POLICYCREATE
}