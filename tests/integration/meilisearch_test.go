//go:build integration

package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func startMeilisearch(ctx context.Context, t *testing.T) (meilisearch.ServiceManager, func()) {
	t.Helper()
	masterKey := "test-master-key-integration"
	req := testcontainers.ContainerRequest{
		Image:        "getmeili/meilisearch:v1.12",
		ExposedPorts: []string{"7700/tcp"},
		Env: map[string]string{
			"MEILI_MASTER_KEY": masterKey,
		},
		WaitingFor: wait.ForLog("Meilisearch is ready").WithStartupTimeout(30 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)
	natPort, err := container.MappedPort(ctx, "7700")
	require.NoError(t, err)

	url := fmt.Sprintf("http://%s:%s", host, natPort.Port())
	client := meilisearch.New(url, meilisearch.WithAPIKey(masterKey))

	cleanup := func() {
		container.Terminate(ctx)
	}

	return client, cleanup
}

func TestIntegration_Meilisearch_Health(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	ctx := context.Background()
	client, cleanup := startMeilisearch(ctx, t)
	defer cleanup()

	resp, err := client.Health()
	require.NoError(t, err)
	assert.Equal(t, "available", resp.Status)
}

func TestIntegration_Meilisearch_CreateIndex(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	ctx := context.Background()
	client, cleanup := startMeilisearch(ctx, t)
	defer cleanup()

	task, err := client.CreateIndex(&meilisearch.IndexConfig{Uid: "obm_test_vmp"})
	require.NoError(t, err)
	assert.NotEmpty(t, task.TaskUID)

	_, err = client.Index("obm_test_vmp").AddDocuments([]map[string]interface{}{
		{"id": 1, "co_seq_id": 1, "no_nm": "Ibuprofen 200mg"},
		{"id": 2, "co_seq_id": 2, "no_nm": "Paracetamol 500mg"},
	}, nil)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	searchResp, err := client.Index("obm_test_vmp").Search("Ibuprofen", &meilisearch.SearchRequest{})
	require.NoError(t, err)
	assert.Greater(t, searchResp.TotalHits, int64(0))
}
