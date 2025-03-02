package utils

import (
	"errors"
	"os"
	"slices"
)

var ErrCredentialsNotFound = errors.New("credentials not found")
var ErrUnrecogniseSecretData = errors.New("unrecognise secret data, please read README")

func ParseEnv(envs []string) ([]string) {
	ret_data := []string{}
	for _, env := range envs {
		env := os.Getenv(env)
		if env == "" {
			ret_data = append(ret_data, "")
		} else {
			ret_data = append(ret_data, env)
		}
	}
	return ret_data
}

func Contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}
