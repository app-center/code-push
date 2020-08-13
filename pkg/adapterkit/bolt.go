package adapterkit

type BoltAdaptOption func(adapter *boltAdapter)

func AdaptBolt() Adaptable {
	return nil
}

type boltAdapter struct {
}

func (a *boltAdapter) Conn() error {
	panic("implement me")
}

func (a *boltAdapter) Close() error {
	panic("implement me")
}
