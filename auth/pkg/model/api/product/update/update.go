package update

type Product struct {
	ProductId int64 `json:"product_id"`
	RenewSsl  *struct {
		Company     string `json:"company"`
		Province    string `json:"province"`
		Country     string `json:"country"`
		CompanyUnit string `json:"company_unit"`
		CommonName  string `json:"common_name"`
		Locality    string `json:"locality"`
		AltDNS      string `json:"alt_dns"`
	} `json:"renew_ssl,omitempty"`
	RenewRSA256 *struct {
		RenewKeyPair *bool `json:"renew_rsa_256,omitempty"`
		RenewConfig  *struct {
			AccessToken struct {
				ExpiresInMinutes int64 `json:"expires_in_minutes"`
			} `json:"access_token"`
			RefreshToken struct {
				ExpiresInMinutes int64 `json:"expires_in_minutes"`
			} `json:"refresh_token"`
		} `json:"renew_config"`
	} `json:"renew_rsa_256,omitempty"`
}
