package views

type Builder interface {
	Key() uint32
	Types() []string
	OnDomainEvent func(DomainEvent) error
}
