package cluster

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/redhatinsights/uhc-auth-proxy/requests/client"
	"github.com/spf13/viper"
)

// GetIdentity is a facade over all the steps required to get an Identity
func GetIdentity(wrapper client.Wrapper, reg Registration) (*Identity, error) {

	acct, err := GetCurrentAccount(wrapper, reg)
	if err != nil {
		return nil, fmt.Errorf("got an err when calling GetCurrentAccount: %w", err)
	}

	return &Identity{
		AccountNumber: acct.Organization.EbsAccountID,
		OrgID:         acct.Organization.ExternalID,
		Type:          "System",
		System: map[string]string{
			"cluster_id": reg.ClusterID,
		},
		Internal: Internal{
			OrgID: acct.Organization.ExternalID,
		},
	}, nil
}

// GetCurrentAccount uses a new flow with direct cluster tokenauth
func GetCurrentAccount(wrapper client.Wrapper, reg Registration) (*Account, error) {
	URL := viper.GetString("CURRENT_ACCOUNT_URL")

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}

	b, err := wrapper.Do(req, URL, reg.ClusterID, reg.AuthorizationToken)
	if err != nil {
		if b != nil {
			res := &AccountError{}
			if json.Unmarshal(b, res) == nil {
				res.Inner = err
				return nil, res
			}
		}

		return nil, err
	}

	res := &Account{}
	if err := json.Unmarshal(b, res); err != nil {
		return nil, err
	}

	return res, nil
}
