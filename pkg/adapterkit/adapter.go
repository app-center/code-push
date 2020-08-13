package adapterkit

type Adaptable interface {
	Conn() error
	Close() error
}
