// the business logic of the backend

package serverusecase

import (
	"context"
	"testing"
)

func TestServer_Ping(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		s       *Server
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Ping(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Server.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
