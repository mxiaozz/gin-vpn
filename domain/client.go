package domain

type Client struct {
	UserName    string `json:"name"`
	CertContent string `json:"cert"`
	CertExpire  int    `json:"expire"`
}
