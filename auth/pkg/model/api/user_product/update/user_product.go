package update

type UserProduct struct {
	UserProductId  int64           `json:"user_product_id"`
	RenewPassword  *RenewPassword  `json:"renew_password,omitempty"`
	ForgotPassword *ForgotPassword `json:"forgot_password,omitempty"`
	ForgotTwoFA    *ForgotTwoFA    `json:"forgot_two_fa,omitempty"`
}

type RenewPassword struct {
	OldPassword       string `json:"old_password"`
	NewPassword       string `json:"new_password"`
	RepeatNewPassword string `json:"repeat_new_password"`
}

type ForgotPassword struct {
	MasterPasswordKey string `json:"master_password_key"`
	NewPassword       string `json:"new_password"`
	RepeatNewPassword string `json:"repeat_new_password"`
}

type ForgotTwoFA struct {
	Password       string `json:"password"`
	MasterTwoFAKey string `json:"master_two_fa_key"`
}
