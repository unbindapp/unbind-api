package k8s

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/services/models"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// RegistryCredential represents credentials for a single registry
type RegistryCredential struct {
	RegistryURL string
	Username    string
	Password    string
}

// CreateMultiRegistryCredentials creates or updates a kubernetes.io/dockerconfigjson secret for multiple container registries
func (self *KubeClient) CreateMultiRegistryCredentials(ctx context.Context, name, namespace string,
	credentials []RegistryCredential, client *kubernetes.Clientset) (*corev1.Secret, error) {

	type DockerConfigEntry struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
		Auth     string `json:"auth,omitempty"`
	}

	type DockerConfigJSON struct {
		Auths map[string]DockerConfigEntry `json:"auths"`
	}

	// Initialize docker config with empty auths map
	dockerConfig := DockerConfigJSON{
		Auths: make(map[string]DockerConfigEntry),
	}

	// Add each registry credential to the config
	for _, cred := range credentials {
		// Create auth string (base64 encoded username:password)
		auth := base64.StdEncoding.EncodeToString([]byte(cred.Username + ":" + cred.Password))

		// Add to docker config
		dockerConfig.Auths[cred.RegistryURL] = DockerConfigEntry{
			Username: cred.Username,
			Password: cred.Password,
			Auth:     auth,
		}
	}

	dockerConfigJSON, err := json.Marshal(dockerConfig)
	if err != nil {
		return nil, err
	}

	// Check if secret already exists
	secret, err := self.GetSecret(ctx, name, namespace, client)
	if err != nil && !apierrors.IsNotFound(err) {
		return nil, err
	}

	// Create or update the secret
	if apierrors.IsNotFound(err) {
		// Create new secret
		newSecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Type: corev1.SecretTypeDockerConfigJson,
			Data: map[string][]byte{
				corev1.DockerConfigJsonKey: dockerConfigJSON,
			},
		}
		return client.CoreV1().Secrets(namespace).Create(ctx, newSecret, metav1.CreateOptions{})
	} else {
		// Update existing secret
		secret.Type = corev1.SecretTypeDockerConfigJson
		if secret.Data == nil {
			secret.Data = make(map[string][]byte)
		}
		secret.Data[corev1.DockerConfigJsonKey] = dockerConfigJSON
		return client.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
	}
}

// GetOrCreateSecret retrieves an existing secret or creates a new one if it doesn't exist
// Returns the secret and a boolean indicating if it was created (true) or retrieved (false)
func (self *KubeClient) GetOrCreateSecret(ctx context.Context, name, namespace string, client *kubernetes.Clientset) (*corev1.Secret, bool, error) {
	// Try to get the secret
	secret, err := client.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err == nil {
		// Secret exists, return it
		return secret, false, nil
	}

	// Create the secret if it doesn't exist
	newSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Type: corev1.SecretTypeOpaque,
	}

	createdSecret, err := client.CoreV1().Secrets(namespace).Create(ctx, newSecret, metav1.CreateOptions{})
	if err != nil {
		return nil, false, err
	}

	return createdSecret, true, nil
}

// GetSecret retrieves a secret by name in the given namespace
func (self *KubeClient) GetSecret(ctx context.Context, name, namespace string, client *kubernetes.Clientset) (*corev1.Secret, error) {
	return client.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
}

// UpdateSecret updates an existing secret with new data
func (self *KubeClient) UpdateSecret(ctx context.Context, name, namespace string, data map[string][]byte, client *kubernetes.Clientset) (*corev1.Secret, error) {
	// Get the current secret
	secret, err := self.GetSecret(ctx, name, namespace, client)
	if err != nil {
		return nil, err
	}

	// Update the data
	secret.Data = data

	return client.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
}

// GetSecretValue retrieves a specific key from a secret
func (self *KubeClient) GetSecretValue(ctx context.Context, name, namespace, key string, client *kubernetes.Clientset) ([]byte, error) {
	secret, err := self.GetSecret(ctx, name, namespace, client)
	if err != nil {
		return nil, err
	}

	value, exists := secret.Data[key]
	if !exists {
		return nil, errors.New("key does not exist in secret")
	}

	return value, nil
}

// GetSecretMap retrieves all key-value pairs from a secret as a map
func (self *KubeClient) GetSecretMap(ctx context.Context, name, namespace string, client *kubernetes.Clientset) (map[string][]byte, error) {
	secret, err := self.GetSecret(ctx, name, namespace, client)
	if err != nil {
		return nil, err
	}

	return secret.Data, nil
}

// DeleteSecret deletes a secret by name in the given namespace
func (self *KubeClient) DeleteSecret(ctx context.Context, name, namespace string, client *kubernetes.Clientset) error {
	err := client.CoreV1().Secrets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	return nil
}

// UpsertSecretValues adds or updates specific keys in a secret without affecting other keys
func (self *KubeClient) UpsertSecretValues(ctx context.Context, name, namespace string, values map[string][]byte, client *kubernetes.Clientset) (*corev1.Secret, error) {
	// Get the current secret
	secret, err := self.GetSecret(ctx, name, namespace, client)
	if err != nil {
		return nil, err
	}

	// If the data map is nil, initialize it
	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}

	// Update the values
	for k, v := range values {
		secret.Data[k] = v
	}

	return client.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
}

// OverwriteSecretValues overwrites all values in a secret with new values
func (self *KubeClient) OverwriteSecretValues(ctx context.Context, name, namespace string, values map[string][]byte, client *kubernetes.Clientset) (*corev1.Secret, error) {
	// Get the current secret
	secret, err := self.GetSecret(ctx, name, namespace, client)
	if err != nil {
		return nil, err
	}

	// Reset values
	secret.Data = make(map[string][]byte)

	// Set the values
	for k, v := range values {
		secret.Data[k] = v
	}

	return client.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
}

// GetAllSecrets retrieves all secrets for the team hierarchy and returns them with their data
func (self *KubeClient) GetAllSecrets(
	ctx context.Context,
	teamID uuid.UUID,
	teamSecret string,
	projectSecrets map[uuid.UUID]string,
	environmentSecrets map[uuid.UUID]string,
	serviceSecrets map[uuid.UUID]string,
	client *kubernetes.Clientset,
	namespace string,
) ([]models.SecretData, error) {
	var result []models.SecretData

	// Process team secret
	if teamSecret != "" {
		secretData, err := self.processSecret(ctx, teamID, schema.VariableReferenceSourceTypeTeam, teamSecret, client, namespace)
		if err != nil {
			return nil, err
		}
		result = append(result, secretData)
	}

	// Process project secrets
	for projectID, secretName := range projectSecrets {
		if secretName == "" {
			continue
		}
		secretData, err := self.processSecret(ctx, projectID, schema.VariableReferenceSourceTypeProject, secretName, client, namespace)
		if err != nil {
			return nil, err
		}
		result = append(result, secretData)
	}

	// Process environment secrets
	for envID, secretName := range environmentSecrets {
		if secretName == "" {
			continue
		}
		secretData, err := self.processSecret(ctx, envID, schema.VariableReferenceSourceTypeEnvironment, secretName, client, namespace)
		if err != nil {
			return nil, err
		}
		result = append(result, secretData)
	}

	// Process service secrets
	for serviceID, secretName := range serviceSecrets {
		if secretName == "" {
			continue
		}
		secretData, err := self.processSecret(ctx, serviceID, schema.VariableReferenceSourceTypeService, secretName, client, namespace)
		if err != nil {
			return nil, err
		}
		result = append(result, secretData)
	}

	return result, nil
}

// Helper function to process a single secret
func (self *KubeClient) processSecret(
	ctx context.Context,
	id uuid.UUID,
	secretType schema.VariableReferenceSourceType,
	secretName string,
	client *kubernetes.Clientset,
	namespace string,
) (models.SecretData, error) {
	// Get the secret data
	secretMap, err := self.GetSecretMap(ctx, secretName, namespace, client)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// If secret doesn't exist, return an empty data map
			return models.SecretData{
				ID:         id,
				Type:       secretType,
				SecretName: secretName,
				Data:       make(map[string][]byte),
			}, nil
		}
		return models.SecretData{}, err
	}

	// Return the secret with its data
	return models.SecretData{
		ID:         id,
		Type:       secretType,
		SecretName: secretName,
		Data:       secretMap,
	}, nil
}
