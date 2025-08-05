package cloudflare

import (
	"errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrCloudflareConfigNotInitialized = errors.New("cloudflare config not initialized")
	ErrCloudflareCallFailed           = errors.New("cloudflare call failed")

	Status_CloudflareError = status.New(codes.Internal, "cloudflare error")
)

func ToStatus(err error) *status.Status {
	if err == nil {
		return nil
	}
	unwrapErr := errors.Unwrap(err)
	if unwrapErr == nil {
		unwrapErr = err
	}
	st, myErr := Status_CloudflareError.WithDetails(
		&errdetails.PreconditionFailure{
			Violations: []*errdetails.PreconditionFailure_Violation{
				{
					Type:        "CLOUDFLARE",
					Subject:     unwrapErr.Error(),
					Description: err.Error(),
				},
			},
		},
	)
	if myErr != nil {
		return Status_CloudflareError
	}
	return st
}
