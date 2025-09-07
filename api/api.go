package api

import (
	"context"
	"errors"
	"math/big"
	"net/http"

	"crypto/x509"
	"encoding/pem"
	"encoding/json"
	"strings"
	"fmt"
	"task/model"

	"github.com/google/uuid"

	openapi_types "github.com/oapi-codegen/runtime/types"
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
	// Fetch the current certicates (stored in a API somewhere)
	ctx := r.Context()
	// Step 1: get all certs from database
	certs, err := a.db.ListCertificates(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Step 2: convert model certs into API certs
	apiCerts := []Certificate{}
	for _, c := range certs {
		apiCert, err := convertModelCertToAPICert(ctx, c)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid certificate: %v", err), http.StatusInternalServerError)
			return 
		}
		apiCerts = append(apiCerts, apiCert)
	}

	// Step 3: Encode as JSON and send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(ListCertificatesResponse{
		Certificates: apiCerts,
	})
}

// (POST /certificates)
func (a *API) PostCertificates(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement me
	ctx := r.Context()

	// Step 1: decode request body
	var req UpdateCertificateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid req body", http.StatusBadRequest)
		return
	}

	// Step 2: ensure client exists
	client, err := a.db.GetClient(ctx, req.ClientId)
	if err != nil {
		http.Error(w, "client not found", http.StatusNotFound)
		return
	}
	// Step 3: ensure certificate exists
	cert, err := a.db.GetCertificate(ctx, req.CertificateId)
	if err != nil {
		http.Error(w, "certificate not found", http.StatusNotFound)
		return
	}
	// Step 4: If there is no change, return conflict
	if client.CertificateID == req.CertificateId {
		http.Error(w, "there is no certificate change", http.StatusConflict)
		return
	}
	// Step 5: Update client cert
	if err := a.db.UpdateClientCertificate(ctx, req.ClientId, req.CertificateId); err != nil {
		http.Error(w, "failed to update", http.StatusInternalServerError)
		return
	}
	// Step 6: Convert model cert into API cert
	apiCert, err := convertModelCertToAPICert(ctx, cert)
	if err != nil {
		return
	}
	// Step 7: Send success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(UpdateCertificateResponse{
		ClientId: req.ClientId,
		Certificate: apiCert,
	})
}

func FormatSerialNumber(serial *big.Int) string {
	// TODO: Implement me (bonus task)
	hexStr := serial.Text(16)
	if len(hexStr)%2 != 0{
		hexStr = "0" + hexStr
	}

	//hexStr = strings.ToUpper(hexStr)

	var parts []string
	for i := 0; i < len(hexStr); i += 2 {
		parts = append(parts, hexStr[i:i+2])
	}

	return strings.Join(parts, ":")
}

func convertModelCertToAPICert(ctx context.Context, mc *model.Certificate) (Certificate, error) {
	block, _ := pem.Decode(mc.CertificatePEMEncoded)
	if block == nil {
		return Certificate{}, fmt.Errorf("failed to decode PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return Certificate{}, err
	}

	return Certificate{
		Id: mc.ID,
		SerialNumber: FormatSerialNumber(cert.SerialNumber),
		NotBefore: openapi_types.Date{Time: cert.NotBefore},
		NotAfter: openapi_types.Date{Time: cert.NotAfter},
	}, nil
}