package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	dotenv "sdk-tester/config"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

var env = dotenv.GetConfig()

func main() {
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		os.Exit(1)
	}

	if !wallet.Exists(env.ClientId) {
		err = populateWallet(wallet)
		if err != nil {
			fmt.Printf("Failed to populate wallet contents: %s\n", err)
			os.Exit(1)
		}
	}

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(env.CcpPath)),
		gateway.WithIdentity(wallet, env.ClientId),
	)
	if err != nil {
		fmt.Printf("failed to connect to gateway: %v", err)
		os.Exit(1)
	}
	defer gw.Close()

	result, err := query(gw)
	if err != nil {
		fmt.Printf("failed to query to contract: %v", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

}

func query(gw *gateway.Gateway) ([]byte, error) {

	network, err := gw.GetNetwork(env.Channel)
	if err != nil {
		return nil, fmt.Errorf("failed to get network: %v", err)
	}

	contract := network.GetContract(env.ChaincodeName)

	result, err := contract.EvaluateTransaction("AdminOperations", "GetBalanceOf", "{\"args\":[\"-\",\"c32cf14e0a11d271e131ecf3d5ac0f8c19d01b23\"]}")
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %v", err)
	}
	return result, err
}

func populateWallet(wallet *gateway.Wallet) error {
	certPath := filepath.Join(env.CertPath, "signcerts", env.PrivateCertName)
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(env.CertPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity(env.MspId, string(cert), string(key))

	err = wallet.Put(env.ClientId, identity)
	if err != nil {
		return err
	}
	return nil
}

func AddAffiliationOrg(caClient *msp.Client, caName string) error {
	affl := env.Affiliation

	fmt.Println("Initializing Affiliation for " + affl)

	affResponse, err := caClient.GetAffiliation(affl)

	if affResponse != nil && err != nil {

		fmt.Println("Affiliation Exists")

		AfInfo := affResponse.AffiliationInfo
		CAName := affResponse.CAName

		fmt.Println("AfInfo : " + AfInfo.Name)
		fmt.Println("CAName : " + CAName)
	} else {

		fmt.Println("Add Affiliation " + affl)

		_, err = caClient.AddAffiliation(&msp.AffiliationRequest{
			Name:   affl,
			Force:  true,
			CAName: caName,
		})

		if err != nil {
			return fmt.Errorf("Failed to add affiliation for CA '%s' : %v ", caName, err)
		}
	}
	fmt.Println("\n Affiliation completed successfully")
	return nil
}
