// Copyright 2023 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.
package data

import (
	"fmt"
	"os"
	"time"

	"terraform-provider-edgecast/test/integration/cmd/populate/config"
	"terraform-provider-edgecast/test/integration/cmd/populate/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"github.com/EdgeCast/ec-sdk-go/edgecast/waf"
	"github.com/EdgeCast/ec-sdk-go/edgecast/waf_bot_manager"
)

const (
	ruleProcessingMaxRetries = 4
	rulewWaitTimeSeconds     = 15
)

func createWAFBotManagerData(cfg config.Config) BotManagerResult {
	wafSvc := internal.Check(waf.New(cfg.SDKConfig))

	// bug work around - need to clear out IDS credentials
	sdkCfgNoIDS := cfg.SDKConfig
	sdkCfgNoIDS.IDSCredentials = edgecast.IDSCredentials{}
	wbmSvc := internal.Check(waf_bot_manager.New(sdkCfgNoIDS))

	id := createBotManager(wafSvc, wbmSvc, cfg.AccountNumber)

	return BotManagerResult{
		BotManagerID: id,
	}
}

func createBotManager(
	wafSvc *waf.WafService,
	wbmSvc *waf_bot_manager.Service,
	accountNumber string,
) string {
	botRuleID := createBotRule(wafSvc, accountNumber)

	customerID := accountNumber
	base64Body := "PCFET0NUWVBFIGh0bWw+CjxodG1sIGxhbmc9ImVuIj4KICA8aGVhZD4KICAgIDxtZXRhIGNoYXJzZXQ9InV0Zi04Ii8+CiAgICA8dGl0bGU+NDAzIHVuYXV0aG9yaXplZDwvdGl0bGU+CiAgICA8bWV0YSBjb250ZW50PSI0MDMgdW5hdXRob3JpemVkIiBwcm9wZXJ0eT0ib2c6dGl0bGUiLz4KICAgIDxtZXRhIGNvbnRlbnQ9IndpZHRoPWRldmljZS13aWR0aCwgaW5pdGlhbC1zY2FsZT0xIiBuYW1lPSJ2aWV3cG9ydCIvPgogICAgPHN0eWxlPgogICAgICBib2R5IHsKICAgICAgICBmb250LWZhbWlseTogc2Fucy1zZXJpZjsKICAgICAgICBsaW5lLWhlaWdodDogMS4yOwogICAgICAgIGZvbnQtc2l6ZTogMThweDsKICAgICAgfQogICAgICBzZWN0aW9uIHsKICAgICAgICBtYXJnaW46IDAgYXV0bzsKICAgICAgICBtYXJnaW4tdG9wOiAxMzBweDsKICAgICAgICB3aWR0aDogNzUlOwogICAgICB9CiAgICAgIGgxIHsKICAgICAgICBmb250LXNpemU6IDUwcHg7CiAgICAgICAgbGluZS1oZWlnaHQ6IDQ1cHg7CiAgICAgICAgZm9udC13ZWlnaHQ6IDcwMDsKICAgICAgICBtYXJnaW4tYm90dG9tOiA3NXB4OwogICAgICAgIHdoaXRlLXNwYWNlOiBub3dyYXA7CiAgICAgIH0KICAgICAgcCB7CiAgICAgICAgbWFyZ2luLWJvdHRvbTogMTBweDsKICAgICAgfQogICAgICBzbWFsbCB7CiAgICAgICAgZm9udC1zaXplOiA4MCU7CiAgICAgICAgY29sb3I6ICMzMzM7CiAgICAgIH0KICAgICAgZm9vdGVyIHsKICAgICAgICBwb3NpdGlvbjogZml4ZWQ7CiAgICAgICAgYm90dG9tOiAwOwogICAgICAgIGxlZnQ6IDA7CiAgICAgICAgcGFkZGluZzogLjdyZW0gMCAuN3JlbSA0cmVtOwogICAgICAgIHdpZHRoOiAxMDAlOwogICAgICAgIGJhY2tncm91bmQ6ICBJbmRpZ287CiAgICAgIH0KICAgICAgZm9vdGVyIGEgewogICAgICAgIGNvbG9yOiB3aGl0ZTsKICAgICAgICBtYXJnaW4tbGVmdDogNDBweDsKICAgICAgICB0ZXh0LWRlY29yYXRpb246IG5vbmU7CiAgICAgIH0KICAgICAgLmQtbm9uZSB7CiAgICAgICAgZGlzcGxheTogbm9uZSAhaW1wb3J0YW50OwogICAgICB9CiAgICAgIC5zZWN0aW9uLWVycm9yIHsKICAgICAgICBjb2xvcjogI2JkMjQyNjsKICAgICAgfQogICAgICAubG9hZGluZyB7CiAgICAgICAgYW5pbWF0aW9uOiAzcyBpbmZpbml0ZSBzbGlkZWluOwogICAgICB9CiAgICAgIEBrZXlmcmFtZXMgc2xpZGVpbiB7CiAgICAgICAgZnJvbSB7CiAgICAgICAgICBtYXJnaW4tbGVmdDogMTAlOwogICAgICAgICAgY29sb3I6IHJnYigxMzQsIDUxLCAyNTUpOwogICAgICAgIH0KCiAgICAgICAgMzAlIHsKICAgICAgICAgIGNvbG9yOiByZ2IoMTM0LCA1MSwgMjU1KTsKICAgICAgICB9CgogICAgICAgIHRvIHsKICAgICAgICAgIG1hcmdpbi1sZWZ0OiAwJTsKICAgICAgICAgIGNvbG9yOiBibGFjazsKICAgICAgICB9CiAgICAgIH0KICAgIDwvc3R5bGU+CiAgPC9oZWFkPgoKICA8Ym9keSBvbmxvYWQ9Im9ubG9hZENvb2tpZUNoZWNrKCkiPgogICAge3tCT1RfSlN9fQogICAgPHNlY3Rpb24+CiAgICAgIDxoMT4KICAgICAgICBWYWxpZGF0aW5nIHlvdXIgYnJvd3NlciAtIGN1c3RvbSBjaGFsbGVuZ2UgcGFnZQogICAgICAgIDxzcGFuIGNsYXNzPSJsb2FkaW5nIj4uPC9zcGFuPjxzcGFuIGNsYXNzPSJsb2FkaW5nIj4uPC9zcGFuPjxzcGFuIGNsYXNzPSJsb2FkaW5nIj4uPC9zcGFuPgogICAgICA8L2gxPgoKICAgICAgPG5vc2NyaXB0PgogICAgICAgIDxoNCBjbGFzcz0ic2VjdGlvbi1lcnJvciI+UGxlYXNlIHR1cm4gSmF2YVNjcmlwdCBvbiBhbmQgcmVsb2FkIHRoZSBwYWdlLjwvaDQ+CiAgICAgIDwvbm9zY3JpcHQ+CgogICAgICA8ZGl2IGlkPSJjb29raWUtZXJyb3IiIGNsYXNzPSJkLW5vbmUiPgogICAgICAgIDxoNCBjbGFzcz0ic2VjdGlvbi1lcnJvciI+UGxlYXNlIGVuYWJsZSBjb29raWVzIGFuZCByZWxvYWQgdGhlIHBhZ2UuPC9oND4KICAgICAgPC9kaXY+CgogICAgICA8cD5UaGlzIG1heSB0YWtlIHVwIHRvIDUgc2Vjb25kczwvcD4KCiAgICAgIDxzbWFsbD5FdmVudCBJRDoge3tFVkVOVF9JRH19PC9zbWFsbD4KICAgIDwvc2VjdGlvbj4KCiAgICA8Zm9vdGVyPgogICAgICA8cD4KICAgICAgICA8YSBocmVmPSJodHRwczovL3d3dy5lZGdlY2FzdC5jb20vc2VjdXJpdHkvIj5Qb3dlcmVkIGJ5IEVkZ2lvPC9hPgogICAgICA8L3A+CiAgICA8L2Zvb3Rlcj4KICA8L2JvZHk+CiAgPHNjcmlwdD4KICAgIGZ1bmN0aW9uIG9ubG9hZENvb2tpZUNoZWNrKCkgewogICAgICBpZiAoIW5hdmlnYXRvci5jb29raWVFbmFibGVkKSB7CiAgICAgICAgZG9jdW1lbnQuZ2V0RWxlbWVudEJ5SWQoJ2Nvb2tpZS1lcnJvcicpLmNsYXNzTGlzdC5yZW1vdmUoJ2Qtbm9uZScpOwogICAgICB9CiAgICB9CiAgPC9zY3JpcHQ+CjwvaHRtbD4="

	botManager := waf_bot_manager.BotManager{
		Name:               internal.Pointer(internal.Unique("-botmanager")),
		SpoofBotActionType: internal.Pointer("ALERT"),
		Actions: &waf_bot_manager.ActionObj{
			ALERT: &waf_bot_manager.AlertAction{
				Name:    internal.Pointer("known_bot action"),
				EnfType: internal.Pointer("ALERT"),
			},
			BLOCK_REQUEST: &waf_bot_manager.BlockRequestAction{
				Name:    internal.Pointer("known_bot action"),
				EnfType: internal.Pointer("BLOCK_REQUEST"),
			},
			BROWSER_CHALLENGE: &waf_bot_manager.BrowserChallengeAction{
				Name:               internal.Pointer("known_bot action"),
				EnfType:            internal.Pointer("BROWSER_CHALLENGE"),
				IsCustomChallenge:  waf_bot_manager.PtrBool(true),
				ResponseBodyBase64: &base64Body,
				Status:             waf_bot_manager.PtrInt32(401),
				ValidForSec:        waf_bot_manager.PtrInt32(3),
			},
			CUSTOM_RESPONSE: &waf_bot_manager.CustomResponseAction{
				Name:               internal.Pointer("ACTION"),
				EnfType:            internal.Pointer("CUSTOM_RESPONSE"),
				ResponseBodyBase64: &base64Body,
				Status:             waf_bot_manager.PtrInt32(403),
				ResponseHeaders: &map[string]string{
					"x-ec-rules": "rejected",
				},
			},
			REDIRECT302: &waf_bot_manager.RedirectAction{
				Name:    internal.Pointer("known_bot action"),
				EnfType: internal.Pointer("REDIRECT_302"),
				Url:     internal.Pointer("http://imouttahere.com"),
			},
		},
		BotsProdId:      &botRuleID,
		CustomerId:      &customerID,
		ExceptionCookie: []string{"yummy-cookie", "yucky-cookie"},
		ExceptionJa3:    []string{"656b9a2f4de6ed4909e157482860ab3d"},
		ExceptionUrl:    []string{"http://asdfasdfasd.com/"},
		ExceptionUserAgent: []string{
			"abc/monkey/banana?abc=howmanybananas",
			"xyz/monkey/banana?abc=howmanybananas",
		},
		InspectKnownBots: waf_bot_manager.PtrBool(true),
		KnownBots: []waf_bot_manager.KnownBotObj{
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
	}

	params := waf_bot_manager.CreateBotManagerParams{
		CustId:         accountNumber,
		BotManagerInfo: botManager,
	}

	// Need to wait for the bot rule processing - retry until completeion.

	var id string

	for i := 1; i <= ruleProcessingMaxRetries; i++ {
		resp, err := wbmSvc.BotManagers.CreateBotManager(params)
		if err == nil {
			id = *resp.Id
			break
		}

		if i == ruleProcessingMaxRetries {
			fmt.Println("bot manager create failed, out of retries")
			os.Exit(1)
		}

		fmt.Printf(
			"bot manager create failed, retrying in %d seconds...\n",
			rulewWaitTimeSeconds)
		time.Sleep(rulewWaitTimeSeconds * time.Second)
	}

	return id
}
