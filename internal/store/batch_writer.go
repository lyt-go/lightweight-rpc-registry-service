package store

type BatchSession interface {
	Put(string) error
	Commit() error
	Rollback() error
	Close() error
}

type BatchSessionFactory interface {
	Begin() (BatchSession, error)
}

type BatchRepository struct{ factory BatchSessionFactory }

func NewBatchRepository(factory BatchSessionFactory) *BatchRepository {
	return &BatchRepository{factory: factory}
}

func (r *BatchRepository) Persist(names []string) (err error) {
	sessions := make([]BatchSession, 0, len(names))
	for _, name := range names {
		session, beginErr := r.factory.Begin()
		if beginErr != nil {
			return beginErr
		}
		sessions = append(sessions, session)
		defer session.Close()
		defer func() { err = session.Commit() }()
		if putErr := session.Put(name); putErr != nil {
			return putErr
		}
	}
	return nil
}
