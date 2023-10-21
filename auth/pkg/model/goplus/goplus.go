package goplus

//
//type TokenSecurityResult struct {
//	ResultType                 string `json:"result_type"`
//	AntiWhaleModifiable        string `json:"anti_whale_modifiable"`
//	BuyTax                     string `json:"buy_tax"`
//	CanTakeBackOwnership       string `json:"can_take_back_ownership"`
//	CannotBuy                  string `json:"cannot_buy"`
//	CannotSellAll              string `json:"cannot_sell_all"`
//	CreatorAddress             string `json:"creator_address"`
//	CreatorBalance             string `json:"creator_balance"`
//	CreatorPercent             string `json:"creator_percent"`
//	Dex                        []DexItem
//	ExternalCall               string `json:"external_call"`
//	HiddenOwner                string `json:"hidden_owner"`
//	HolderCount                string `json:"holder_count"`
//	Holders                    []HoldersItem
//	HoneypotWithSameCreator    string `json:"honeypot_with_same_creator"`
//	IsAirdropScam              string `json:"is_airdrop_scam"`
//	IsAntiWhale                string `json:"is_antiwhale"`
//	IsBlacklisted              string `json:"is_blacklisted"`
//	IsHoneypot                 string `json:"is_honeypot"`
//	IsInDex                    string `json:"is_in_dex"`
//	IsMintable                 string `json:"is_mintable"`
//	IsOpenSource               string `json:"is_open_source"`
//	IsProxy                    string `json:"is_proxy"`
//	IsTrueToken                string `json:"is_true_token"`
//	IsWhitelisted              string `json:"is_whitelisted"`
//	LpHolderCount              string `json:"lp_holder_count"`
//	LpHolders                  []HoldersItem
//	LpTotalSupply              string `json:"lp_total_supply"`
//	Note                       string `json:"note"`
//	OtherPotentialRisks        string `json:"other_potential_risks"`
//	OwnerAddress               string `json:"owner_address"`
//	OwnerBalance               string `json:"owner_balance"`
//	OwnerChangeBalance         string `json:"owner_change_balance"`
//	OwnerPercent               string `json:"owner_percent"`
//	PersonalSlippageModifiable string `json:"personal_slippage_modifiable"`
//	Selfdestruct               string `json:"selfdestruct"`
//	SellTax                    string `json:"sell_tax"`
//	SlippageModifiable         string `json:"slippage_modifiable"`
//	TokenName                  string `json:"token_name"`
//	TokenSymbol                string `json:"token_symbol"`
//	TotalSupply                string `json:"total_supply"`
//	TradingCooldown            string `json:"trading_cooldown"`
//	TransferPausable           string `json:"transfer_pausable"`
//	TrustList                  string `json:"trust_list"`
//}
//
//type DexItem struct {
//	Liquidity string `json:"liquidity"`
//	Name      string `json:"name"`
//	Pair      string `json:"pair"`
//}
//
//type HoldersItem struct {
//	Address      string        `json:"address"`
//	Balance      string        `json:"balance"`
//	IsContract   int32         `json:"is_contract"`
//	IsLocked     int32         `json:"is_locked"`
//	LockedDetail []interface{} `json:"locked_detail"`
//	Percent      string        `json:"percent"`
//	Tag          string        `json:"tag"`
//}
//
///*
//type LockedDetailItem struct {
//}*/
