package siwebsocket

import "context"

// id:            key "9099901_165695708447103008123"
// clientID:      key_sub "9099901"
// clientGroupID: channel "90999"
// hubAddr:       value "http://172.16.130.144:45501"
// hubPath:       route "/ws/message/adpos/emg/_push"
// status "0"
type Router interface {
	Store(ctx context.Context, ID string, userID, userGroupID string, hubAddr, hubPath string) error
	Delete(ctx context.Context, ID string) error
	DeleteByHubAddr(ctx context.Context, hubAddr string) error
}

type NopRouter struct{}

func (n *NopRouter) Store(ctx context.Context, ID string, userID, userGroupID string, hubAddr, hubPath string) error {
	return nil
}
func (n *NopRouter) Delete(ctx context.Context, ID string) error {
	return nil
}
func (n *NopRouter) DeleteByHubAddr(ctx context.Context, hubAddr string) error {
	return nil
}
