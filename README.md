# Certificate API

At Empyrean, many of our clients have SSO (Single-Sign On) connections to many different providers. These SSO connections have an associated certificate and we must manage the lifecycle of these certificates since they are only valid for a period of time.

## Context

The API is defined using an [OpenAPI](https://www.openapis.org/) specification which can be viewed in [`specification.yaml`](./specification.yaml) and rendered using the [online editor](https://editor-next.swagger.io/). A code generator transforms the OpenAPI specification into Go primitives that can be used to implement the API, these can be viewed in [`api.gen.go`](./api/api.gen.go).

We will only be focusing on two endpoints:

* `GET /api/v1/certificates` Returns a list of certificates
* `POST /api/v1/certificates` Updates the certificate for a client

## Challenge

The task is to make the tests pass, the tests can be found in [`api_test.go`](./api/api_test.go). The tests can be run with `go test -v ./...` or via any other method of your choosing. You shouldn't need to modify [`api_test.go`](./api/api_test.go) or rerun the code generation to complete this task.

### Bonus Challenge

Rather than rendering a serial number of a certificate as a decimal, it is usually more common to render it in hexadecimal format with colons (`:`) as a separator every two characters.

For example, the serial number `520097931764758754948131491264192108474121041456` could be represented as `5b:19:ff:73:a8:d9:d2:56:cd:fe:8b:07:cf:29:eb:5b:4d:53:b6:30`.

Remove `t.Skip()` in the `TestFormatSerialNumber` test and implement the stubbed out `FormatSerialNumber` function in [`api.go`](./api/api.go).

## Hints

* The [`x509`](https://pkg.go.dev/crypto/x509) and [`pem`](https://pkg.go.dev/encoding/pem) packages in the standard library are useful for working with certificates.
* `openssl` can parse PEM encoded certificates, the data in [`certificates.json`](./api/certificates.json) can be viewed on the command line with `echo 'BASE64_PEM_ENCODED_CERTIFICATE' | base64 -d | openssl x509 -text -in /dev/stdin`.
