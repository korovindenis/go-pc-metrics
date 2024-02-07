package ip_test

import (
	"testing"

	"github.com/korovindenis/go-pc-metrics/internal/ip"
	"github.com/stretchr/testify/assert"
)

func TestGetOutbound(t *testing.T) {
	localIP := ip.GetOutbound()
	assert.NotNil(t, localIP, "Local IP should not be nil")
}

func TestCheckInSubnet(t *testing.T) {
	testCases := []struct {
		name    string
		ip      string
		subnet  string
		want    bool
		wantErr bool
	}{
		{
			name:    "IP in subnet",
			ip:      "192.168.1.1",
			subnet:  "192.168.1.0/24",
			want:    true,
			wantErr: false,
		},
		{
			name:    "IP not in subnet",
			ip:      "10.0.0.1",
			subnet:  "192.168.1.0/24",
			want:    false,
			wantErr: false,
		},
		{
			name:    "Invalid subnet",
			ip:      "192.168.1.1",
			subnet:  "invalid_subnet",
			want:    false,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ip.CheckInSubnet(tc.ip, tc.subnet)
			if tc.wantErr {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
				assert.Equal(t, tc.want, got, "Unexpected result")
			}
		})
	}
}
