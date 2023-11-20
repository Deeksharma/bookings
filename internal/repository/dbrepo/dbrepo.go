package dbrepo

import (
	"github.com/Deeksharma/bookings/internal/config"
	"github.com/Deeksharma/bookings/internal/repository"
	"github.com/jackc/pgx/v4/pgxpool"
)

type postgresDBRepo struct { // this should implement DatabaseRepo interface
	App *config.AppConfig
	DB  *pgxpool.Pool
	// mybe use sync.Mutex
}

type testDBRepo struct { // this should implement DatabaseRepo interface
	App *config.AppConfig
	DB  *pgxpool.Pool
}


// the above one is for postgres, if you want to use another db just make a new struct with appconfig and db connection pools
// and make them implement the DatabaseRepo interface

func NewPostgresRepo(conn *pgxpool.Pool, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
