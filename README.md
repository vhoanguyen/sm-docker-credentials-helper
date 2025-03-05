# SecretManager Docker Credential Helper

The SecretManager Docker Credential Helper is a [credential helper](https://github.com/docker/docker-credential-helpers) that makes it easier to use Docker with AWS Secrets Manager. When you run the helper, it will fetch the Docker credentials from AWS Secrets Manager and output them in a format that Docker can understand.


## Prerequisites

You must have AWS credentials configured on your machine. The helper will use these credentials to fetch the Docker credentials from AWS Secrets Manager.

You also need to create a secret in AWS Secrets Manager that contains the Docker credentials. The secret should be a JSON object with the following structure:

```json
{
  "https://index.docker.io/v1/":
    {
        "username": "username",
        "password": "TOKEN"
    },
  "ghcr.io":
    {
      "username": "username",
      "password": "TOKEN"
    },
  "registry.gitlab.com":
    {
      "username": "username",
      "password": "TOKEN"
    }
}
```

##  Installation

Build binary files (`VERSION=1.0.0 make build`) and copy to the docker credential helper path
```bash
cp ./bin/sm-login-darwin-arm64 /usr/local/bin/docker-credential-sm-login

OR

cp ./bin/sm-login-linux-amd64 /usr/local/bin/docker-credential-sm-login
```

## Environment variable


| Environment Variable         | Sample Value  | Description                                                        |
| ---------------------------- | ------------- | ------------------------------------------------------------------ |
| DOCKER_SECRET_NAME   | YOUR_SECRET_NAME | The name of the secret in AWS Secrets Manager that contains the Docker credentials. |
| AWS_PROFILE                  | YOUR_AWS_PROFILE | The name of the AWS profile to use. If you run locally |


## Configuration

~/.docker/config.json
```json
{
  "credsStore": "ecr-login",
  "credHelpers": {
    "https://index.docker.io/v1/": "sm-login",
    "registry.gitlab.com": "sm-login",
    "ghcr.io": "sm-login"
  }
}

```

## Local testing

```bash
% AWS_PROFILE=YOUR_AWS_PROFILE DOCKER_SECRET_NAME=YOUR_SECRET_NAME ./bin/sm-login-linux-amd64 list
```
```bash
echo https://index.docker.io/v1/ | AWS_PROFILE=YOUR_AWS_PROFILE DOCKER_SECRET_NAME=YOUR_SECRET_NAME ./bin/sm-login-linux-amd64 get

```
## Troubleshooting

Logs from the SecretManager Docker Credential Helper are stored in `~/.sm/`.