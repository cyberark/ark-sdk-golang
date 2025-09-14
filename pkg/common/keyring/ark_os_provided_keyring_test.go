package keyring

import (
	"errors"
	"testing"

	"github.com/99designs/keyring"
)

// mockKeyringStore implements keyring.Keyring for testing.
type mockKeyringStore struct {
	setErr    error
	getErr    error
	removeErr error
	item      keyring.Item
}

func (m *mockKeyringStore) Set(item keyring.Item) error {
	if m.setErr != nil {
		return m.setErr
	}
	m.item = item
	return nil
}

func (m *mockKeyringStore) Get(key string) (keyring.Item, error) {
	if m.getErr != nil {
		return keyring.Item{}, m.getErr
	}
	return m.item, nil
}

func (m *mockKeyringStore) Remove(key string) error {
	if m.removeErr != nil {
		return m.removeErr
	}
	m.item = keyring.Item{}
	return nil
}

// The following methods are required to satisfy keyring.Keyring but are unused in these tests.
func (m *mockKeyringStore) Keys() ([]string, error) { return nil, nil }
func (m *mockKeyringStore) GetMetadata(key string) (keyring.Metadata, error) {
	return keyring.Metadata{}, nil
}
func (m *mockKeyringStore) Type() keyring.BackendType { return keyring.KeychainBackend }

func patchOSKeyringOpen(store keyring.Keyring, err error) func() {
	origOpen := OSKeyringOpen
	OSKeyringOpen = func(cfg keyring.Config) (keyring.Keyring, error) {
		return store, err
	}
	return func() {
		OSKeyringOpen = origOpen
	}
}

// mockFallbackKeyring implements ArkKeyringImpl for fallback logic.
type mockFallbackKeyring struct {
	setCalled    bool
	getCalled    bool
	deleteCalled bool
	password     string
	shouldError  bool
}

func (m *mockFallbackKeyring) SetPassword(service, user, password string) error {
	m.setCalled = true
	m.password = password
	if m.shouldError {
		return errors.New("fallback error")
	}
	return nil
}
func (m *mockFallbackKeyring) GetPassword(service, user string) (string, error) {
	m.getCalled = true
	return m.password, nil
}
func (m *mockFallbackKeyring) DeletePassword(service, user string) error {
	m.deleteCalled = true
	m.password = ""
	return nil
}

func (m *mockFallbackKeyring) ClearAllPasswords() error {
	m.deleteCalled = true
	m.password = ""
	return nil
}

func TestArkOSProvidedKeyring_SetPassword(t *testing.T) {
	tests := []struct {
		name             string
		openErr          error
		setErr           error
		fallbackCalled   bool
		fallbackError    error
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name:           "success_case_os_keyring",
			openErr:        nil,
			setErr:         nil,
			fallbackCalled: false,
			expectedError:  false,
		},
		{
			name:           "error_case_open_fallback",
			openErr:        errors.New("open error"),
			setErr:         nil,
			fallbackCalled: true,
			expectedError:  false,
		},
		{
			name:           "error_case_set_fallback",
			openErr:        nil,
			setErr:         errors.New("set error"),
			fallbackCalled: true,
			expectedError:  false,
		},
		{
			name:             "error_case_fallback_returns_error",
			openErr:          errors.New("open error"),
			setErr:           nil,
			fallbackCalled:   true,
			fallbackError:    errors.New("fallback error"),
			expectedError:    true,
			expectedErrorMsg: "fallback error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockKeyringStore{setErr: tt.setErr}
			restore := patchOSKeyringOpen(store, tt.openErr)
			defer restore()

			mockFallback := &mockFallbackKeyring{}
			if tt.fallbackError != nil {
				mockFallback.shouldError = true
			}
			keyring := NewArkOSProvidedKeyring(mockFallback)

			err := keyring.SetPassword("svc", "user", "pwd")
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				if tt.expectedErrorMsg != "" && err.Error() != tt.expectedErrorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.expectedErrorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
			if tt.fallbackCalled && !mockFallback.setCalled {
				t.Errorf("Expected fallback SetPassword to be called")
			}
			if !tt.fallbackCalled && mockFallback.setCalled {
				t.Errorf("Fallback should not be called on OS keyring success")
			}
		})
	}
}

func TestArkOSProvidedKeyring_GetPassword(t *testing.T) {
	tests := []struct {
		name           string
		openErr        error
		getErr         error
		itemData       string
		fallbackCalled bool
		fallbackPwd    string
		expectedPwd    string
	}{
		{
			name:           "success_case_os_keyring",
			openErr:        nil,
			getErr:         nil,
			itemData:       "os_pwd",
			fallbackCalled: false,
			expectedPwd:    "os_pwd",
		},
		{
			name:           "error_case_open_fallback",
			openErr:        errors.New("open error"),
			getErr:         nil,
			itemData:       "",
			fallbackCalled: true,
			fallbackPwd:    "fallback_pwd",
			expectedPwd:    "fallback_pwd",
		},
		{
			name:           "error_case_get_fallback",
			openErr:        nil,
			getErr:         errors.New("get error"),
			itemData:       "",
			fallbackCalled: true,
			fallbackPwd:    "fallback_pwd",
			expectedPwd:    "fallback_pwd",
		},
		{
			name:           "edge_case_key_not_found",
			openErr:        nil,
			getErr:         keyring.ErrKeyNotFound,
			itemData:       "",
			fallbackCalled: false,
			expectedPwd:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockKeyringStore{
				getErr: tt.getErr,
				item:   keyring.Item{Data: []byte(tt.itemData)},
			}
			restore := patchOSKeyringOpen(store, tt.openErr)
			defer restore()

			mockFallback := &mockFallbackKeyring{password: tt.fallbackPwd}
			keyring := NewArkOSProvidedKeyring(mockFallback)

			pwd, err := keyring.GetPassword("svc", "user")
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if pwd != tt.expectedPwd {
				t.Errorf("Expected password %q, got %q", tt.expectedPwd, pwd)
			}
			if tt.fallbackCalled && !mockFallback.getCalled {
				t.Errorf("Expected fallback GetPassword to be called")
			}
			if !tt.fallbackCalled && mockFallback.getCalled {
				t.Errorf("Fallback should not be called on OS keyring success")
			}
		})
	}
}

func TestArkOSProvidedKeyring_DeletePassword(t *testing.T) {
	tests := []struct {
		name           string
		openErr        error
		removeErr      error
		fallbackCalled bool
	}{
		{
			name:           "success_case_os_keyring",
			openErr:        nil,
			removeErr:      nil,
			fallbackCalled: false,
		},
		{
			name:           "error_case_open_fallback",
			openErr:        errors.New("open error"),
			removeErr:      nil,
			fallbackCalled: true,
		},
		{
			name:           "error_case_remove_fallback",
			openErr:        nil,
			removeErr:      errors.New("remove error"),
			fallbackCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockKeyringStore{removeErr: tt.removeErr}
			restore := patchOSKeyringOpen(store, tt.openErr)
			defer restore()

			mockFallback := &mockFallbackKeyring{}
			keyring := NewArkOSProvidedKeyring(mockFallback)

			err := keyring.DeletePassword("svc", "user")
			if (tt.openErr != nil || tt.removeErr != nil) && err != nil {
				if !mockFallback.deleteCalled {
					t.Errorf("Expected fallback DeletePassword to be called")
				}
			}
			if tt.openErr == nil && tt.removeErr == nil && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.openErr == nil && tt.removeErr == nil && mockFallback.deleteCalled {
				t.Errorf("Fallback should not be called on OS keyring success")
			}
		})
	}
}

func TestArkOSProvidedKeyring_FallbackError(t *testing.T) {
	store := &mockKeyringStore{setErr: errors.New("set error")}
	restore := patchOSKeyringOpen(store, nil)
	defer restore()

	mockFallback := &mockFallbackKeyring{}
	mockFallback.shouldError = true
	keyring := NewArkOSProvidedKeyring(mockFallback)

	err := keyring.SetPassword("svc", "user", "pwd")
	if err == nil || err.Error() != "fallback error" {
		t.Errorf("Expected fallback error, got %v", err)
	}
}

func TestArkOSProvidedKeyring_ConcurrentUsage(t *testing.T) {
	store := &mockKeyringStore{}
	restore := patchOSKeyringOpen(store, nil)
	defer restore()

	mockFallback := &mockFallbackKeyring{}
	keyring := NewArkOSProvidedKeyring(mockFallback)

	const n = 10
	done := make(chan bool, n)
	for i := 0; i < n; i++ {
		go func(idx int) {
			_ = keyring.SetPassword("svc", "user", "pwd")
			_, _ = keyring.GetPassword("svc", "user")
			_ = keyring.DeletePassword("svc", "user")
			done <- true
		}(i)
	}
	for i := 0; i < n; i++ {
		<-done
	}
}

type mockKeyringStoreClear struct {
	keysErr   error
	removeErr error
	keys      []string
	removed   []string
}

func (m *mockKeyringStoreClear) Keys() ([]string, error) {
	return m.keys, m.keysErr
}
func (m *mockKeyringStoreClear) Remove(key string) error {
	m.removed = append(m.removed, key)
	return m.removeErr
}
func (m *mockKeyringStoreClear) Set(item keyring.Item) error          { return nil }
func (m *mockKeyringStoreClear) Get(key string) (keyring.Item, error) { return keyring.Item{}, nil }
func (m *mockKeyringStoreClear) GetMetadata(key string) (keyring.Metadata, error) {
	return keyring.Metadata{}, nil
}
func (m *mockKeyringStoreClear) Type() keyring.BackendType { return keyring.KeychainBackend }

func TestArkOSProvidedKeyring_ClearAllPasswords(t *testing.T) {
	tests := []struct {
		name           string
		openErr        error
		keysErr        error
		keys           []string
		removeErr      error
		expectFallback bool
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:           "success_case_os_keyring",
			openErr:        nil,
			keysErr:        nil,
			keys:           []string{"ark_sdk_golang_user1", "other_key"},
			removeErr:      nil,
			expectFallback: false,
			expectedError:  false,
		},
		{
			name:           "error_case_open_fallback",
			openErr:        errors.New("open error"),
			expectFallback: true,
			expectedError:  false,
		},
		{
			name:           "error_case_keys_fallback",
			openErr:        nil,
			keysErr:        errors.New("keys error"),
			expectFallback: true,
			expectedError:  false,
		},
		{
			name:           "edge_case_no_keys_to_remove",
			openErr:        nil,
			keysErr:        nil,
			keys:           []string{"other_key"},
			removeErr:      nil,
			expectFallback: false,
			expectedError:  false,
		},
		{
			name:           "error_case_remove_fails",
			openErr:        nil,
			keysErr:        nil,
			keys:           []string{"ark_sdk_golang_user1"},
			removeErr:      errors.New("remove error"),
			expectFallback: false,
			expectedError:  false, // Current logic ignores remove errors
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockKeyringStoreClear{
				keysErr:   tt.keysErr,
				removeErr: tt.removeErr,
				keys:      tt.keys,
			}
			restore := patchOSKeyringOpen(store, tt.openErr)
			defer restore()

			mockFallback := &mockFallbackKeyring{}
			keyring := NewArkOSProvidedKeyring(mockFallback)

			err := keyring.ClearAllPasswords()
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				if tt.expectedErrMsg != "" && err.Error() != tt.expectedErrMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.expectedErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
			if tt.expectFallback && !mockFallback.deleteCalled {
				t.Errorf("Expected fallback ClearAllPasswords to be called")
			}
			if !tt.expectFallback && mockFallback.deleteCalled {
				t.Errorf("Fallback should not be called on OS keyring success")
			}
			// Validate that only ARK keys are removed
			if tt.name == "success_case_os_keyring" && len(store.removed) != 1 {
				t.Errorf("Expected 1 ARK key removed, got %d", len(store.removed))
			}
		})
	}
}
