// Copyright 2023 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.
package waf_bot_manager_test

import (
	"reflect"
	"testing"

	"terraform-provider-edgecast/edgecast/resources/waf_bot_manager"

	botmanager "github.com/EdgeCast/ec-sdk-go/edgecast/waf_bot_manager"
	"github.com/go-test/deep"
)

func TestFlattenActions(t *testing.T) {
	t.Parallel()
	base64Body := "PCFET0NUWVBFIGh0bWw+CjxodG1sIGxhbmc9ImVuIj4KICA8aGVhZD4KICAgIDxtZXRhIGNoYXJzZXQ9InV0Zi04Ii8+CiAgICA8dGl0bGU+NDAzIHVuYXV0aG9yaXplZDwvdGl0bGU+CiAgICA8bWV0YSBjb250ZW50PSI0MDMgdW5hdXRob3JpemVkIiBwcm9wZXJ0eT0ib2c6dGl0bGUiLz4KICAgIDxtZXRhIGNvbnRlbnQ9IndpZHRoPWRldmljZS13aWR0aCwgaW5pdGlhbC1zY2FsZT0xIiBuYW1lPSJ2aWV3cG9ydCIvPgogICAgPHN0eWxlPgogICAgICBib2R5IHsKICAgICAgICBmb250LWZhbWlseTogc2Fucy1zZXJpZjsKICAgICAgICBsaW5lLWhlaWdodDogMS4yOwogICAgICAgIGZvbnQtc2l6ZTogMThweDsKICAgICAgfQogICAgICBzZWN0aW9uIHsKICAgICAgICBtYXJnaW46IDAgYXV0bzsKICAgICAgICBtYXJnaW4tdG9wOiAxMzBweDsKICAgICAgICB3aWR0aDogNzUlOwogICAgICB9CiAgICAgIGgxIHsKICAgICAgICBmb250LXNpemU6IDUwcHg7CiAgICAgICAgbGluZS1oZWlnaHQ6IDQ1cHg7CiAgICAgICAgZm9udC13ZWlnaHQ6IDcwMDsKICAgICAgICBtYXJnaW4tYm90dG9tOiA3NXB4OwogICAgICAgIHdoaXRlLXNwYWNlOiBub3dyYXA7CiAgICAgIH0KICAgICAgcCB7CiAgICAgICAgbWFyZ2luLWJvdHRvbTogMTBweDsKICAgICAgfQogICAgICBzbWFsbCB7CiAgICAgICAgZm9udC1zaXplOiA4MCU7CiAgICAgICAgY29sb3I6ICMzMzM7CiAgICAgIH0KICAgICAgZm9vdGVyIHsKICAgICAgICBwb3NpdGlvbjogZml4ZWQ7CiAgICAgICAgYm90dG9tOiAwOwogICAgICAgIGxlZnQ6IDA7CiAgICAgICAgcGFkZGluZzogLjdyZW0gMCAuN3JlbSA0cmVtOwogICAgICAgIHdpZHRoOiAxMDAlOwogICAgICAgIGJhY2tncm91bmQ6ICBJbmRpZ287CiAgICAgIH0KICAgICAgZm9vdGVyIGEgewogICAgICAgIGNvbG9yOiB3aGl0ZTsKICAgICAgICBtYXJnaW4tbGVmdDogNDBweDsKICAgICAgICB0ZXh0LWRlY29yYXRpb246IG5vbmU7CiAgICAgIH0KICAgICAgLmQtbm9uZSB7CiAgICAgICAgZGlzcGxheTogbm9uZSAhaW1wb3J0YW50OwogICAgICB9CiAgICAgIC5zZWN0aW9uLWVycm9yIHsKICAgICAgICBjb2xvcjogI2JkMjQyNjsKICAgICAgfQogICAgICAubG9hZGluZyB7CiAgICAgICAgYW5pbWF0aW9uOiAzcyBpbmZpbml0ZSBzbGlkZWluOwogICAgICB9CiAgICAgIEBrZXlmcmFtZXMgc2xpZGVpbiB7CiAgICAgICAgZnJvbSB7CiAgICAgICAgICBtYXJnaW4tbGVmdDogMTAlOwogICAgICAgICAgY29sb3I6IHJnYigxMzQsIDUxLCAyNTUpOwogICAgICAgIH0KCiAgICAgICAgMzAlIHsKICAgICAgICAgIGNvbG9yOiByZ2IoMTM0LCA1MSwgMjU1KTsKICAgICAgICB9CgogICAgICAgIHRvIHsKICAgICAgICAgIG1hcmdpbi1sZWZ0OiAwJTsKICAgICAgICAgIGNvbG9yOiBibGFjazsKICAgICAgICB9CiAgICAgIH0KICAgIDwvc3R5bGU+CiAgPC9oZWFkPgoKICA8Ym9keSBvbmxvYWQ9Im9ubG9hZENvb2tpZUNoZWNrKCkiPgogICAge3tCT1RfSlN9fQogICAgPHNlY3Rpb24+CiAgICAgIDxoMT4KICAgICAgICBWYWxpZGF0aW5nIHlvdXIgYnJvd3NlciAtIGN1c3RvbSBjaGFsbGVuZ2UgcGFnZQogICAgICAgIDxzcGFuIGNsYXNzPSJsb2FkaW5nIj4uPC9zcGFuPjxzcGFuIGNsYXNzPSJsb2FkaW5nIj4uPC9zcGFuPjxzcGFuIGNsYXNzPSJsb2FkaW5nIj4uPC9zcGFuPgogICAgICA8L2gxPgoKICAgICAgPG5vc2NyaXB0PgogICAgICAgIDxoNCBjbGFzcz0ic2VjdGlvbi1lcnJvciI+UGxlYXNlIHR1cm4gSmF2YVNjcmlwdCBvbiBhbmQgcmVsb2FkIHRoZSBwYWdlLjwvaDQ+CiAgICAgIDwvbm9zY3JpcHQ+CgogICAgICA8ZGl2IGlkPSJjb29raWUtZXJyb3IiIGNsYXNzPSJkLW5vbmUiPgogICAgICAgIDxoNCBjbGFzcz0ic2VjdGlvbi1lcnJvciI+UGxlYXNlIGVuYWJsZSBjb29raWVzIGFuZCByZWxvYWQgdGhlIHBhZ2UuPC9oND4KICAgICAgPC9kaXY+CgogICAgICA8cD5UaGlzIG1heSB0YWtlIHVwIHRvIDUgc2Vjb25kczwvcD4KCiAgICAgIDxzbWFsbD5FdmVudCBJRDoge3tFVkVOVF9JRH19PC9zbWFsbD4KICAgIDwvc2VjdGlvbj4KCiAgICA8Zm9vdGVyPgogICAgICA8cD4KICAgICAgICA8YSBocmVmPSJodHRwczovL3d3dy5lZGdlY2FzdC5jb20vc2VjdXJpdHkvIj5Qb3dlcmVkIGJ5IEVkZ2lvPC9hPgogICAgICA8L3A+CiAgICA8L2Zvb3Rlcj4KICA8L2JvZHk+CiAgPHNjcmlwdD4KICAgIGZ1bmN0aW9uIG9ubG9hZENvb2tpZUNoZWNrKCkgewogICAgICBpZiAoIW5hdmlnYXRvci5jb29raWVFbmFibGVkKSB7CiAgICAgICAgZG9jdW1lbnQuZ2V0RWxlbWVudEJ5SWQoJ2Nvb2tpZS1lcnJvcicpLmNsYXNzTGlzdC5yZW1vdmUoJ2Qtbm9uZScpOwogICAgICB9CiAgICB9CiAgPC9zY3JpcHQ+CjwvaHRtbD4="

	tests := []struct {
		name string
		args botmanager.ActionObj
		want interface{}
	}{
		{
			name: "Happy Path",
			args: botmanager.ActionObj{
				ALERT: &botmanager.AlertAction{
					Id:   botmanager.PtrString("1"),
					Name: botmanager.PtrString("my alert action"),
				},
				CUSTOM_RESPONSE: &botmanager.CustomResponseAction{
					Id:                 botmanager.PtrString("2"),
					Name:               botmanager.PtrString("my custom_response action"),
					ResponseBodyBase64: &base64Body,
					Status:             botmanager.PtrInt32(403),
					ResponseHeaders: &map[string]string{
						"x-ec-rules": "rejected",
					},
				},
				BLOCK_REQUEST: &botmanager.BlockRequestAction{
					Id:   botmanager.PtrString("3"),
					Name: botmanager.PtrString("my block_request action"),
				},
				REDIRECT302: &botmanager.RedirectAction{
					Id:   botmanager.PtrString("4"),
					Name: botmanager.PtrString("my redirect_302 action"),
					Url:  botmanager.PtrString("http://imouttahere.com"),
				},
				BROWSER_CHALLENGE: &botmanager.BrowserChallengeAction{
					Id:                 botmanager.PtrString("5"),
					Name:               botmanager.PtrString("my browser_challenge action"),
					IsCustomChallenge:  botmanager.PtrBool(true),
					ResponseBodyBase64: &base64Body,
					ValidForSec:        botmanager.PtrInt32(3),
					Status:             botmanager.PtrInt32(401),
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
			args: botmanager.ActionObj{},
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
		args []botmanager.KnownBotObj
		want []map[string]interface{}
	}{
		{
			name: "Happy Path",
			args: []botmanager.KnownBotObj{
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
			args: make([]botmanager.KnownBotObj, 0),
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
