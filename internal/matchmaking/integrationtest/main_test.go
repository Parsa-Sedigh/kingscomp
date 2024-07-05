package integrationtest

import (
	"fmt"
	"github.com/Parsa-Sedigh/kingscomp/internal/repository/redis"
	"github.com/Parsa-Sedigh/kingscomp/pkg/testhelper"
	"github.com/ory/dockertest/v3"
	"os"
	"testing"
)

var redisPort string

func TestMain(m *testing.M) {
	if !testhelper.IsIntegration() {
		return
	}

	pool := testhelper.StartDockerPool()

	// set up redis container for tests
	redisResource := testhelper.StartDockerInstance(pool, "redis/redis-stack-server", "latest", func(res *dockertest.Resource) error {
		port := res.GetPort("6379/tcp")
		_, err := redis.NewRedisClient(fmt.Sprintf("localhost:%s", port))

		return err
	})

	fmt.Println(redisResource.GetPort("6379/tcp"))

	// now run the tests
	exitCode := m.Run()
	redisResource.Close()
	os.Exit(exitCode)
}
