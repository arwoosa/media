package db

import (
	"errors"

	"github.com/arwoosa/vulpes/relation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrRelation = errors.New("relation error")
)

func ToStatus(err error) *status.Status {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, ErrRelation):
		return relation.ToStatus(err)
	default:
		unwrapErr := errors.Unwrap(err)
		if unwrapErr == nil {
			unwrapErr = err
		}
		baseErrStatus := status.New(codes.Internal, err.Error())
		st, myErr := baseErrStatus.WithDetails(
			&errdetails.PreconditionFailure{
				Violations: []*errdetails.PreconditionFailure_Violation{
					{
						Type:        "MEDIA_DB",
						Subject:     unwrapErr.Error(),
						Description: err.Error(),
					},
				},
			},
		)
		if myErr != nil {
			return baseErrStatus
		}
		return st
	}
}
