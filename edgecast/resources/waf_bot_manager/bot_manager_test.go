// Copyright 2023 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.
package waf_bot_manager_test

import (
	"reflect"
	"testing"

	"terraform-provider-edgecast/edgecast/resources/waf_bot_manager"

	"terraform-provider-edgecast/edgecast/helper"

	sdkbotmanager "github.com/EdgeCast/ec-sdk-go/edgecast/waf_bot_manager"
	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestFlattenActions(t *testing.T) {
	t.Parallel()
	base64Body := "PCFET0NUWVBFIGh0bWw+CjxodG1sIGxhbmc9ImVuIj4KICA8aGVhZD4KICAgIDxtZXRhIGNoYXJzZXQ9InV0Zi04Ii8+CiAgICA8dGl0bGU+NDAzIHVuYXV0aG9yaXplZDwvdGl0bGU+CiAgICA8bWV0YSBjb250ZW50PSI0MDMgdW5hdXRob3JpemVkIiBwcm9wZXJ0eT0ib2c6dGl0bGUiLz4KICAgIDxtZXRhIGNvbnRlbnQ9IndpZHRoPWRldmljZS13aWR0aCwgaW5pdGlhbC1zY2FsZT0xIiBuYW1lPSJ2aWV3cG9ydCIvPgogICAgPHN0eWxlPgogICAgICBib2R5IHsKICAgICAgICBmb250LWZhbWlseTogc2Fucy1zZXJpZjsKICAgICAgICBsaW5lLWhlaWdodDogMS4yOwogICAgICAgIGZvbnQtc2l6ZTogMThweDsKICAgICAgfQogICAgICBzZWN0aW9uIHsKICAgICAgICBtYXJnaW46IDAgYXV0bzsKICAgICAgICBtYXJnaW4tdG9wOiAxMzBweDsKICAgICAgICB3aWR0aDogNzUlOwogICAgICB9CiAgICAgIGgxIHsKICAgICAgICBmb250LXNpemU6IDUwcHg7CiAgICAgICAgbGluZS1oZWlnaHQ6IDQ1cHg7CiAgICAgICAgZm9udC13ZWlnaHQ6IDcwMDsKICAgICAgICBtYXJnaW4tYm90dG9tOiA3NXB4OwogICAgICAgIHdoaXRlLXNwYWNlOiBub3dyYXA7CiAgICAgIH0KICAgICAgcCB7CiAgICAgICAgbWFyZ2luLWJvdHRvbTogMTBweDsKICAgICAgfQogICAgICBzbWFsbCB7CiAgICAgICAgZm9udC1zaXplOiA4MCU7CiAgICAgICAgY29sb3I6ICMzMzM7CiAgICAgIH0KICAgICAgZm9vdGVyIHsKICAgICAgICBwb3NpdGlvbjogZml4ZWQ7CiAgICAgICAgYm90dG9tOiAwOwogICAgICAgIGxlZnQ6IDA7CiAgICAgICAgcGFkZGluZzogLjdyZW0gMCAuN3JlbSA0cmVtOwogICAgICAgIHdpZHRoOiAxMDAlOwogICAgICAgIGJhY2tncm91bmQ6ICBJbmRpZ287CiAgICAgIH0KICAgICAgZm9vdGVyIGEgewogICAgICAgIGNvbG9yOiB3aGl0ZTsKICAgICAgICBtYXJnaW4tbGVmdDogNDBweDsKICAgICAgICB0ZXh0LWRlY29yYXRpb246IG5vbmU7CiAgICAgIH0KICAgICAgLmQtbm9uZSB7CiAgICAgICAgZGlzcGxheTogbm9uZSAhaW1wb3J0YW50OwogICAgICB9CiAgICAgIC5zZWN0aW9uLWVycm9yIHsKICAgICAgICBjb2xvcjogI2JkMjQyNjsKICAgICAgfQogICAgICAubG9hZGluZyB7CiAgICAgICAgYW5pbWF0aW9uOiAzcyBpbmZpbml0ZSBzbGlkZWluOwogICAgICB9CiAgICAgIEBrZXlmcmFtZXMgc2xpZGVpbiB7CiAgICAgICAgZnJvbSB7CiAgICAgICAgICBtYXJnaW4tbGVmdDogMTAlOwogICAgICAgICAgY29sb3I6IHJnYigxMzQsIDUxLCAyNTUpOwogICAgICAgIH0KCiAgICAgICAgMzAlIHsKICAgICAgICAgIGNvbG9yOiByZ2IoMTM0LCA1MSwgMjU1KTsKICAgICAgICB9CgogICAgICAgIHRvIHsKICAgICAgICAgIG1hcmdpbi1sZWZ0OiAwJTsKICAgICAgICAgIGNvbG9yOiBibGFjazsKICAgICAgICB9CiAgICAgIH0KICAgIDwvc3R5bGU+CiAgPC9oZWFkPgoKICA8Ym9keSBvbmxvYWQ9Im9ubG9hZENvb2tpZUNoZWNrKCkiPgogICAge3tCT1RfSlN9fQogICAgPHNlY3Rpb24+CiAgICAgIDxoMT4KICAgICAgICBWYWxpZGF0aW5nIHlvdXIgYnJvd3NlciAtIGN1c3RvbSBjaGFsbGVuZ2UgcGFnZQogICAgICAgIDxzcGFuIGNsYXNzPSJsb2FkaW5nIj4uPC9zcGFuPjxzcGFuIGNsYXNzPSJsb2FkaW5nIj4uPC9zcGFuPjxzcGFuIGNsYXNzPSJsb2FkaW5nIj4uPC9zcGFuPgogICAgICA8L2gxPgoKICAgICAgPG5vc2NyaXB0PgogICAgICAgIDxoNCBjbGFzcz0ic2VjdGlvbi1lcnJvciI+UGxlYXNlIHR1cm4gSmF2YVNjcmlwdCBvbiBhbmQgcmVsb2FkIHRoZSBwYWdlLjwvaDQ+CiAgICAgIDwvbm9zY3JpcHQ+CgogICAgICA8ZGl2IGlkPSJjb29raWUtZXJyb3IiIGNsYXNzPSJkLW5vbmUiPgogICAgICAgIDxoNCBjbGFzcz0ic2VjdGlvbi1lcnJvciI+UGxlYXNlIGVuYWJsZSBjb29raWVzIGFuZCByZWxvYWQgdGhlIHBhZ2UuPC9oND4KICAgICAgPC9kaXY+CgogICAgICA8cD5UaGlzIG1heSB0YWtlIHVwIHRvIDUgc2Vjb25kczwvcD4KCiAgICAgIDxzbWFsbD5FdmVudCBJRDoge3tFVkVOVF9JRH19PC9zbWFsbD4KICAgIDwvc2VjdGlvbj4KCiAgICA8Zm9vdGVyPgogICAgICA8cD4KICAgICAgICA8YSBocmVmPSJodHRwczovL3d3dy5lZGdlY2FzdC5jb20vc2VjdXJpdHkvIj5Qb3dlcmVkIGJ5IEVkZ2lvPC9hPgogICAgICA8L3A+CiAgICA8L2Zvb3Rlcj4KICA8L2JvZHk+CiAgPHNjcmlwdD4KICAgIGZ1bmN0aW9uIG9ubG9hZENvb2tpZUNoZWNrKCkgewogICAgICBpZiAoIW5hdmlnYXRvci5jb29raWVFbmFibGVkKSB7CiAgICAgICAgZG9jdW1lbnQuZ2V0RWxlbWVudEJ5SWQoJ2Nvb2tpZS1lcnJvcicpLmNsYXNzTGlzdC5yZW1vdmUoJ2Qtbm9uZScpOwogICAgICB9CiAgICB9CiAgPC9zY3JpcHQ+CjwvaHRtbD4="

	tests := []struct {
		name string
		args sdkbotmanager.ActionObj
		want interface{}
	}{
		{
			name: "Happy Path",
			args: sdkbotmanager.ActionObj{
				ALERT: &sdkbotmanager.AlertAction{
					Id:   sdkbotmanager.PtrString("1"),
					Name: sdkbotmanager.PtrString("my alert action"),
				},
				CUSTOM_RESPONSE: &sdkbotmanager.CustomResponseAction{
					Id:                 sdkbotmanager.PtrString("2"),
					Name:               sdkbotmanager.PtrString("my custom_response action"),
					ResponseBodyBase64: &base64Body,
					Status:             sdkbotmanager.PtrInt32(403),
					ResponseHeaders: &map[string]string{
						"x-ec-rules": "rejected",
					},
				},
				BLOCK_REQUEST: &sdkbotmanager.BlockRequestAction{
					Id:   sdkbotmanager.PtrString("3"),
					Name: sdkbotmanager.PtrString("my block_request action"),
				},
				REDIRECT302: &sdkbotmanager.RedirectAction{
					Id:   sdkbotmanager.PtrString("4"),
					Name: sdkbotmanager.PtrString("my redirect_302 action"),
					Url:  sdkbotmanager.PtrString("http://imouttahere.com"),
				},
				BROWSER_CHALLENGE: &sdkbotmanager.BrowserChallengeAction{
					Id:                 sdkbotmanager.PtrString("5"),
					Name:               sdkbotmanager.PtrString("my browser_challenge action"),
					IsCustomChallenge:  sdkbotmanager.PtrBool(true),
					ResponseBodyBase64: &base64Body,
					ValidForSec:        sdkbotmanager.PtrInt32(3),
					Status:             sdkbotmanager.PtrInt32(401),
				},
			},
			want: []map[string]interface{}{
				{
					"alert": map[string]interface{}{
						"id":   "1",
						"name": "my alert action",
					},
					"custom_response": map[string]interface{}{
						"id":                   "2",
						"name":                 "my custom_response action",
						"response_body_base64": base64Body,
						"status":               int32(403),
						"response_headers": map[string]string{
							"x-ec-rules": "rejected",
						},
					},
					"block_request": map[string]interface{}{
						"id":   "3",
						"name": "my block_request action",
					},
					"redirect_302": map[string]interface{}{
						"id":   "4",
						"name": "my redirect_302 action",
						"url":  "http://imouttahere.com",
					},
					"browser_challenge": map[string]interface{}{
						"id":                   "5",
						"name":                 "my browser_challenge action",
						"is_custom_challenge":  true,
						"response_body_base64": base64Body,
						"valid_for_sec":        int32(3),
						"status":               int32(401),
					},
				},
			},
		},
		{
			name: "Empty input",
			args: sdkbotmanager.ActionObj{},
			// we expect a single empty map
			want: []map[string]interface{}{
				make(map[string]interface{}, 0),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := waf_bot_manager.FlattenActions(tt.args)

			if !reflect.DeepEqual(got, tt.want) {
				// deep.Equal doesn't compare pointer values, so we just use it
				// to generate a human friendly diff
				diff := deep.Equal(got, tt.want)
				t.Errorf("Diff: %+v", diff)
				t.Fatalf("%s: Expected %+v but got %+v",
					tt.name,
					tt.want,
					got,
				)
			}
		})
	}
}

func TestFlattenKnownBots(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []sdkbotmanager.KnownBotObj
		want []map[string]interface{}
	}{
		{
			name: "Happy Path",
			args: []sdkbotmanager.KnownBotObj{
				{
					ActionType: "ALERT",
					BotToken:   "google",
				},
				{
					ActionType: "ALERT",
					BotToken:   "facebook",
				},
				{
					ActionType: "BLOCK_REQUEST",
					BotToken:   "twitter",
				},
				{
					ActionType: "CUSTOM_RESPONSE",
					BotToken:   "yandex",
				},
				{
					ActionType: "REDIRECT_302",
					BotToken:   "semrush",
				},
			},
			want: []map[string]interface{}{
				{
					"action_type": "ALERT",
					"bot_token":   "google",
				},
				{
					"action_type": "ALERT",
					"bot_token":   "facebook",
				},
				{
					"action_type": "BLOCK_REQUEST",
					"bot_token":   "twitter",
				},
				{
					"action_type": "CUSTOM_RESPONSE",
					"bot_token":   "yandex",
				},
				{
					"action_type": "REDIRECT_302",
					"bot_token":   "semrush",
				},
			},
		},
		{
			name: "Empty input",
			args: make([]sdkbotmanager.KnownBotObj, 0),
			want: make([]map[string]interface{}, 0),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := waf_bot_manager.FlattenKnownBots(tt.args)

			if !reflect.DeepEqual(got, tt.want) {
				// deep.Equal doesn't compare pointer values, so we just use it
				// to generate a human friendly diff
				diff := deep.Equal(got, tt.want)
				t.Errorf("Diff: %+v", diff)
				t.Fatalf("%s: Expected %+v but got %+v",
					tt.name,
					tt.want,
					got,
				)
			}
		})
	}
}

func TestExpandActions(t *testing.T) {
	status403 := int32(403)
	status401 := int32(401)
	validforsec := int32(3)

	cases := []struct {
		name          string
		input         any
		expectedPtr   *sdkbotmanager.ActionObj
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: []interface{}{
				map[string]interface{}{
					"alert": []interface{}{
						map[string]interface{}{
							"id":   "1",
							"name": "known_bot action",
						},
					},
					"custom_response": []interface{}{
						map[string]any{
							"id":                   "1",
							"name":                 "known_bot action",
							"response_body_base64": "base64string",
							"status":               403,
							"response_headers": map[string]any{
								"header1": "x-ec-rules",
								"header2": "rejected",
							},
						},
					},
					"block_request": []interface{}{
						map[string]any{
							"id":   "1",
							"name": "known_bot action",
						},
					},
					"redirect_302": []interface{}{
						map[string]any{
							"id":   "1",
							"name": "known_bot action",
							"url":  "http://imouttahere.com",
						},
					},
					"browser_challenge": []interface{}{
						map[string]any{
							"id":                   "1",
							"name":                 "known_bot action",
							"is_custom_challenge":  true,
							"response_body_base64": "base64string",
							"status":               401,
							"valid_for_sec":        3,
						},
					},
				},
			},
			expectedPtr: &sdkbotmanager.ActionObj{
				ALERT: &sdkbotmanager.AlertAction{
					Id:   helper.WrapStringInPtr("1"),
					Name: helper.WrapStringInPtr("known_bot action"),
				},
				CUSTOM_RESPONSE: &sdkbotmanager.CustomResponseAction{
					Id:                 helper.WrapStringInPtr("1"),
					Name:               helper.WrapStringInPtr("known_bot action"),
					ResponseBodyBase64: helper.WrapStringInPtr("base64string"),
					Status:             &status403,
					ResponseHeaders: &map[string]string{
						"header1": "x-ec-rules",
						"header2": "rejected",
					},
				},
				BLOCK_REQUEST: &sdkbotmanager.BlockRequestAction{
					Id:   helper.WrapStringInPtr("1"),
					Name: helper.WrapStringInPtr("known_bot action"),
				},
				REDIRECT302: &sdkbotmanager.RedirectAction{
					Id:   helper.WrapStringInPtr("1"),
					Name: helper.WrapStringInPtr("known_bot action"),
					Url:  helper.WrapStringInPtr("http://imouttahere.com"),
				},
				BROWSER_CHALLENGE: &sdkbotmanager.BrowserChallengeAction{
					Id:                 helper.WrapStringInPtr("1"),
					Name:               helper.WrapStringInPtr("known_bot action"),
					IsCustomChallenge:  helper.WrapBoolInPtr(true),
					ResponseBodyBase64: helper.WrapStringInPtr("base64string"),
					Status:             &status401,
					ValidForSec:        &validforsec,
				},
			},
			expectSuccess: true,
		},
		{
			name:          "nil input",
			input:         nil,
			expectedPtr:   nil,
			expectSuccess: false,
		},
		{
			name:          "Error path - incorrect input type",
			input:         "not a []interface{}",
			expectedPtr:   nil,
			expectSuccess: false,
		},
	}

	for _, v := range cases {
		actualPtr, err := waf_bot_manager.ExpandActions(v.input)

		if v.expectSuccess {
			if err == nil {
				actual := *actualPtr
				expected := *v.expectedPtr

				if !reflect.DeepEqual(actual, expected) {
					// deep.Equal doesn't compare pointer values, so we just use it to
					// generate a human friendly diff
					diff := deep.Equal(actual, expected)
					t.Errorf("Diff: %+v", diff)
					t.Fatalf("%s: Expected %+v but got %+v",
						v.name,
						expected,
						actual,
					)
				}
			} else {
				t.Fatalf("%s: Encountered error where one was not expected: %+v",
					v.name,
					err,
				)
			}
		} else {
			if err == nil {
				t.Fatalf("%s: Expected error, but got no error", v.name)
			}
		}
	}
}

func TestExpandBotManager(t *testing.T) {

	cases := []struct {
		name          string
		input         map[string]any
		expectedPtr   *sdkbotmanager.BotManager
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: map[string]interface{}{
				"customer_id":  "ABC",
				"name":         "my bot manager",
				"bots_prod_id": "123",
				"exception_cookie": []any{
					"yummy-cookie",
					"yucky-cookie",
				},
				"exception_ja3": []any{
					"656b9a2f4de6ed4909e157482860ab3d",
				},
				"exception_url": []any{
					"myurl",
				},
				"exception_user_agent": []any{
					"useragent 1", "useragent 2",
				},
				"inspect_known_bots": true,
				"known_bot": []any{
					map[string]any{
						"action_type": "ALERT",
						"bot_token":   "google",
					},
					map[string]any{
						"action_type": "ALERT",
						"bot_token":   "facebook",
					},
				},
				"spoof_bot_action_type": "ALERT",
			},

			expectedPtr: &sdkbotmanager.BotManager{
				CustomerId:         helper.WrapStringInPtr("ABC"),
				Name:               helper.WrapStringInPtr("my bot manager"),
				BotsProdId:         helper.WrapStringInPtr("123"),
				ExceptionCookie:    []string{"yummy-cookie", "yucky-cookie"},
				ExceptionJa3:       []string{"656b9a2f4de6ed4909e157482860ab3d"},
				ExceptionUrl:       []string{"myurl"},
				ExceptionUserAgent: []string{"useragent 1", "useragent 2"},
				InspectKnownBots:   helper.WrapBoolInPtr(true),
				KnownBots: []sdkbotmanager.KnownBotObj{
					{
						ActionType: "ALERT",
						BotToken:   "google",
					},
					{
						ActionType: "ALERT",
						BotToken:   "facebook",
					},
				},
				SpoofBotActionType: helper.WrapStringInPtr("ALERT"),
			},
			expectSuccess: true,
		},
		{
			name:          "nil input",
			input:         nil,
			expectedPtr:   nil,
			expectSuccess: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var rd *schema.ResourceData
			if tt.input != nil {
				rd = schema.TestResourceDataRaw(
					t,
					waf_bot_manager.GetBotManagerSchema(),
					tt.input)
			}

			actualPtr, err := waf_bot_manager.ExpandBotManager(rd)

			if tt.expectSuccess {
				if len(err) == 0 {
					actual := actualPtr
					expected := tt.expectedPtr

					if !reflect.DeepEqual(actual, expected) {
						// deep.Equal doesn't compare pointer values, so we just use it to
						// generate a human friendly diff
						diff := deep.Equal(actual, expected)
						t.Errorf("Diff: %+v", diff)
						t.Fatalf("%s: Expected %+v but got %+v",
							tt.name,
							expected,
							actual,
						)
					}
				} else {
					t.Fatalf("%s: Encountered error where one was not expected: %+v",
						tt.name,
						err,
					)
				}
			} else {
				if err == nil {
					t.Fatalf("%s: Expected error, but got no error", tt.name)
				}
			}

		})
	}
}
