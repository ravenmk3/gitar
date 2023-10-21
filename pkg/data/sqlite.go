package data

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func NewSqlite3DataStore(dsn string) DataStore {
	store := &Sqlite3DataStore{
		dsn: dsn,
		db:  nil,
	}
	return store
}

type Sqlite3DataStore struct {
	dsn string
	db  *sqlx.DB
}

func (me *Sqlite3DataStore) String() string {
	return fmt.Sprintf("Sqlite3DataStore { DataSource: %s }", me.dsn)
}

func (me *Sqlite3DataStore) Open() error {
	db, err := sqlx.Open("sqlite3", me.dsn)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	me.db = db
	return me.initDatabase()
}

func (me *Sqlite3DataStore) Close() error {
	if me.db == nil {
		return errors.New("db is not open")
	}
	return me.db.Close()
}

func (me *Sqlite3DataStore) initDatabase() error {
	cmd := `
	CREATE TABLE IF NOT EXISTS [git_repo] (
		[repo] TEXT NOT NULL PRIMARY KEY
	);

	CREATE TABLE IF NOT EXISTS [commit_downloaded] (
		[id] TEXT NOT NULL PRIMARY KEY
	);

	CREATE TABLE IF NOT EXISTS [commit_mailed] (
		[id] TEXT NOT NULL PRIMARY KEY
	);
	`
	_, err := me.db.Exec(cmd)
	if err != nil {
		return err
	}
	return nil
}

func (me *Sqlite3DataStore) queryExists(table, field string, value any) (bool, error) {
	cmd := fmt.Sprintf("SELECT count(*) FROM [%s] WHERE [%s] = ?;", table, field)
	rows, err := me.db.Query(cmd, value)
	if err != nil {
		return false, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	rows.Next()
	var count int32
	err = rows.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (me *Sqlite3DataStore) RepoExists(repo string) (bool, error) {
	return me.queryExists("git_repo", "repo", repo)
}

func (me *Sqlite3DataStore) SaveRepo(repo string) error {
	exists, err := me.RepoExists(repo)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	cmd := "INSERT INTO [git_repo] ([repo]) VALUES(?);"
	_, err = me.db.Exec(cmd, repo)
	return err
}

func (me *Sqlite3DataStore) IsCommitDownloaded(id string) (bool, error) {
	return me.queryExists("commit_downloaded", "id", id)
}

func (me *Sqlite3DataStore) SetCommitDownloaded(id string) error {
	cmd := "INSERT INTO [commit_downloaded] ([id]) VALUES(?);"
	_, err := me.db.Exec(cmd, id)
	return err
}

func (me *Sqlite3DataStore) IsCommitMailed(id string) (bool, error) {
	return me.queryExists("commit_mailed", "id", id)
}

func (me *Sqlite3DataStore) SetCommitMailed(id string) error {
	cmd := "INSERT INTO [commit_mailed] ([id]) VALUES(?);"
	_, err := me.db.Exec(cmd, id)
	return err
}
