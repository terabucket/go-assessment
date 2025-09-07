package model

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrCertificateNotFound = errors.New("certificate not found")
	ErrClientNotFound      = errors.New("client not found")
)

type Certificate struct {
	ID                    uuid.UUID `json:"certificate_id"`
	CertificatePEMEncoded []byte    `json:"certificate_pem_encoded"`
}

type Client struct {
	ID            uuid.UUID `json:"client_id"`
	CertificateID uuid.UUID `json:"certificate_id"`
}
