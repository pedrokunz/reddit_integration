package redditpostsync

type postgresRepository struct{}

func Postgres() Repository {
	return postgresRepository{}
}

func (p postgresRepository) PostCreate(input PostCreateInput) (PostCreateOutput, error) {
	return PostCreateOutput{}, nil
}
