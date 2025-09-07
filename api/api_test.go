package api_test

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"slices"
	"strings"
	"testing"
	"time"

	"task/api"
	"task/model"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
	"github.com/oapi-codegen/testutil"
	assert "github.com/stretchr/testify/require"
)

const (
	cert27bfe9fc = `-----BEGIN CERTIFICATE-----
MIIFqzCCA5OgAwIBAgIUWxn/c6jZ0lbN/osHzynrW01TtjAwDQYJKoZIhvcNAQEL
BQAwZTELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAklMMRAwDgYDVQQHDAdDaGljYWdv
MQ0wCwYDVQQKDARBY21lMRIwEAYDVQQLDAlBY21lQ2VydHMxFDASBgNVBAMMC2V4
YW1wbGUuY29tMB4XDTI1MDUxNjE1MzQ0OVoXDTI2MDcyNTE1MzQ0OVowZTELMAkG
A1UEBhMCVVMxCzAJBgNVBAgMAklMMRAwDgYDVQQHDAdDaGljYWdvMQ0wCwYDVQQK
DARBY21lMRIwEAYDVQQLDAlBY21lQ2VydHMxFDASBgNVBAMMC2V4YW1wbGUuY29t
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEArhRjRve4PyHMlXdVFk1m
7xMTLeoz089k+xxrF/ABoeUnVDpChE/psUc24VNaK0XDfK716wOtPqqKXe0jsb/4
3YSABlJc4J55r8Z8O6tqLUeU3P+eg1P0n4VMyV3V8dD/kuDL6jLVM5Ife0ZofWSc
pcrQDSkpByzA9JPI5w+IZwJDMVkV2Xb8TUY5DzORDa0ugH+LawymZNk2YQeY5Vh2
/KHwssca9fqFlPms6DPuLGKcdV5Z4FKfffs7T0GvJkmPqpAgybdwoyJ0RNqVZL1M
Q9qPJR4Fj6Q0i9b/owhqDSXOJvrj0ecYub+63EO8kCBkCoZrEssQvsBVw+4ct/sm
3BPA183Mlzv0Wq9PbXcj4LHvwJjlHjtwook6vdZ12JBs6grSxHbusX5c4IMdY+jO
xOwcCamPkkuQeFb4catoZzB8zNxX4RQB6LsB6BCqJjspqTUO+Iz4L4HDQWhMjRen
ijKczhcq11PzhUiBjvUVWHjD5DECIVO3VB9z9OA5RZNWuJEyU/3Yiq4qGgvs+wqy
6KpDNcaLATFTw/wKocqEiPPjK1bnQ5KMvgljv7Vt3Jst8WiqSUwYFi3YnOCNNpbk
exXao9ySUBB79FPN5DpDVWGZUAZYXI11tsUAbeWepgSlVtGpu8CbKplwbq7rtRlI
mvNTML2D0eRDeQzTUI7FK7sCAwEAAaNTMFEwHQYDVR0OBBYEFIF6C8+eTzw4I3mb
iPawd96Ol7I7MB8GA1UdIwQYMBaAFIF6C8+eTzw4I3mbiPawd96Ol7I7MA8GA1Ud
EwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggIBAIPleVRV27teq2yDE2OFLUI2
hiiQbEI6z6E0u/+gGeHBNGkM0/ovT96Vf62r9C30HmUt+6gMGI1Hvp64GQ2/Umf8
AFGd2fLr8bjSPSih61uz5Ijh/vODxAhdV3lKncvAVeQFJK2pVkJjE7JvQTpfIJj+
YvaXts0b5plRIYwbBCeV06vC95B09lbSQKGxSGQXP4eh5pJrYdSy60svEAXaxyRN
GKsi8ZuFX3hoCJc7SIxcL2O0BIrsiMplsQKx1l1t5j1zvDM9yv0LuASdR9NL3Bu9
Xf8qHz4pbV+UpL6dM4qjPudX9C4Jq00FJ8Z0meSPDmYa6FqkLiJktb88XlX9ZvI3
NqPuQWdngvthf1OMyfk4tJqmdLg8Jd90VTuZ2W9IY1JuDQ9xGoOgBXqt2tl2cRNJ
UFd7RjkJaueB/GNtyQPTcAmPWxPdFh5WsRwEREFCEwzNsSI5bDn3mKtaJrkeBE67
XSasQgWaVo3+JShkAoMJ6BzhdchmHY3wYEC5gy1ltsW7wB7WxxMVqsxkFITrr1mt
4dU99RSxxV19ogtUcwmWK3l2E1HzTWgwtFt8Aw+4dXtp8lfj+BaA/Pdv03Sn/dZx
xaKDVfj+iRnJ0XFttG3DA1rHwQHte7Cj+32pt5w4HYgI9BDhyOpNkXJH7Iw4wI1d
pA0G4ZYLe/UpdWGv5Qlp
-----END CERTIFICATE-----`

	certf0e5137f = `-----BEGIN CERTIFICATE-----
MIIFqzCCA5OgAwIBAgIUNTz/+9GEbqt9gjzfn05HUoHJLZ0wDQYJKoZIhvcNAQEL
BQAwZTELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAklMMRAwDgYDVQQHDAdDaGljYWdv
MQ0wCwYDVQQKDARBY21lMRIwEAYDVQQLDAlBY21lQ2VydHMxFDASBgNVBAMMC2V4
YW1wbGUuY29tMB4XDTI1MDUxNjE1MzcxNFoXDTI2MTIyNTE1MzcxNFowZTELMAkG
A1UEBhMCVVMxCzAJBgNVBAgMAklMMRAwDgYDVQQHDAdDaGljYWdvMQ0wCwYDVQQK
DARBY21lMRIwEAYDVQQLDAlBY21lQ2VydHMxFDASBgNVBAMMC2V4YW1wbGUuY29t
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAqthC7otfWmDnQHYhrDnN
ofgA5Z+spGhel9PPrvb0yoRnqJTOCey6WfgOddczir2a843N5a9JnEsdevJdMfwA
KWbkCBUcBVmLzxPbWfVywTWjw9TmZ7pJVEyEuNyPlIPzQ7X4OG0HCVGe9hYG9lmk
iZvuKR8Pe0PMu8Xw5/0pMupwSj9OehPSFQaqnYtG7GbbyQKV9ywEkPXEPX/jmv42
MPKp3zwCJvMHi/2thPQ1Hel9Qrw0/osvApNJCSoPyPW7ejJo6LFYDQgoHSF/6MWs
yCUEEwjOw8hFReiV4UaPU6zgaLqI2byYHwelld/ZPP9EVviHTvvvOWYEjm2PfFww
kLih78lkV7QUyQVdxw0Na1sN864BExzbd+eVAoJsmDdfZuNt/gTAb5fEdWXeyAJz
eirf1UkRByUGW9+pn2eeWt90UG0xXc2WhYGoFYoIF6UzfFJA/CCVqqljt/aqPWz3
4YUiK9OHjgjdPRsvGP06GNg4LJKhNTcJJb7U0HifD+4FIQL7lhEEc/VU7Q81W4Lv
CyDs+ZDmEI6SMo9jSlNqPP68InAS2CiElrnNavBb3L46Pd42k+uG/1oMmBRSLo+6
wfhIuIqiLqr0cIHkcVaFAnEMHlEjwgs5kumsq4fB4lWMF056nIXs6FMsuu6C0VEG
DPaGhwYcPJ7Bz6DUOEQkMJ8CAwEAAaNTMFEwHQYDVR0OBBYEFFBiKyzbIfEyYUYo
GiXiMGq0YeBnMB8GA1UdIwQYMBaAFFBiKyzbIfEyYUYoGiXiMGq0YeBnMA8GA1Ud
EwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggIBAC/j5kdXabsrWT0JedPkoW23
DKQ6tLqDeo5+70MSH+RQBrBOzmp3GQYaotYllHCJVCu3MBbuiYNKF5JaBUNe+IPM
zAXiCyMiumuTIn8R9NpBzJDE+Mb/dIcDi9N3NZS87R8f+Q2zHlmOr7QeKiahPo/I
XpJMpRpBRk0x3L4L5AxyVNbBIN1mdY/yL1Tz8WSUrls3xcCBBiNa24BEp0OXpyfV
EOgu3NrEJN6K6ayamfDjFHJzXRuZXsrwOd1YGoeg75UQW+qTCSLReuazbwq+U6O3
7qlp8pHjA70YBw1r+zwyDMlAMeU9tnnm2IWElOygTzUrWqtfuoUJjXjt9pZl1LW4
Cup8wPK5hK4RF7acIfHZiMKaaMMwe7f2QrAXB6l1uaXHtmDuyQAD5l8Jcg5sYPvz
XnWax65Lib1QmWvbp3anCsljcnZ7FIwBhsw6Z72DbgOqMNejYH69Jy7olbyEtuHR
EHRIJua6cK+2FLEH99vAE3brJet+nlbtqqJsebI1Okn/rAuW6Eq1vVhFtVgmiCmm
zD8LFRNy41DrVNAzIjSfhpCZkKq9jE+z+5pO1aQam+olwJcBr2snPkhGIZVrF5x6
Q85tGcim1r4dbqzpwbWESx+nNNIqiL4uPUgbshk80Lfa9dpX6AjAReYmFOrUVZyD
8v0L177Ltm4aaa46Liu1
-----END CERTIFICATE-----`
)

//go:embed certificates.json
var certificatesJSON []byte

//go:embed clients.json
var clientsJSON []byte

func TestGetCertificates(t *testing.T) {
	tests := []struct {
		name     string
		req      *testutil.RequestBuilder
		status   int
		expected api.ListCertificatesResponse
	}{
		{
			name: "lists all certificates",
			req: testutil.
				NewRequest().
				Get("/api/v1/certificates").
				WithAcceptJson(),
			status: http.StatusOK,
			expected: api.ListCertificatesResponse{
				[]api.Certificate{
					{
						Id:                    uuid.MustParse("27bfe9fc-b80b-46a1-a967-78658db0aeec"),
						NotBefore:             types.Date{Time: time.Date(2025, time.May, 16, 0, 0, 0, 0, time.UTC)},
						NotAfter:              types.Date{Time: time.Date(2026, time.July, 25, 0, 0, 0, 0, time.UTC)},
						SerialNumber:          "520097931764758754948131491264192108474121041456",
						CertificatePemEncoded: cert27bfe9fc,
					},
					{
						Id:                    uuid.MustParse("f0e5137f-03e1-4ca9-8dd9-b79da983d6be"),
						NotBefore:             types.Date{Time: time.Date(2025, time.May, 16, 0, 0, 0, 0, time.UTC)},
						NotAfter:              types.Date{Time: time.Date(2026, time.December, 25, 0, 0, 0, 0, time.UTC)},
						SerialNumber:          "303936854887858307147587267379171672533028646301",
						CertificatePemEncoded: certf0e5137f,
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := http.NewServeMux()

			opts := api.StdHTTPServerOptions{
				BaseURL:    "/api/v1",
				BaseRouter: m,
			}

			api.HandlerWithOptions(api.New(newStore(t)), opts)

			rr := test.req.GoWithHTTPHandler(t, m).Recorder

			assert.Equal(t, test.status, rr.Code, rr.Body.String())

			var actual api.ListCertificatesResponse
			assert.NoError(t, json.NewDecoder(rr.Body).Decode(&actual))

			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestPostCertificates(t *testing.T) {
	tests := []struct {
		name    string
		req     *testutil.RequestBuilder
		status  int
		failed  *api.Error
		success *api.UpdateCertificateResponse
	}{
		{
			name: "updates the client's certificate successfully",
			req: testutil.
				NewRequest().
				Post("/api/v1/certificates").
				WithJsonBody(api.PostCertificatesJSONRequestBody{
					ClientId:      uuid.MustParse("63838416-B316-418A-8CC8-9EFE3411136C"),
					CertificateId: uuid.MustParse("f0e5137f-03e1-4ca9-8dd9-b79da983d6be"),
				}),
			status: http.StatusOK,
			success: &api.UpdateCertificateResponse{
				ClientId: uuid.MustParse("63838416-B316-418A-8CC8-9EFE3411136C"),
				Certificate: api.Certificate{
					Id:                    uuid.MustParse("f0e5137f-03e1-4ca9-8dd9-b79da983d6be"),
					NotBefore:             types.Date{Time: time.Date(2025, time.May, 16, 0, 0, 0, 0, time.UTC)},
					NotAfter:              types.Date{Time: time.Date(2026, time.December, 25, 0, 0, 0, 0, time.UTC)},
					SerialNumber:          "303936854887858307147587267379171672533028646301",
					CertificatePemEncoded: certf0e5137f,
				},
			},
		},
		{
			name: "invalid client fails",
			req: testutil.
				NewRequest().
				Post("/api/v1/certificates").
				WithJsonBody(api.PostCertificatesJSONRequestBody{
					ClientId:      uuid.MustParse("49973760-cfdc-462c-be8c-d8db6d44f093"),
					CertificateId: uuid.MustParse("f0e5137f-03e1-4ca9-8dd9-b79da983d6be"),
				}),
			status: http.StatusNotFound,
			failed: &api.Error{
				Code:    http.StatusNotFound,
				Message: "client not found",
			},
		},
		{
			name: "invalid certificate fails",
			req: testutil.
				NewRequest().
				Post("/api/v1/certificates").
				WithJsonBody(api.PostCertificatesJSONRequestBody{
					ClientId:      uuid.MustParse("63838416-B316-418A-8CC8-9EFE3411136C"),
					CertificateId: uuid.MustParse("35d5113b-419c-4620-bb23-6d5d5b1cd361"),
				}),
			status: http.StatusNotFound,
			failed: &api.Error{
				Code:    http.StatusNotFound,
				Message: "certificate not found",
			},
		},
		{
			name: "updating with the same certificate fails",
			req: testutil.
				NewRequest().
				Post("/api/v1/certificates").
				WithJsonBody(api.PostCertificatesJSONRequestBody{
					ClientId:      uuid.MustParse("63838416-B316-418A-8CC8-9EFE3411136C"),
					CertificateId: uuid.MustParse("27bfe9fc-b80b-46a1-a967-78658db0aeec"),
				}),
			status: http.StatusBadRequest,
			failed: &api.Error{
				Code:    http.StatusBadRequest,
				Message: "current certificate matches existing certificate",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := http.NewServeMux()

			opts := api.StdHTTPServerOptions{
				BaseURL:    "/api/v1",
				BaseRouter: m,
			}

			api.HandlerWithOptions(api.New(newStore(t)), opts)

			rr := test.req.GoWithHTTPHandler(t, m).Recorder

			assert.Equal(t, test.status, rr.Code, rr.Body.String())

			if test.status == http.StatusOK {
				var actual api.UpdateCertificateResponse
				assert.NoError(t, json.NewDecoder(rr.Body).Decode(&actual))

				assert.Equal(t, *test.success, actual)
			} else {
				var actual api.Error
				assert.NoError(t, json.NewDecoder(rr.Body).Decode(&actual))
				assert.Equal(t, *test.failed, actual)
			}
		})
	}
}

func TestFormatSerialNumber(t *testing.T) {
	t.Skip()

	bigInt := func(t *testing.T, s string) *big.Int {
		t.Helper()
		i := new(big.Int)
		_, err := fmt.Sscan(s, i)
		assert.NoError(t, err)
		return i
	}

	tests := []struct {
		serialNumber *big.Int
		expected     string
	}{
		{
			serialNumber: bigInt(t, "520097931764758754948131491264192108474121041456"),
			expected:     "5b:19:ff:73:a8:d9:d2:56:cd:fe:8b:07:cf:29:eb:5b:4d:53:b6:30",
		},
		{
			serialNumber: bigInt(t, "303936854887858307147587267379171672533028646301"),
			expected:     "35:3c:ff:fb:d1:84:6e:ab:7d:82:3c:df:9f:4e:47:52:81:c9:2d:9d",
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			actual := api.FormatSerialNumber(test.serialNumber)
			assert.Equal(t, test.expected, actual)
		})
	}
}

type store struct {
	clients      map[uuid.UUID]model.Client
	certificates map[uuid.UUID]model.Certificate
}

func newStore(t *testing.T) *store {
	t.Helper()

	s := &store{
		clients:      make(map[uuid.UUID]model.Client),
		certificates: make(map[uuid.UUID]model.Certificate),
	}

	var clients []model.Client
	assert.NoError(t, json.Unmarshal(clientsJSON, &clients))

	for _, client := range clients {
		s.clients[client.ID] = client
	}

	var certificates []model.Certificate
	assert.NoError(t, json.Unmarshal(certificatesJSON, &certificates))

	for _, certificate := range certificates {
		certificate.CertificatePEMEncoded = bytes.TrimSpace(certificate.CertificatePEMEncoded)
		s.certificates[certificate.ID] = certificate
	}

	return s
}

func (s *store) GetClient(_ context.Context, clientID uuid.UUID) (*model.Client, error) {
	c, ok := s.clients[clientID]
	if !ok {
		return nil, model.ErrClientNotFound
	}
	return &c, nil
}

func (s *store) GetCertificate(ctx context.Context, certificateID uuid.UUID) (*model.Certificate, error) {
	c, ok := s.certificates[certificateID]
	if !ok {
		return nil, model.ErrCertificateNotFound
	}
	return &c, nil
}

func (s *store) ListCertificates(ctx context.Context) ([]*model.Certificate, error) {
	certificates := make([]*model.Certificate, 0)

	for _, certificate := range s.certificates {
		certificates = append(certificates, &certificate)
	}

	slices.SortStableFunc(certificates, func(a, b *model.Certificate) int {
		return strings.Compare(a.ID.String(), b.ID.String())
	})
	return certificates, nil
}

func (s *store) UpdateClientCertificate(ctx context.Context, clientID, certificateID uuid.UUID) error {
	c, err := s.GetClient(ctx, clientID)
	if err != nil {
		return err
	}
	c.CertificateID = certificateID

	return nil
}
