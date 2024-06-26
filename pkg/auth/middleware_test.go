package auth

import (
	"net/http"
	"reflect"
	"testing"

	jwt "github.com/golang-jwt/jwt/v4"
	jwtReq "github.com/golang-jwt/jwt/v4/request"
)

func Test_requestContainsToken(t *testing.T) {
	tests := []struct {
		name    string
		r       *http.Request
		headers []string
		want    bool
	}{
		{
			name: "Get token using default Authorization header",
			r: &http.Request{
				Header: map[string][]string{
					"Authorization": {"Bearer: eyJhbGciOiJIUzI1NiJ9.eyJ0ZW5hbnRfaWQiOiIxMjMiLCJ2ZXJzaW9uIjoiMSJ9.QmTSzJbvlB5_QmSmYb3nrpnUK4xuK9iWACc5xl8mmLU"},
				},
			},
			headers: []string{"Authorization"},
			want:    true,
		},
		{
			name: "Get token from X-Id-Token header first then Authorization (both present)",
			r: &http.Request{
				Header: map[string][]string{
					"Authorization": {"Bearer: eyJhbGciOiJIUzI1NiJ9.eyJ0ZW5hbnRfaWQiOiIxMjMiLCJ2ZXJzaW9uIjoiMSJ9.QmTSzJbvlB5_QmSmYb3nrpnUK4xuK9iWACc5xl8mmLU"},
					"X-Id-Token":    {"eyJhbGciOiJIUzI1NiJ9.eyJ0ZW5hbnRfaWQiOiIxMjMiLCJ2ZXJzaW9uIjoiMSJ9.QmTSzJbvlB5_QmSmYb3nrpnUK4xuK9iWACc5xl8mmLU"},
				},
			},
			headers: []string{"X-Id-Token", "Authorization"},
			want:    true,
		},
		{
			name: "Get token from X-Id-Token header first then Authorization (X-Id-Token present)",
			r: &http.Request{
				Header: map[string][]string{
					"X-Id-Token": {"eyJhbGciOiJIUzI1NiJ9.eyJ0ZW5hbnRfaWQiOiIxMjMiLCJ2ZXJzaW9uIjoiMSJ9.QmTSzJbvlB5_QmSmYb3nrpnUK4xuK9iWACc5xl8mmLU"},
				},
			},
			headers: []string{"X-Id-Token", "Authorization"},
			want:    true,
		},
		{
			name: "Get token using lowercase headers",
			r: &http.Request{
				Header: map[string][]string{
					"Authorization": {"Bearer: eyJhbGciOiJIUzI1NiJ9.eyJ0ZW5hbnRfaWQiOiIxMjMiLCJ2ZXJzaW9uIjoiMSJ9.QmTSzJbvlB5_QmSmYb3nrpnUK4xuK9iWACc5xl8mmLU"},
					"X-Id-Token":    {"eyJhbGciOiJIUzI1NiJ9.eyJ0ZW5hbnRfaWQiOiIxMjMiLCJ2ZXJzaW9uIjoiMSJ9.QmTSzJbvlB5_QmSmYb3nrpnUK4xuK9iWACc5xl8mmLU"},
				},
			},
			headers: []string{"x-id-token", "authorization"},
			want:    true,
		},
		{
			name:    "Missing Authorization or X-Id-Token headers",
			r:       &http.Request{},
			headers: []string{"X-Id-Token", "Authorization"},
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := requestContainsToken(tt.r, tt.headers); got != tt.want {
				t.Errorf("requestContainsToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildHeaderExtractor(t *testing.T) {
	tests := []struct {
		name         string
		extraHeaders []string
		want         jwtReq.Extractor
	}{
		{
			name: "Default",
			want: jwtReq.MultiExtractor{
				jwtReq.AuthorizationHeaderExtractor,
			},
		},
		{
			name:         "One extra header",
			extraHeaders: []string{"X-Id-Token"},
			want: jwtReq.MultiExtractor{
				jwtReq.HeaderExtractor{"X-Id-Token"},
				jwtReq.AuthorizationHeaderExtractor,
			},
		},
		{
			name:         "Two extra headers",
			extraHeaders: []string{"X-Id-Token", "JWT-ID"},
			want: jwtReq.MultiExtractor{
				jwtReq.HeaderExtractor{"X-Id-Token"},
				jwtReq.HeaderExtractor{"JWT-ID"},
				jwtReq.AuthorizationHeaderExtractor,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildHeaderExtractor(tt.extraHeaders); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildHeaderExtractor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractTenantID(t *testing.T) {
	claim := jwt.MapClaims{
		"tenant_id":                     "1234",
		"https://example.com/tenant_id": "1234",
		"empty":                         "",
	}
	tests := []struct {
		name          string
		claim         jwt.MapClaims
		tenantIDClaim string
		want          string
		wantErr       bool
	}{
		{
			name:          "Tenant ID Claim",
			claim:         claim,
			tenantIDClaim: "tenant_id",
			want:          "1234",
			wantErr:       false,
		},
		{
			name:          "URL format Claim",
			claim:         claim,
			tenantIDClaim: "https://example.com/tenant_id",
			want:          "1234",
			wantErr:       false,
		},
		{
			name:          "Empty tenant id",
			claim:         claim,
			tenantIDClaim: "empty",
			want:          "",
			wantErr:       true,
		},
		{
			name:          "Missing Claim",
			claim:         claim,
			tenantIDClaim: "absent",
			want:          "",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractTenantID(tt.claim, tt.tenantIDClaim)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractTenantID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractTenantID() = %v, want %v", got, tt.want)
			}
		})
	}
}
