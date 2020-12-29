package db

import (
	"fmt"
	"github.com/Droidion/implementing-change-game/models"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rotisserie/eris"
	"os"
)

var PgConn *pgxpool.Pool

// PostgresConnect tries to connect to Postgres database
func PostgresConnect() {
	var err error
	PgConn, err = pgxpool.Connect(Ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		fmt.Println(eris.ToString(err, true))
		os.Exit(1)
	}
}

// GetUserByLogin searches for a user with a given login in a database
func GetUserByLogin(login string) (*models.User, error) {
	var users []*models.User
	err := pgxscan.Select(Ctx, PgConn, &users, `SELECT id, login, password FROM users WHERE login=$1`, login)
	if err != nil {
		return nil, eris.Wrap(err, "problem with quering user from a db")
	}
	if users == nil {
		return nil, eris.New("no users found in a db")
	}
	return users[0], nil
}
