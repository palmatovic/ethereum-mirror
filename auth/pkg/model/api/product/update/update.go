package update

type Product struct {
	ProductId   int64  `json:"product_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SslSetup    *struct {
		Company     string `json:"company"`
		Province    string `json:"province"`
		Country     string `json:"country"`
		CompanyUnit string `json:"company_unit"`
		CommonName  string `json:"common_name"`
		Locality    string `json:"locality"`
		AltDNS      string `json:"alt_dns"`
	} `json:"ssl_setup,omitempty"`
	JwtConfig *struct {
		RenewRSA256      bool `json:"renew_rsa_256"`
		RenewTokenConfig *struct {
			AccessToken struct {
				ExpiresInMinutes int64 `json:"expires_in_minutes"`
			} `json:"access_token"`
			RefreshToken struct {
				ExpiresInMinutes int64 `json:"expires_in_minutes"`
			} `json:"refresh_token"`
		} `json:"renew_token_config"`
	} `json:"jwt_config,omitempty"`
	RenewAES256 bool `json:"renew_aes_256"`
}
