package create

type Product struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SslSetup    SslSetup  `json:"ssl_setup"`
	JwtConfig   JwtConfig `json:"jwt_config"`
}

type SslSetup struct {
	Company     string `json:"company"`
	Province    string `json:"province"`
	Country     string `json:"country"`
	CompanyUnit string `json:"company_unit"`
	CommonName  string `json:"common_name"`
	Locality    string `json:"locality"`
	AltDNS      string `json:"alt_dns"`
}

type JwtConfig struct {
	AccessToken  Token `json:"access_token"`
	RefreshToken Token `json:"refresh_token"`
}

type Token struct {
	ExpiresInMinutes int64 `json:"expires_in_minutes"`
}
