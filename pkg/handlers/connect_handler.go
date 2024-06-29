package handlers

import "context"

type Connection interface {
	InitiateConnectionRequest(ctx context.Context, requester string, accepter string) (bool, error)
	AcceptConnectionRequest(ctx context.Context, accepter string) (bool, error)
	RejectConnectionRequest(ctx context.Context, accepter string) (bool, error)
	ConnectionDetails(ctx context.Context)
}

var _ Connection = &connectionStruct{}

type connectionStruct struct {
}

func (connstruct *connectionStruct) InitiateConnectionRequest(ctx context.Context, requester string, accepter string) (bool, error) {
	return false, nil
}

func (connstruct *connectionStruct) AcceptConnectionRequest(ctx context.Context, accepter string) (bool, error) {
	return false, nil
}

func (connstruct *connectionStruct) RejectConnectionRequest(ctx context.Context, accepter string) (bool, error) {
	return false, nil
}

func (connstruct *connectionStruct) ConnectionDetails(ctx context.Context) {}
