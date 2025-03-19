package k8s

import (
	"context"
	"errors"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetOrCreateSecret retrieves an existing secret or creates a new one if it doesn't exist
// Returns the secret and a boolean indicating if it was created (true) or retrieved (false)
func (self *KubeClient) GetOrCreateSecret(ctx context.Context, name, namespace string) (*corev1.Secret, bool, error) {
	// Try to get the secret
	secret, err := self.clientset.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
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

	createdSecret, err := self.clientset.CoreV1().Secrets(namespace).Create(ctx, newSecret, metav1.CreateOptions{})
	if err != nil {
		return nil, false, err
	}

	return createdSecret, true, nil
}

// GetSecret retrieves a secret by name in the given namespace
func (self *KubeClient) GetSecret(ctx context.Context, name, namespace string) (*corev1.Secret, error) {
	return self.clientset.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
}

// UpdateSecret updates an existing secret with new data
func (self *KubeClient) UpdateSecret(ctx context.Context, name, namespace string, data map[string][]byte) (*corev1.Secret, error) {
	// Get the current secret
	secret, err := self.GetSecret(ctx, name, namespace)
	if err != nil {
		return nil, err
	}

	// Update the data
	secret.Data = data

	return self.clientset.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
}

// GetSecretValue retrieves a specific key from a secret
func (self *KubeClient) GetSecretValue(ctx context.Context, name, namespace, key string) ([]byte, error) {
	secret, err := self.GetSecret(ctx, name, namespace)
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
func (self *KubeClient) GetSecretMap(ctx context.Context, name, namespace string) (map[string][]byte, error) {
	secret, err := self.GetSecret(ctx, name, namespace)
	if err != nil {
		return nil, err
	}

	return secret.Data, nil
}

// DeleteSecret deletes a secret by name in the given namespace
func (self *KubeClient) DeleteSecret(ctx context.Context, name, namespace string) error {
	err := self.clientset.CoreV1().Secrets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	return nil
}

// AddSecretValues adds or updates specific keys in a secret without affecting other keys
func (self *KubeClient) AddSecretValues(ctx context.Context, name, namespace string, values map[string][]byte) (*corev1.Secret, error) {
	// Get the current secret
	secret, err := self.GetSecret(ctx, name, namespace)
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

	return self.clientset.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
}
