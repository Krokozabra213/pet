package chatdomain

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func timeToProto(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}
