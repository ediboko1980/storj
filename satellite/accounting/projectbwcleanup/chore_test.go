// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package projectbwcleanup_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"storj.io/common/pb"
	"storj.io/common/testcontext"
	"storj.io/common/testrand"
	"storj.io/storj/private/testplanet"
	"storj.io/storj/satellite"
)

var testBytes int64 = 5000

func TestProjectAllocatedBandwidthRetainNegative(t *testing.T) {
	// -1 to validate we don't delete the current month
	testProjectAllocatedBandwidthRetain(t, -1)
}

func TestProjectAllocatedBandwidthRetainZero(t *testing.T) {
	// 0 to validate we don't delete the current month
	testProjectAllocatedBandwidthRetain(t, 0)
}

func TestProjectAllocatedBandwidthRetainTwo(t *testing.T) {
	testProjectAllocatedBandwidthRetain(t, 2)
}

func testProjectAllocatedBandwidthRetain(t *testing.T, retain int) {
	testplanet.Run(t, testplanet.Config{
		SatelliteCount: 1, StorageNodeCount: 0, UplinkCount: 0,
		Reconfigure: testplanet.Reconfigure{
			Satellite: func(log *zap.Logger, index int, config *satellite.Config) {
				config.ProjectBWCleanup.RetainMonths = retain
			},
		},
	}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
		months := retain + 1 // number of months to retain + 1 for current month
		// must have at least have the current month
		if months < 1 {
			months = 1
		}

		satellite := planet.Satellites[0]
		satellite.Accounting.ProjectBWCleanup.Loop.Pause()
		ordersDB := satellite.DB.Orders()

		projectID := testrand.UUID()
		bucketName := testrand.BucketName()
		now := time.Now()
		interval := time.Date(now.Year(), now.Month(), 15, 12, 0, 0, 0, time.UTC)

		for i := 0; i <= months; i++ {
			err := ordersDB.UpdateBucketBandwidthAllocation(ctx, projectID, []byte(bucketName), pb.PieceAction_GET, testBytes, interval.AddDate(0, -i, 0))
			require.NoError(t, err)
		}
		for i := 0; i <= months; i++ {
			m := interval.AddDate(0, -i, 0)
			bytes, err := satellite.Accounting.ProjectUsage.GetProjectAllocatedBandwidth(ctx, projectID, m.Year(), m.Month())
			require.NoError(t, err)
			require.EqualValues(t, testBytes, bytes)
		}

		satellite.Accounting.ProjectBWCleanup.Loop.TriggerWait()

		for i := 0; i <= months; i++ {
			m := interval.AddDate(0, -i, 0)
			bytes, err := satellite.Accounting.ProjectUsage.GetProjectAllocatedBandwidth(ctx, projectID, m.Year(), m.Month())

			if i < months || retain < 0 { // there should always be the current month
				require.NoError(t, err)
				require.EqualValues(t, testBytes, bytes)
			} else {
				require.NoError(t, err)
				require.EqualValues(t, 0, bytes)
			}

		}
	})
}
