package sicore

type Client interface {
	Stop() error
	Wait() error
	Send(b []byte) error
	SendAndWait(b []byte) error

	GetID() string
	GetUserID() string
	GetUserGroupID() string
}

type Hub interface {
	Add(c Client) error
	Remove(c Client) error
}

type NopHub struct{}

func (o NopHub) Add(c Client) error {
	return nil
}
func (o NopHub) Remove(c Client) error {
	return nil
}
