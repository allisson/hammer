package grpc

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validationStatusError(code codes.Code, msg string, err error) *status.Status {
	st := status.New(code, msg)

	switch t := err.(type) {
	case validation.Errors:
		br := &errdetails.BadRequest{}
		for id, e := range t {
			ve := e.(validation.Error)
			v := &errdetails.BadRequest_FieldViolation{
				Field:       id,
				Description: ve.Message(),
			}
			br.FieldViolations = append(br.FieldViolations, v)
		}
		newSt, err := st.WithDetails(br)
		if err == nil {
			return newSt
		}
	}

	return st
}
