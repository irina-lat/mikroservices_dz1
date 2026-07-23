//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	pb "shared/pkg/proto/inventory/v1"
)

type InventoryIntegrationTestSuite struct {
	suite.Suite
	env    *TestEnvironment
	client pb.InventoryServiceClient
	teardown func()
}

func (s *InventoryIntegrationTestSuite) SetupSuite() {
	env, client, teardown := SetupTest(s.T())
	s.env = env
	s.client = client
	s.teardown = teardown
}

func (s *InventoryIntegrationTestSuite) TearDownSuite() {
	if s.teardown != nil {
		s.teardown()
	}
}

func (s *InventoryIntegrationTestSuite) TestListParts_ReturnsAllParts() {
	ctx := context.Background()

	resp, err := s.client.ListParts(ctx, &pb.ListPartsRequest{
		Filter: &pb.PartsFilter{},
	})
	require.NoError(s.T(), err)
	require.NotNil(s.T(), resp)
	require.Greater(s.T(), len(resp.Parts), 0)
}

func (s *InventoryIntegrationTestSuite) TestGetPart_ByUUID() {
	ctx := context.Background()

	// сначала получаем список
	listResp, err := s.client.ListParts(ctx, &pb.ListPartsRequest{
		Filter: &pb.PartsFilter{},
	})
	require.NoError(s.T(), err)
	require.Greater(s.T(), len(listResp.Parts), 0)

	firstUUID := listResp.Parts[0].Uuid

	resp, err := s.client.GetPart(ctx, &pb.GetPartRequest{
		Uuid: firstUUID,
	})
	require.NoError(s.T(), err)
	require.NotNil(s.T(), resp)
	require.NotNil(s.T(), resp.Part)
	require.Equal(s.T(), firstUUID, resp.Part.Uuid)
}

func (s *InventoryIntegrationTestSuite) TestListParts_FilterByCategory() {
	ctx := context.Background()

	resp, err := s.client.ListParts(ctx, &pb.ListPartsRequest{
		Filter: &pb.PartsFilter{
			Categories: []pb.Category{pb.Category_CATEGORY_ENGINE},
		},
	})
	require.NoError(s.T(), err)
	require.NotNil(s.T(), resp)

	// проверяем, что все детали в ответе имеют категорию ENGINE
	for _, part := range resp.Parts {
		require.Equal(s.T(), pb.Category_CATEGORY_ENGINE, part.Category)
	}
}

func (s *InventoryIntegrationTestSuite) TestGetPart_NotFound() {
	ctx := context.Background()

	resp, err := s.client.GetPart(ctx, &pb.GetPartRequest{
		Uuid: "non-existent-uuid",
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), resp)
}

func TestInventoryIntegration(t *testing.T) {
	suite.Run(t, new(InventoryIntegrationTestSuite))
}