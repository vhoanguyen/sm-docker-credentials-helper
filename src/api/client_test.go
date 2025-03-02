package api

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type SecretTestData struct {
	Username string
	Password string
}

// MockSecretsManagerGetSecretAPI is a mock of SecretsManagerGetSecretAPI interface
type MockSecretsManagerGetSecretAPI struct {
	ctrl     *gomock.Controller
	recorder *MockSecretsManagerGetSecretAPIMockRecorder
}

// MockSecretsManagerGetSecretAPIMockRecorder is the mock recorder for MockSecretsManagerGetSecretAPI
type MockSecretsManagerGetSecretAPIMockRecorder struct {
	mock *MockSecretsManagerGetSecretAPI
}

// NewMockSecretsManagerGetSecretAPI creates a new mock instance
func NewMockSecretsManagerGetSecretAPI(ctrl *gomock.Controller) *MockSecretsManagerGetSecretAPI {
	mock := &MockSecretsManagerGetSecretAPI{ctrl: ctrl}
	mock.recorder = &MockSecretsManagerGetSecretAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSecretsManagerGetSecretAPI) EXPECT() *MockSecretsManagerGetSecretAPIMockRecorder {
	return m.recorder
}

// GetSecretValue mocks base method
func (m *MockSecretsManagerGetSecretAPI) GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecretValue", ctx, params)
	ret0, _ := ret[0].(*secretsmanager.GetSecretValueOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecretValue indicates an expected call of GetSecretValue
func (mr *MockSecretsManagerGetSecretAPIMockRecorder) GetSecretValue(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecretValue", reflect.TypeOf((*MockSecretsManagerGetSecretAPI)(nil).GetSecretValue), ctx, params)
}


func TestSecretManagerObject_getSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAPI := NewMockSecretsManagerGetSecretAPI(ctrl)
	ctx := context.TODO()
	logger := logrus.New()

	secretManager := SecretManagerObject{
		SecretName:      "test-secret",
		SecretClientAPI: mockAPI,
		Cxt:             ctx,
		Logger:          logger,
	}

	t.Run("successful getSecret", func(t *testing.T) {
		secretValue := `{"https://example.com": {"Username": "testuser", "Password": "testpass"}}`
		mockAPI.EXPECT().GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId:     aws.String("test-secret"),
			VersionStage: aws.String("AWSCURRENT"),
		}).Return(&secretsmanager.GetSecretValueOutput{
			SecretString: aws.String(secretValue),
		}, nil)

		expectedSecretData := map[string]SecretData{
			"https://example.com": {
				Username: "testuser",
				Password: "testpass",
			},
		}

		secretData, err := secretManager.getSecret()
		assert.NoError(t, err)
		assert.Equal(t, expectedSecretData, secretData)
	})

	t.Run("failed to get secret", func(t *testing.T) {
		mockAPI.EXPECT().GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId:     aws.String("test-secret"),
			VersionStage: aws.String("AWSCURRENT"),
		}).Return(nil, awserr.New("ResourceNotFoundException", "secret not found", nil))

		secretData, err := secretManager.getSecret()
		assert.Error(t, err)
		assert.Nil(t, secretData)
	})

	t.Run("invalid secret format", func(t *testing.T) {
		invalidSecretValue := `{"https://example.com": "invalid_credentials"}`
		mockAPI.EXPECT().GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId:     aws.String("test-secret"),
			VersionStage: aws.String("AWSCURRENT"),
		}).Return(&secretsmanager.GetSecretValueOutput{
			SecretString: aws.String(invalidSecretValue),
		}, nil)

		secretData, err := secretManager.getSecret()
		assert.Error(t, err)
		assert.Nil(t, secretData)
	})
}
func TestSecretManagerObject_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAPI := NewMockSecretsManagerGetSecretAPI(ctrl)
	ctx := context.TODO()
	logger := logrus.New()

	secretManager := SecretManagerObject{
		SecretName:      "test-secret",
		SecretClientAPI: mockAPI,
		Cxt:             ctx,
		Logger:          logger,
	}

	t.Run("unsupported serverURL", func(t *testing.T) {
		serverURL := "https://unsupported.com"
		username, password, err := secretManager.Get(serverURL)
		assert.Error(t, err)
		assert.Equal(t, "", username)
		assert.Equal(t, "", password)
		assert.Equal(t, fmt.Sprintf("unsupported serverURL: %s", serverURL), err.Error())
	})

	t.Run("successful get", func(t *testing.T) {
		serverURL := "https://index.docker.io/v1/"
		secretValue := `{"https://index.docker.io/v1/": {"Username": "testuser", "Password": "testpass"}}`
		mockAPI.EXPECT().GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId:     aws.String("test-secret"),
			VersionStage: aws.String("AWSCURRENT"),
		}).Return(&secretsmanager.GetSecretValueOutput{
			SecretString: aws.String(secretValue),
		}, nil)

		username, password, err := secretManager.Get(serverURL)
		assert.NoError(t, err)
		assert.Equal(t, "testuser", username)
		assert.Equal(t, "testpass", password)
	})

	t.Run("secret not found for serverURL", func(t *testing.T) {
		serverURL := "https://index.docker.io/v1/"
		secretValue := `{"https://example.com": {"Username": "testuser", "Password": "testpass"}}`
		mockAPI.EXPECT().GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId:     aws.String("test-secret"),
			VersionStage: aws.String("AWSCURRENT"),
		}).Return(&secretsmanager.GetSecretValueOutput{
			SecretString: aws.String(secretValue),
		}, nil)

		username, password, err := secretManager.Get(serverURL)
		assert.NoError(t, err)
		assert.Equal(t, "", username)
		assert.Equal(t, "", password)
	})

	t.Run("failed to get secret", func(t *testing.T) {
		serverURL := "https://index.docker.io/v1/"
		mockAPI.EXPECT().GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId:     aws.String("test-secret"),
			VersionStage: aws.String("AWSCURRENT"),
		}).Return(nil, awserr.New("ResourceNotFoundException", "secret not found", nil))

		username, password, err := secretManager.Get(serverURL)
		assert.Error(t, err)
		assert.Equal(t, "", username)
		assert.Equal(t, "", password)
	})
}

