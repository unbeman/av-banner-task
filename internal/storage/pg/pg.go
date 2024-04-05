package pg

type PG struct {
}

func NewPG(dsn string) (*PG, error) {
	pg := &PG{}
	return pg, nil
}
