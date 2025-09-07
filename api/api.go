package api

import (
	"context"
	"errors"
	"math/big"
	"net/http"

	"task/model"

	"github.com/google/uuid"
)

var ErrNoCertificateChange = errors.New("current certificate matches existing certificate")

type Store interface {
	GetClient(ctx context.Context, clientID uuid.UUID) (*model.Client, error)
	GetCertificate(ctx context.Context, certificateID uuid.UUID) (*model.Certificate, error)
	ListCertificates(ctx context.Context) ([]*model.Certificate, error)
	UpdateClientCertificate(ctx context.Context, clientID, certificateID uuid.UUID) error
}

type API struct {
	ServerInterface
	db Store
}

func New(db Store) *API {
	return &API{
		db: db,
	}
}

// (GET /certificates)
func (a *API) GetCertificates(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement me
}

// (POST /certificates)
func (a *API) PostCertificates(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement me
}

func FormatSerialNumber(serial *big.Int) string {
	// TODO: Implement me (bonus task)
	return ""
}
