package sqlstash

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/store/entity"
	"github.com/grafana/grafana/pkg/services/store/entity/db/dbimpl"
	"github.com/grafana/grafana/pkg/setting"
)

func TestCreate(t *testing.T) {
	sqlStore := db.InitTestDB(t)

	entityDB, err := dbimpl.ProvideEntityDB(
		sqlStore,
		setting.NewCfg(),
		featuremgmt.WithFeatures(featuremgmt.FlagUnifiedStorage))
	require.NoError(t, err)

	s, err := ProvideSQLEntityServer(entityDB)
	require.NoError(t, err)

	tests := []struct {
		name             string
		ent              *entity.Entity
		errIsExpected    bool
		statusIsExpected bool
	}{
		{
			"request with key and entity creator",
			&entity.Entity{
				Group:     "playlist.grafana.app",
				Resource:  "playlists",
				Namespace: "default",
				Name:      "set-minimum-uid",
				Key:       "/playlist.grafana.app/playlists/default/set-minimum-uid",
				CreatedBy: "set-minimum-creator",
			},
			false,
			true,
		},
		{
			"request with no entity creator",
			&entity.Entity{
				Key: "/playlist.grafana.app/playlists/default/set-only-key",
			},
			true,
			false,
		},
		{
			"request with no key",
			&entity.Entity{
				CreatedBy: "entity-creator",
			},
			true,
			true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := entity.CreateEntityRequest{
				Entity: &entity.Entity{
					Key:       tc.ent.Key,
					CreatedBy: tc.ent.CreatedBy,
				},
			}
			resp, err := s.Create(context.Background(), &req)

			if tc.errIsExpected {
				require.Error(t, err)

				if tc.statusIsExpected {
					require.Equal(t, entity.CreateEntityResponse_ERROR, resp.Status)
				}

				return
			}

			require.Nil(t, err)
			require.Equal(t, entity.CreateEntityResponse_CREATED, resp.Status)
			require.NotNil(t, resp)
			require.Nil(t, resp.Error)

			read, err := s.Read(context.Background(), &entity.ReadEntityRequest{
				Key: tc.ent.Key,
			})
			require.NoError(t, err)
			require.NotNil(t, read)

			require.Greater(t, len(read.Guid), 0)
			require.Greater(t, read.ResourceVersion, int64(0))

			expectedETag := createContentsHash(tc.ent.Body, tc.ent.Meta, tc.ent.Status)
			require.Equal(t, expectedETag, read.ETag)
			require.Equal(t, tc.ent.Origin, read.Origin)
			require.Equal(t, tc.ent.Group, read.Group)
			require.Equal(t, tc.ent.Resource, read.Resource)
			require.Equal(t, tc.ent.Namespace, read.Namespace)
			require.Equal(t, tc.ent.Name, read.Name)
			require.Equal(t, tc.ent.Subresource, read.Subresource)
			require.Equal(t, tc.ent.GroupVersion, read.GroupVersion)
			require.Equal(t, tc.ent.Key, read.Key)
			require.Equal(t, tc.ent.Folder, read.Folder)
			require.Equal(t, tc.ent.Meta, read.Meta)
			require.Equal(t, tc.ent.Body, read.Body)
			require.Equal(t, tc.ent.Status, read.Status)
			require.Equal(t, tc.ent.Title, read.Title)
			require.Equal(t, tc.ent.Size, read.Size)
			require.Equal(t, tc.ent.CreatedAt, read.CreatedAt)
			require.Equal(t, tc.ent.CreatedBy, read.CreatedBy)
			require.Equal(t, tc.ent.UpdatedAt, read.UpdatedAt)
			require.Equal(t, tc.ent.UpdatedBy, read.UpdatedBy)
			require.Equal(t, tc.ent.Description, read.Description)
			require.Equal(t, tc.ent.Slug, read.Slug)
			require.Equal(t, tc.ent.Message, read.Message)
			require.Equal(t, tc.ent.Labels, read.Labels)
			require.Equal(t, tc.ent.Fields, read.Fields)
			require.Equal(t, tc.ent.Errors, read.Errors)

			// check entity_history

			sess, err := entityDB.GetSession()
			require.NoError(t, err)

			// #TODO update the query used here
			rows, err := sess.Query(context.Background(), "SELECT * FROM entity_history;")
			require.NoError(t, err)

			defer func() { _ = rows.Close() }()

			ss, ok := s.(*sqlEntityServer)
			require.True(t, ok)

			require.True(t, rows.Next())

			// #TODO figure out what the request should contain
			ent, err := ss.rowToEntity(context.Background(), rows, &entity.ReadEntityRequest{
				Key: tc.ent.Key,
			})
			fmt.Printf("%+v", ent)
			require.Equal(t, tc.ent.Key, ent.Key) // #TODO figure out why it panics here

		})
	}
}

// #TODO for later can also return other vals besides entity.EntityStoreServer if we need them
// #TODO we need access to the entityDB in the test. Refactor set up.
// func setUpTestServer(t *testing.T) entity.EntityStoreServer {
// 	sqlStore := db.InitTestDB(t)

// 	entityDB, err := dbimpl.ProvideEntityDB(
// 		sqlStore,
// 		setting.NewCfg(),
// 		featuremgmt.WithFeatures(featuremgmt.FlagUnifiedStorage))
// 	require.NoError(t, err)

// 	s, err := ProvideSQLEntityServer(entityDB)
// 	require.NoError(t, err)
// 	return s
// }
