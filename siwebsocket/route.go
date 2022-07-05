package siwebsocket

// id:            key "9099901_165695708447103008123"
// clientID:      key_sub "9099901"
// clientGroupID: channel "90999"
// hubAddr:       value "http://172.16.130.144:45501"
// hubPath:       route "/ws/message/adpos/emg/_push"
// status "0"
type ClientStorage interface {
	Store(ID string, clientID, clientGroupID string, hubAddr, hubPath string) error
	Delete(ID string) error
	DeleteByHubAddr(hubAddr string) error
}

type NopRouteStorage struct{}

func (n *NopRouteStorage) Store(ID string, clientID, clientGroupID string, hubAddr, hubPath string) error {
	return nil
}
func (n *NopRouteStorage) Delete(ID string) error {
	return nil
}
func (n *NopRouteStorage) DeleteByHubAddr(hubAddr string) error {
	return nil
}
