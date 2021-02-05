#Please update data in <> in order to run your terraform
credentials = {
    api_token = "4SJdxOPq0B7KZf7IVPXIeofwPNnwKDyk"
    ids_client_secret = "CDbbMJw7FFJ11a7433ti1l9XgJHKr2Wk"
    ids_client_id = "31ef8e8f-0120-4112-8554-3eb11e83d58b"
    ids_scope = "ec.rules ec.rules.admin ec.rules.deploy_dist"
}

test_customer_info = {
    account_number = "5F534"
    customeruserid = "349706"
    portaltypeid = 1
}

#valid values are "production" and "staging"
rulesEngineEnvironment = "staging"

# example policy that can be referenced in main.tf
# httplarge_policy = <<POLICYCREATE
#     {
#     '@type': 'policy',
#     'name':'test policy1182021-40',
#     'description':'This is a test policy of PolicyCreate.',
#     'state':'draft',
#     'rules': [
#         {
#             '@type':'rule',
#             'name':'test rule1',
#             'description': 'This is a test rule1.',
#             'matches': [{
#                 'type': 'match.origin.customer-origin.literal',
#                 'value': '/000000/Origin-X/',
#                 'features': [{
#                     'type': 'feature.caching.compress-file-types',
#                     'media-types': ['text/plain text/html', 'text/css application/x-javascript', 'text/javascript']
#                 }]
#                 }, {
#                     'type': 'match.request.request-header.wildcard',
#                     'name': 'User-Agent',
#                     'result': 'nomatch',
#                     'value': '*MSIE\\ 5*Mac* *MSIE\\ 4* *Mozilla/4* *compatible;*',
#                     'ignore-case': 'True',
#                     'features': [{
#                         'type': 'feature.caching.compress-file-types',
#                         'media-types': ['text/plain text/html', 'text/css application/x-javascript', 'text/javascript']
#                     }]
#                 }]
#         }
#     ]
# }
# POLICYCREATE