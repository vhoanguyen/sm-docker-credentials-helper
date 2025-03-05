package api

import (
	"context"
	"encoding/json"
	"fmt"
	"sm-login/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/sirupsen/logrus"
)


var supportedURLs = []string{"https://index.docker.io/v1/"} // only support docker hub

type SecretData struct {
	Username string
	Password string
}

type SecretManagerAPI interface {
	GetSecretValue(
		ctx context.Context,
		input *secretsmanager.GetSecretValueInput,
		optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)

}

type SecretManagerObject struct {
	SecretName   string
	SecretClientAPI SecretManagerAPI
	Cxt          context.Context
	Logger       *logrus.Logger
}

//lint:ignore ST1006
func (self SecretManagerObject) Add(creds *credentials.Credentials) error {
	self.Logger.Error("Add is not implemented")
	return nil
}

//lint:ignore ST1006
func (sm SecretManagerObject) Delete(serverURL string) error {
	sm.Logger.Error("Delete is not implemented")
	return nil
}

//lint:ignore ST1006
func (self SecretManagerObject) Store(creds *credentials.Credentials) error {
	self.Logger.Error("Store is not implemented")
	return nil
}

//lint:ignore ST1006
func (self SecretManagerObject) Get(serverURL string) (string, string, error) {
	if !utils.Contains(supportedURLs, serverURL) {
		self.Logger.Errorf("unsupported serverURL: %s", serverURL)
		return "", "", fmt.Errorf("unsupported serverURL: %s", serverURL)
	}
	secret, err := self.getSecret()
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "ResourceNotFoundException":
				self.Logger.Errorf("ResourceNotFoundException: %v", aerr.Error())
			default:
				self.Logger.Errorf(aerr.Error())
			}
		}
		return "", "", err

	}
	data, ok := secret[serverURL]
	if !ok {
		self.Logger.Errorf("secret not found for serverURL: %s", serverURL)
		return "", "", fmt.Errorf("secret not found for serverURL: %s", serverURL)
	}
	return data.Username, data.Password, nil
}

//lint:ignore ST1006
func (self SecretManagerObject) List() (map[string]string, error) {
	secret, err := self.getSecret()
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "ResourceNotFoundException":
				self.Logger.Errorf("ResourceNotFoundException: %v", aerr.Error())
			default:
				self.Logger.Errorf(aerr.Error())
			}
		}
		return nil, err

	}
	result := make(map[string]string)
	for serverURL, data := range secret {
		result[serverURL] = data.Username
	}
	return result, nil

}


//lint:ignore ST1006
func (self SecretManagerObject) getSecret() (map[string]SecretData, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: &self.SecretName,
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}
	secret, err := self.SecretClientAPI.GetSecretValue(
		self.Cxt,
		input,
	)
	if err != nil {
		if self.Logger != nil {
			self.Logger.Errorf("failed to get secret: %v", err.Error())
		}
		return nil, err
	}
	var secretData map[string]SecretData
	err = json.Unmarshal([]byte(*secret.SecretString), &secretData)
	if err != nil {
		if self.Logger != nil {
			self.Logger.Errorf("Unrecognise Secret Data, please read README: %v", err)
		}
		return nil, err
	}
	return secretData, nil
}