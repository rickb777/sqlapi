package testenv

import (
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"os"
	"syscall"
	"testing"
	"time"
)

func SetEnvironmentForLocalPostgres() {
	log.Println("Attempting to connect to local postgres")
	mustSetEnv("DB_DRIVER", "pgx")
	mustSetEnv("DB_DIALECT", "postgres")
	mustSetEnv("PGHOST", "localhost")
	mustSetEnv("PGPORT", "5432")
	mustSetEnv("PGDATABASE", "test")
	mustSetEnv("PGUSER", "testuser")
	mustSetEnv("PGPASSWORD", "TestPasswd.9.9.9")
}

func SetUpDockerDbForTest(m *testing.M, repo string, runTestSetup func()) {
	log.Printf("Attempting to connect to docker %s\n", repo)
	mustSetEnv("DB_DRIVER", "pgx")
	mustSetEnv("DB_DIALECT", "postgres")
	mustSetEnv("PGHOST", "localhost")
	mustSetEnv("PGPORT", "15432")
	mustSetEnv("PGDATABASE", "postgres")
	mustSetEnv("PGUSER", "postgres")
	mustSetEnv("PGPASSWORD", "simple")

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker %s: %s", repo, err)
	}

	// pulls an image, creates a container based on it and runs it
	opts := &dockertest.RunOptions{
		Repository: repo,
		Tag:        "12-alpine",
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432/tcp": {{HostPort: "15432/tcp"}},
		},
		Env: []string{"PGPASSWORD=simple", "POSTGRES_PASSWORD=simple"},
	}
	resource, err := pool.RunWithOptions(opts)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// docker always takes some time to start
	time.Sleep(1950 * time.Millisecond)

	runTestSetup()

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func mustSetEnv(name, defaultValue string) {
	err := syscall.Setenv(name, defaultValue)
	if err != nil {
		log.Fatalf("Failed to set %q=%q; %v\n", name, defaultValue, err)
	}
}
