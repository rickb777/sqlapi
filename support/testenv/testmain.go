package testenv

import (
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/log/testingadapter"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	pkgerrors "github.com/pkg/errors"
)

func Shebang(m *testing.M, connectFunc func(lgr tracelog.Logger, logLevel tracelog.LogLevel, tries int) error) {
	flag.Parse()

	var lvl tracelog.LogLevel = tracelog.LogLevelWarn
	if testing.Verbose() {
		lvl = tracelog.LogLevelInfo
	}

	lgr := testingadapter.NewLogger(simpleLogger{})

	switch os.Getenv("DB_DRIVER") {
	case "sqlite3":
		mustUnsetEnv("DB_DRIVER")
		mustUnsetEnv("DB_URL")
		mustUnsetEnv("PGHOST")
		mustUnsetEnv("PGPORT")
		mustUnsetEnv("PGDATABASE")
		mustUnsetEnv("PGUSER")
		mustUnsetEnv("PGPASSWORD")
		mustUnsetEnv("MYUSER")
		mustUnsetEnv("MYPASSWORD")
		err := connectFunc(lgr, lvl, 1)
		if err == nil {
			log.Printf("Test using Travis DB\n")
			os.Exit(m.Run())
		}
	}

	//abortIfNotConnectionError("Cannot connect to Travis DB", err)

	setEnvironmentDockerDb()

	// second connection attempt: use pre-existing Docker DB, if it exists
	//log.Printf("----- Second attempt ----- (%s)\n", dfltDriver)
	//err := connectFunc(lgr, lvl, 1)
	//if err == nil {
	//	log.Printf("Test using local DB\n")
	//	os.Exit(m.Run())
	//}
	//
	//abortIfNotConnectionError("Cannot connect to pre-existing Docker DB", err)

	// third connection attempt: spin up DB in Docker container and connect to it
	//log.Printf("----- Third attempt -----\n")
	setUpDockerDbForTest(m, "postgres", func() error {
		return connectFunc(lgr, lvl, 0)
	})
}

//-------------------------------------------------------------------------------------------------
// PostgresQL URL general form
// postgresql://[user[:password]@][netloc][:port][,...][/dbname][?param1=value1&...]

func SetDefaultDbDriver(dfltDriver string) string {
	dbDriver := os.Getenv("DB_DRIVER")
	if dbDriver == "" {
		dbDriver = dfltDriver
		mustSetEnv("DB_DRIVER", dbDriver)
	}
	return dbDriver
}

func setEnvironmentForTravisDB(dfltDriver string) {
	dbDriver := SetDefaultDbDriver(dfltDriver)
	log.Printf("set environment for Travis %s DB\n", dbDriver)

	switch dbDriver {
	case "sqlite3":
		mustUnsetEnv("DB_DRIVER")
		mustUnsetEnv("DB_URL")
		mustUnsetEnv("PGHOST")
		mustUnsetEnv("PGPORT")
		mustUnsetEnv("PGDATABASE")
		mustUnsetEnv("PGUSER")
		mustUnsetEnv("PGPASSWORD")
		mustUnsetEnv("MYUSER")
		mustUnsetEnv("MYPASSWORD")

	case "postgres", "pgx":
		log.Println("Attempting to connect to local postgres")
		mustSetEnv("PGHOST", "localhost")
		mustSetEnv("PGPORT", "15432")
		mustSetEnv("PGDATABASE", "postgres")
		mustSetEnv("PGUSER", "postgres")
		mustSetEnv("PGPASSWORD", "")
		mustSetEnv("DB_URL", "postgres:@/postgres")

	case "mysql":
		log.Println("Attempting to connect to local mysql")
		mustSetEnv("DB_DRIVER", "mysql")
		mustSetEnv("DB_URL", "travis:@/test")
		mustSetEnv("MYUSER", "travis")
		mustSetEnv("MYPASSWORD", "")
	}
}

func setEnvironmentDockerDb() {
	dbDriver := os.Getenv("DB_DRIVER")
	log.Printf("set environment for Docker %s DB\n", dbDriver)

	switch dbDriver {
	case "postgres", "pgx":
		mustSetEnv("PGHOST", "localhost")
		mustSetEnv("PGPORT", "15432")
		mustSetEnv("PGDATABASE", "postgres")
		mustSetEnv("PGUSER", "postgres")
		mustSetEnv("PGPASSWORD", "simple")
		mustSetEnv("PGSSLMODE", "disable")
		mustSetEnv("DB_URL", "postgres:simple@localhost:15432/postgres")
	}
}

func setUpDockerDbForTest(m *testing.M, repo string, runTestSetup func() error) {
	log.Printf("Spinning up docker %s\n", repo)

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker %s for %s: %s", repo, os.Getenv("DB_DRIVER"), err)
	}

	// pulls an image, creates a container based on it and runs it
	opts := &dockertest.RunOptions{
		Name:       "postgres4test",
		Repository: repo,
		Tag:        "13-alpine",
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432/tcp": {{HostPort: "15432/tcp"}},
		},
		Env: []string{"PGPASSWORD=simple", "POSTGRES_PASSWORD=simple"},
	}

	resource, err := pool.RunWithOptions(opts)
	if err != nil {
		e2 := pkgerrors.Cause(err)
		if e2 != docker.ErrContainerAlreadyExists {
			switch e3 := e2.(type) {
			case *docker.Error:
				if e3.Status != 500 {
					log.Fatalf("Could not start docker+postgres resource: %v, %#v", err, e3)
				}
			default:
				log.Fatalf("Could not start docker+postgres resource: %v", err)
			}
		}
	}

	// docker always takes some time to start
	time.Sleep(5000 * time.Millisecond)

	err = runTestSetup()
	if err != nil {
		if resource != nil {
			pool.Purge(resource)
		}
		log.Fatalf("Could not connect to DB in docker+postgres: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if resource != nil {
		if err = pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge docker+postgres resource: %s", err)
		}
	}

	os.Exit(code)
}

//-------------------------------------------------------------------------------------------------

type simpleLogger struct{}

func (l simpleLogger) Log(args ...interface{}) {
	if testing.Verbose() {
		log.Println(args...)
	}
}

//-------------------------------------------------------------------------------------------------

func abortIfNotConnectionError(msg string, err error) {
	var connErr *net.OpError
	if !errors.As(err, &connErr) {
		log.Fatalf("%s: %s", msg, err)
	}
}

//-------------------------------------------------------------------------------------------------

func mustUnsetEnv(name string) {
	err := syscall.Unsetenv(name)
	if err != nil {
		log.Fatalf("Failed to unset %q; %v\n", name, err)
	}
}

func mustSetEnv(name, defaultValue string) {
	err := syscall.Setenv(name, defaultValue)
	if err != nil {
		log.Fatalf("Failed to set %q=%q; %v\n", name, defaultValue, err)
	}
}
