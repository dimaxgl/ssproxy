package database

import (
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"log"
	"github.com/mitchellh/mapstructure"
	"github.com/dimaxgl/ssproxy/store"
)

const (
	StoreName = `database`

	AddUserExecQuery = `INSERT INTO users ("user", "password") VALUES (:user,:password)`
	GetUserQuery     = `SELECT password FROM users WHERE "user" = :user`
	UserColumn       = `user`
)

type dbStore struct {
	conn *sqlx.DB
	opts *dbOptions
}

type dbOptions struct {
	DriverName string    `yaml:"driverName"`
	Dsn        string    `yaml:"dsn"`
	UserColumn string    `yaml:"userColumn"`
	Queries    dbQueries `yaml:"queries"`
}

type dbQueries struct {
	AddUserExecQuery string
	GetUserQuery     string
}

func (d dbStore) Initialize(params map[string]interface{}) (store.Store, error) {
	var options dbOptions
	var err error

	if err = mapstructure.Decode(params, &options); err != nil {
		return nil, errors.Wrap(err, `failed to decode params`)
	}

	// set default user column if not set
	if options.UserColumn == `` {
		options.UserColumn = UserColumn
	}

	// set default adding user sql query if not set
	if options.Queries.AddUserExecQuery == `` {
		options.Queries.AddUserExecQuery = AddUserExecQuery
	}

	// set default search user sql query if not set
	if options.Queries.GetUserQuery == `` {
		options.Queries.GetUserQuery = GetUserQuery
	}

	d.opts = &options
	if d.conn, err = sqlx.Connect(options.DriverName, options.Dsn); err != nil {
		return nil, errors.Wrap(err, `failed to create sql connection`)
	}

	return d, nil
}

func (d dbStore) Valid(user, password string) bool {

	var passwordHash string

	if rows, err := d.conn.NamedQuery(d.opts.Queries.GetUserQuery, map[string]interface{}{
		d.opts.UserColumn: user,
	}); err != nil {
		log.Println(err)
		return false
	} else {
		defer rows.Close()
		for rows.Next() {
			if err = rows.Scan(&passwordHash); err != nil {
				log.Println(err)
				return false
			} else {
				log.Println(password, passwordHash)
				return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) == nil
			}
		}
	}
	return false
}

func (d dbStore) Add(user, password string) error {

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, `failed to encrypt password using bcrypt`)
	}

	if result, err := d.conn.NamedExec(d.opts.Queries.AddUserExecQuery, map[string]interface{}{
		`user`: user, `password`: string(passwordHash),
	}); err != nil {
		return errors.Wrap(err, `failed to execute addUser query`)
	} else {
		if rowsCount, err := result.RowsAffected(); err != nil {
			return errors.Wrap(err, `failed to get count of affected rows`)
		} else {
			if rowsCount <= 0 {
				return errors.Errorf("affected rows: %d, maybe problem in your database?", rowsCount)
			}
		}
	}
	return nil
}

func init() {
	store.Register(StoreName, dbStore{})
}
