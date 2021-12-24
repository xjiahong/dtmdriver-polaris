package driver

import (
	"testing"
)

func Test_polarisDriver_ParseServerMethod(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name       string
		args       args
		wantServer string
		wantMethod string
		wantErr    bool
	}{
		{
			name:       "ip:port",
			args:       args{uri: "127.0.0.1:8080/package.service/method"},
			wantServer: "127.0.0.1:8080",
			wantMethod: "/package.service/method",
		},
		{
			name:       "polaris",
			args:       args{uri: "polaris://service/package.service/method?namespace=Test"},
			wantServer: "polaris:///service?namespace=Test",
			wantMethod: "/package.service/method",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &polarisDriver{}
			gotServer, gotMethod, err := p.ParseServerMethod(tt.args.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseServerMethod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotServer != tt.wantServer {
				t.Errorf("ParseServerMethod() gotServer = %v, want %v", gotServer, tt.wantServer)
			}
			if gotMethod != tt.wantMethod {
				t.Errorf("ParseServerMethod() gotMethod = %v, want %v", gotMethod, tt.wantMethod)
			}
		})
	}
}
