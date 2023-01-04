package v3

import (
	"context"
	"testing"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/gardener/ext-authz-server/pkg/auth"
	"google.golang.org/genproto/googleapis/rpc/code"
)

const service auth.Services = `^outbound\|1194\|\|vpn-seed-server(-[0-4])?\..*\.svc\.cluster\.local$`

func TestCheck(t *testing.T) {

	type testData struct {
		headerValue    string
		expectedStatus code.Code
	}

	tests := []testData{
		{"outbound|1194||vpn-seed-server.foo.svc.cluster.local", code.Code_OK},
		{"outbound|1194||vpn-seed-server-0.foo.svc.cluster.local", code.Code_OK},
		{"outbound|1194||vpn-seed-server-2.foo.svc.cluster.local", code.Code_OK},
		{"outbound|1194||vpn-seed-server-11.foo.svc.cluster.local", code.Code_PERMISSION_DENIED},
		{"outbound|443||kube-apiserver.foo.svc.cluster.local", code.Code_PERMISSION_DENIED},
		{"", code.Code_PERMISSION_DENIED},
	}

	server := New(service)

	ctx := context.Background()

	for _, test := range tests {
		req := &envoy_service_auth_v3.CheckRequest{
			Attributes: &envoy_service_auth_v3.AttributeContext{
				Request: &envoy_service_auth_v3.AttributeContext_Request{
					Http: &envoy_service_auth_v3.AttributeContext_HttpRequest{
						Headers: map[string]string{
							"reversed-vpn": test.headerValue,
						},
					},
				},
			},
		}

		resp, err := server.Check(ctx, req)

		if err != nil {
			t.Errorf("Check returned an error: %v", err)
		}

		if code.Code(resp.Status.Code) != test.expectedStatus {
			t.Errorf("Unexpected response status: %v", resp.Status.Code)
		}
	}
}
