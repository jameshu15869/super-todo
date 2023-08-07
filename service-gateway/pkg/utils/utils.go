package utils

import (
	"fmt"
	"supertodo/gateway/pb"
	"time"
)

func FormatJsonError(name string, err error) *pb.JsonError {
	return &pb.JsonError{
		Name:      name,
		Message:   err.Error(),
		Timestamp: fmt.Sprintf("%s", time.Now()),
	}
}
