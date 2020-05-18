package cf

type ResourceBuilder interface {
	JSON() ([]byte, error)
	Name() string
}
