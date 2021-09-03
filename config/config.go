package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	CcpPath               string
	MspId                 string
	AdminId               string
	AdminSecret           string
	ClientId              string
	CretificateAuthorites string
	Affiliation           string
	Channel               string
	ChaincodeName         string
	AsLocalhost           string
	CertPath              string
	PrivateCertName       string
}

func (c *Config) loadEnv() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := godotenv.Load(wd + "/.env"); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	c.CcpPath = os.Getenv("CCP_PATH")
	c.MspId = os.Getenv("MSP_ID")
	c.AdminId = os.Getenv("ADMIN_ID")
	c.AdminSecret = os.Getenv("ADMIN_SECRET")
	c.ClientId = os.Getenv("CLIENT_ID")
	c.CretificateAuthorites = os.Getenv("CERTIFICATE_AUTORITIES")
	c.Affiliation = os.Getenv("AFFILIATION")
	c.Channel = os.Getenv("CHANNEL")
	c.ChaincodeName = os.Getenv("CHAINCODE_NAME")
	c.AsLocalhost = os.Getenv("AS_LOCALHOST")
	c.CertPath = os.Getenv("CERT_PATH")
	c.PrivateCertName = os.Getenv("PRIVATE_CERT_NAME")

	os.Setenv("DISCOVERY_AS_LOCALHOST", c.AsLocalhost)
}

func GetConfig() *Config {
	c := &Config{}
	c.loadEnv()
	return c
}
