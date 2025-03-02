package main

import (
	"context"
	"fmt"
	"os"
	"sm-login/api"
	"sm-login/log"
	"sm-login/utils"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/sirupsen/logrus"
)

func init(){
	log.LogrusConfig()
	credentials.Name = "sm-login"
	credentials.Version = "0.1.0"
}

func main(){
	envs := utils.ParseEnv([]string{"AWS_REGION", "DOCKER_SECRET_NAME"})
	if envs[0] == "" {
		envs[0] = "ap-southeast-2"
	}
	if envs[1] == "" {
		fmt.Println("DOCKER_SECRET_NAME is not set")
		os.Exit(1)
	}
	sess, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(envs[0]))
	if err != nil {
		logrus.Errorf("failed to load config: %v", err)
		os.Exit(1)
	}
	secretClient := secretsmanager.NewFromConfig(sess)

	credentials.Serve(
		api.SecretManagerObject{
			SecretName: envs[1],
			SecretClientAPI: secretClient,
			Cxt: context.TODO(),
			Logger: logrus.StandardLogger(),
		})
}

