resource "ec_waf_managed_rule" "managed_rule_1" {
    account_number                = "<account_number>"
    name                          = "Terraform Managed Rule #1"
    ruleset_id                    = "ECRS"
    ruleset_version               = "2020-05-01"
    policies                      = [
        "r4020_tw_cpanel.conf.json",
        "r4040_tw_drupal.conf.json",
        "r4030_tw_iis.conf.json",
        "r4070_tw_joomla.conf.json",
        "r4050_tw_microsoft_sharepoint.conf.json",
        "r4060_tw_wordpress.conf.json",
        "r5040_cross_site_scripting.conf.json",
        "r2000_ec_custom_rule.conf.json",
        "r5021_http_attack.conf.json",
        "r5020_http_protocol_violation.conf.json",
        "r5043_java_attack.conf.json",
        "r5030_local_file_inclusion.conf.json",
        "r5033_php_injection.conf.json",
        "r5032_remote_code_execution.conf.json",
        "r5031_remote_file_inclusion.conf.json",
        "r5010_scanner_detection.conf.json",
        "r5042_session_fixation.conf.json",
        "r5041_sql_injection.conf.json",
        "r6000_blocking_evaluation.conf.json"
    ]

    disabled_rules {
                policy_id = "r2000_ec_custom_rule.conf.json"
                rule_id = "431003"
            } 

    general_settings {
        anomaly_threshold = 10
        arg_length =  8000
        arg_name_length = 1024
        combined_file_sizes = 6291456
        ignore_cookie = ["ignoredCookie"]
        ignore_header = ["ignoredHeaders"]
        ignore_query_args = ["ignoredQuery"]
        json_parser = true
        max_num_args = 512
        paranoia_level = 1
        process_request_body = true
        response_header_name = "X-EC-Security-Audit"
        total_arg_length = 64000
        validate_utf8_encoding = true
        xml_parser = true
    }

    rule_target_updates {
        is_regex = true
        target = "ARGS"
        target_match = "ignoredArgumentException"
        rule_id = "431000"
    }
    
    rule_target_updates {
        is_regex = true
        target = "REQUEST_COOKIES"
        target_match = "ignoredCookiesException"
        rule_id = "431000"
    }

    rule_target_updates  {
        is_regex = true
        target = "REQUEST_HEADERS"
        target_match = "ignoredHeaderException"
        rule_id = "431000"
    }
}