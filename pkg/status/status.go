package status

type Status interface {
	ID() string
	TypeName() string
}
