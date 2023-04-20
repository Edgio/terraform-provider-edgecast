resource "edgecast_waf_botmanager" "botmanager_1" {

	customer_id = "ABCDE"
	name = "My Bot Manager"
	bots_prod_id = "123bot1" # Must be an existing production Bot Rule
	actions {
		alert{
			name = "known_bot action"
		}
		custom_response{
			name = "ACTION"
			response_body_base64 = base64encode(file("response_body.html"))
			status = 403
			response_headers = {
				"header1" = "x-ec-rules"
				"header2" = "rejected"
			}
		}
		block_request{
			name = "known_bot action"
		}
		redirect_302{
			name = "known_bot action"
			url = "http://imouttahere.com"
		}
		browser_challenge{
			name = "known_bot action"
			is_custom_challenge = true
			response_body_base64 = base64encode(file("response_body.html"))
			status = 401
			valid_for_sec = 35
		}
		recaptcha{
			name = "my recaptcha"
			status = 401
			valid_for_sec = 35
			failed_action_type = "ALERT"
		}
	}
	exception_cookie = ["yummy-cookie", "yucky-cookie"]
	exception_ja3 = ["656b9a2f4de6ed4909e157482860ab3d"]
	exception_url = ["http://asdfasdfasd.com/"]
	exception_user_agent = ["abc/monkey/banana?abc=howmanybananas", "xyz/monkey/banana?abc=howmanybananas",]
	inspect_known_bots = true
	known_bot{
		action_type = "ALERT"
		bot_token = "google"
	}
	known_bot{
		action_type = "ALERT"
		bot_token = "facebook"
	}
	known_bot{
		action_type = "BLOCK_REQUEST"
		bot_token = "twitter"
	}
	known_bot{
		action_type = "CUSTOM_RESPONSE"
		bot_token = "yandex"
	}
	known_bot{
		action_type = "REDIRECT_302"
		bot_token = "semrush"
	}
	spoof_bot_action_type = "ALERT"
}
