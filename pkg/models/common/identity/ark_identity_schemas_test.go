package identity

import "testing"

func TestGetTenantID(t *testing.T) {
	tests := []struct {
		name     string
		podFqdn  string
		expected string
	}{
		{
			name:     "valid_pod_fqdn",
			podFqdn:  "tenant1.namespace.svc.cluster.local",
			expected: "tenant1",
		},
		{
			name:     "tenant_only",
			podFqdn:  "tenant1",
			expected: "tenant1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tenant := PodFqdnResult{PodFqdn: tt.podFqdn}
			result := tenant.GetTenantID()
			if result != tt.expected {
				t.Errorf("Expected tenant ID '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
