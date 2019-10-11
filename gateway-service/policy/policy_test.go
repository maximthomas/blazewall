package policy

import (
	"testing"

	"github.com/maximthomas/blazewall/gateway-service/models"
)

func TestPolicyValidator(t *testing.T) {
	tests := []struct {
		name    string
		p       PolicyValidator
		want    bool
		session *models.Session
	}{
		{
			name: "allowed policy validator test",
			p:    AllowedPolicyValidator{},
			want: true,
		},
		{
			name: "denied policy validator test",
			p:    DeniedPolicyValidator{},
			want: false,
		},
		{
			name: "realms policy positive validator test",
			p: RealmsPolicyValidator{
				Realms: []string{
					"staff",
				},
			},
			want: true,
			session: &models.Session{
				Realm: "staff",
			},
		},

		{
			name: "realms policy negative validator test",
			p: RealmsPolicyValidator{
				Realms: []string{
					"staff",
				},
			},
			want: false,
			session: &models.Session{
				Realm: "users",
			},
		},

		{
			name: "realms policy missing session",
			p: RealmsPolicyValidator{
				Realms: []string{
					"staff",
				},
			},
			want:    false,
			session: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.p.ValidatePolicy(nil, test.session)
			assertPolicyPassed(t, got, test.want)
		})
	}
}

func assertPolicyPassed(t *testing.T, got, want bool) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct policy result, got %t, want %t", got, want)
	}
}
