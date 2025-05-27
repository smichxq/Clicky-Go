package model

import (
	"fmt"
	"net/url"
)

type Code2SessionReq struct {
	Base      string
	AppId     string
	Secret    string
	JsCode    string
	GrantType string
}

func (c2s *Code2SessionReq) GetReqUrl() (string, error) {
	// 1. Parse the base URL. This helps validate if the base is a valid URL structure.
	parsedBaseUrl, err := url.Parse(c2s.Base)
	if err != nil {
		// Return an empty string and an error if the base URL is invalid.
		// Wrapping the error provides more context.
		return "", fmt.Errorf("invalid base URL '%s': %w", c2s.Base, err)
	}

	// 2. Create a new url.Values map to hold the query parameters.
	params := url.Values{}

	// 3. Add the parameters to the map.
	// The keys here match the desired output format (e.g., "js_code").
	params.Add("appid", c2s.AppId)
	params.Add("secret", c2s.Secret)
	params.Add("js_code", c2s.JsCode)
	params.Add("grant_type", c2s.GrantType)

	// 4. Encode the parameters and assign them to the RawQuery field of the parsed URL.
	// The Encode() method automatically handles URL-escaping special characters.
	parsedBaseUrl.RawQuery = params.Encode()

	// 5. Return the string representation of the fully constructed URL.
	return parsedBaseUrl.String(), nil
}

type Code2SessionResp struct {
	ErrCode    int32  `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
	Rid        string `json:"rid,omitempty"`
	SessionKey string `json:"session_key"`
	OpenId     string `json:"openid"`
	UnionId    string `json:"unionid,omitempty"` // UnionId is optional, so we use the omitempty tag
}
