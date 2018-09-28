package main

import (
	"github.com/golang/glog"
	vaultapi "github.com/hashicorp/vault/api"
)

func VaultClientMaker(vaultAddr string) (client *vaultapi.Client, err error) {
	config := *vaultapi.DefaultConfig()

	config.Address = vaultAddr
	client, err = vaultapi.NewClient(&config)
	
	return client, err
}

func VaultTokenValidation(vaultAddr string, token string) (string, error) {
	glog.V(4).Infof("Validating token: %v", token)
	_, err := VaultClientMaker(vaultAddr)
	if err != nil {
		return "", err
	}
	
	// TODO - Look into security measures like setting TTL on Token
	//client.SetToken(token)
	//AuthClient := client.Auth().Token()
	//_, err = AuthClient.RenewSelf(600)
	//
	//if err != nil {
	//	return "", err
	//}
	
	return token, nil
}


func VaultGetSecret(vaultAddr string, path string, method string, token string) (map[string]interface{}, error) {
	client, err := VaultClientMaker(vaultAddr)
	if err != nil {
		glog.Fatalf("Error getting vault secrets: %s", err.Error())
	}
	
	client.SetToken(token)
	switch method {
	case "read":
		secret, err := client.Logical().Read(path)
		if err != nil {
			return nil, err
		}
		return secret.Data, nil
	case "write":
		Options := make(map[string]interface{})
		secret, err := client.Logical().Write(path, Options)
		if err != nil {
			return nil, err
		}
		return secret.Data, nil
	}
	
	return nil, err
}
