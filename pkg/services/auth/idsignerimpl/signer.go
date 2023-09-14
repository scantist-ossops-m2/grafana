package idsignerimpl

import (
	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"

	"github.com/grafana/grafana/pkg/services/auth"
	"github.com/grafana/grafana/pkg/services/signingkeys"
)

type localSigner struct {
	signer jose.Signer
}

func newLocalSigner(signingKeyService signingkeys.Service) (*localSigner, error) {
	key := signingKeyService.GetServerPrivateKey() // FIXME: replace with signing specific key

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.ES256, Key: key}, &jose.SignerOptions{
		ExtraHeaders: map[jose.HeaderKey]interface{}{
			"kid": "default", // FIXME: replace with specific key id
		},
	})
	if err != nil {
		return nil, err
	}

	return &localSigner{
		signer: signer,
	}, nil
}

func (t *localSigner) SignToken(claims *jwt.Claims, assertions *auth.IDAssertions) (string, error) {
	builder := jwt.Signed(t.signer).Claims(claims).Claims(assertions)

	// Build the JWT.
	token, err := builder.CompactSerialize()
	if err != nil {
		return "", err
	}

	return token, nil
}
