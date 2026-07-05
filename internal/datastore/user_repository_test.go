package datastore

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/bytedance/go-noba-trial/internal/domain"
)

func TestUserRepositoryCreate(t *testing.T) {
	db := newStubDB(t, &stubDB{
		query: func(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
			if !strings.Contains(query, `INSERT INTO "users"`) {
				t.Fatalf("expected insert query, got %q", query)
			}

			expectedArgs := []any{"Jane Doe", "jane@example.com", "Example user"}
			for i, expected := range expectedArgs {
				if args[i].Value != expected {
					t.Fatalf("expected arg %d to be %q, got %#v", i, expected, args[i].Value)
				}
			}

			return &stubRows{
				columns: []string{"uuid", "fullname", "email", "description"},
				values: [][]driver.Value{{
					"f4b2fe41-4b68-42a9-8db2-8563dc5c7eb9",
					"Jane Doe",
					"jane@example.com",
					"Example user",
				}},
			}, nil
		},
	})
	defer db.Close()

	repository := NewUserRepository(db)
	user, err := repository.Create(context.Background(), domain.CreateUserInput{
		Fullname:    "Jane Doe",
		Email:       "jane@example.com",
		Description: "Example user",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Email != "jane@example.com" {
		t.Fatalf("expected email jane@example.com, got %q", user.Email)
	}
}

func TestUserRepositoryCreateReturnsQueryError(t *testing.T) {
	expectedErr := errors.New("insert failed")
	db := newStubDB(t, &stubDB{
		query: func(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
			return nil, expectedErr
		},
	})
	defer db.Close()

	repository := NewUserRepository(db)
	_, err := repository.Create(context.Background(), domain.CreateUserInput{
		Fullname: "Jane Doe",
		Email:    "jane@example.com",
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected query error, got %v", err)
	}
}

func TestUserRepositoryList(t *testing.T) {
	db := newStubDB(t, &stubDB{
		query: func(_ context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
			if !strings.Contains(query, `SELECT uuid, fullname, email, description`) {
				t.Fatalf("expected select query, got %q", query)
			}

			return &stubRows{
				columns: []string{"uuid", "fullname", "email", "description"},
				values: [][]driver.Value{
					{
						"f4b2fe41-4b68-42a9-8db2-8563dc5c7eb9",
						"Jane Doe",
						"jane@example.com",
						"Example user",
					},
					{
						"5c8a9be7-8d62-4d1f-8cf1-80cb5f1d3c55",
						"John Smith",
						"john@example.com",
						"Another user",
					},
				},
			}, nil
		},
	})
	defer db.Close()

	repository := NewUserRepository(db)
	users, err := repository.List(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}

	if users[0].Email != "jane@example.com" {
		t.Fatalf("expected first email jane@example.com, got %q", users[0].Email)
	}
}

func TestUserRepositoryListReturnsScanError(t *testing.T) {
	db := newStubDB(t, &stubDB{
		query: func(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
			return &stubRows{
				columns: []string{"uuid", "fullname", "email", "description"},
				values: [][]driver.Value{{
					"f4b2fe41-4b68-42a9-8db2-8563dc5c7eb9",
					"Jane Doe",
				}},
			}, nil
		},
	})
	defer db.Close()

	repository := NewUserRepository(db)
	_, err := repository.List(context.Background())

	if err == nil {
		t.Fatal("expected scan error")
	}
}

func TestEnsureSchemaExecutesUserSchema(t *testing.T) {
	var executedQuery string
	db := newStubDB(t, &stubDB{
		exec: func(_ context.Context, query string, _ []driver.NamedValue) (driver.Result, error) {
			executedQuery = query
			return driver.RowsAffected(0), nil
		},
	})
	defer db.Close()

	if err := EnsureSchema(context.Background(), db); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(executedQuery, userSchema) {
		t.Fatal("expected EnsureSchema to execute userSchema")
	}
}

func TestEnsureSchemaReturnsExecError(t *testing.T) {
	expectedErr := errors.New("schema failed")
	db := newStubDB(t, &stubDB{
		exec: func(context.Context, string, []driver.NamedValue) (driver.Result, error) {
			return nil, expectedErr
		},
	})
	defer db.Close()

	err := EnsureSchema(context.Background(), db)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected exec error, got %v", err)
	}
}

type stubDB struct {
	query func(context.Context, string, []driver.NamedValue) (driver.Rows, error)
	exec  func(context.Context, string, []driver.NamedValue) (driver.Result, error)
}

var (
	stubDBsMu sync.Mutex
	stubDBs   = map[string]*stubDB{}
)

func newStubDB(t *testing.T, stub *stubDB) *sql.DB {
	t.Helper()

	registerStubSQLDriver()

	name := strings.ReplaceAll(t.Name(), "/", "_")

	stubDBsMu.Lock()
	stubDBs[name] = stub
	stubDBsMu.Unlock()

	t.Cleanup(func() {
		stubDBsMu.Lock()
		delete(stubDBs, name)
		stubDBsMu.Unlock()
	})

	db, err := sql.Open("datastore_stub", name)
	if err != nil {
		t.Fatalf("failed to open stub database: %v", err)
	}

	return db
}

var registerStubSQLDriverOnce sync.Once

func registerStubSQLDriver() {
	registerStubSQLDriverOnce.Do(func() {
		sql.Register("datastore_stub", stubDriver{})
	})
}

type stubDriver struct{}

func (stubDriver) Open(name string) (driver.Conn, error) {
	stubDBsMu.Lock()
	stub := stubDBs[name]
	stubDBsMu.Unlock()

	if stub == nil {
		return nil, errors.New("stub database not found")
	}

	return stubConn{stub: stub}, nil
}

type stubConn struct {
	stub *stubDB
}

func (stubConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepare is not implemented")
}

func (stubConn) Close() error {
	return nil
}

func (stubConn) Begin() (driver.Tx, error) {
	return nil, errors.New("transactions are not implemented")
}

func (c stubConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if c.stub.query == nil {
		return nil, errors.New("query is not implemented")
	}

	return c.stub.query(ctx, query, args)
}

func (c stubConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if c.stub.exec == nil {
		return nil, errors.New("exec is not implemented")
	}

	return c.stub.exec(ctx, query, args)
}

type stubRows struct {
	columns []string
	values  [][]driver.Value
	index   int
}

func (r *stubRows) Columns() []string {
	return r.columns
}

func (r *stubRows) Close() error {
	return nil
}

func (r *stubRows) Next(dest []driver.Value) error {
	if r.index >= len(r.values) {
		return io.EOF
	}

	copy(dest, r.values[r.index])
	r.index++
	return nil
}
