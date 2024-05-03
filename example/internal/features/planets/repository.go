package planets

type IReadRepository interface {
	Get() error
}

type IWriteRepository interface {
	Update() error
}
