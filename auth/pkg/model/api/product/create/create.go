package create

type Product struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	SslSetup    struct {
		Company     string `json:"company"`
		Province    string `json:"province"`
		Country     string `json:"country"`
		CompanyUnit string `json:"company_unit"`
		CommonName  string `json:"common_name"`
		Locality    string `json:"locality"`
		AltDNS      string `json:"alt_dns"`
	} `json:"ssl_setup"`
	JwtConfig struct {
		AccessToken struct {
			ExpiresInMinutes int64 `json:"expires_in_minutes"`
		} `json:"access_token"`
		RefreshToken struct {
			ExpiresInMinutes int64 `json:"expires_in_minutes"`
		} `json:"refresh_token"`
	} `json:"jwt_config"`
}
