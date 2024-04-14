package cfg

import "os"

type (
	CfgDB struct {
		HostDB     string
		UsernameDB string
		PasswordDB string
		NameDB     string
		PortDB     string
	}

	CfgTokens struct {
		AdminToken string
		UserToken  string
	}
)

func NewCfgDB() *CfgDB {
	return &CfgDB{}
}

func (cfg *CfgDB) SetFromEnv() {
	cfg.HostDB = os.Getenv("HOST_DB")

	cfg.UsernameDB = os.Getenv("USERNAME_DB")

	cfg.PasswordDB = os.Getenv("PASSWORD_DB")

	cfg.NameDB = os.Getenv("NAME_DB")

	cfg.PortDB = os.Getenv("PORT_DB")
}

func NewCfgTokens() *CfgTokens {
	return &CfgTokens{}
}

func (cfg *CfgTokens) SetFromEnv() {
	cfg.AdminToken = os.Getenv("ADMIN_TOKEN")

	cfg.UserToken = os.Getenv("USER_TOKEN")
}
