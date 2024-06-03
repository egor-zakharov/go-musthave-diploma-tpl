package migrator

type Migrator interface {
	Run() error
}
