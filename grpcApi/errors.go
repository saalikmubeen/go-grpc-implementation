package grpcApi

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func fieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func invalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {

	badRequest := &errdetails.BadRequest{FieldViolations: violations}
	statusInvalid := status.New(codes.InvalidArgument, "invalid parameters")

	// Add more details about those invalid parameters to the status object.
	statusDetails, err := statusInvalid.WithDetails(badRequest)

	// There is something wrong with the badRequest details
	if err != nil {
		return statusInvalid.Err()
	}

	return statusDetails.Err()

	// statusInvalid.Err() will return the error of following structure:
	// 	{
	//   "code": 3,
	//   "message": "invalid parameters",
	//   "details": [
	//     {
	//       "@type": "type.googleapis.com/google.rpc.BadRequest",
	//       "field_violations": [
	//         {
	//           "field": "password",
	//           "description": "must contain from 6-100 characters"
	//         },
	//         {
	//           "field": "email",
	//           "description": "must contain from 3-200 characters"
	//         }
	//       ]
	//     }
	//   ]
	// }

}

func unauthenticatedError(err error) error {
	return status.Errorf(codes.Unauthenticated, "unauthorized: %s", err)
}
