package obswatcher

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infrastructure/obsws"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infrastructure/sharedmem"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/utils"
)

const (
	componentName          = "obswatcher"
	syncPeriod             = 10 * time.Second
	startPreparetionPeriod = 60 * time.Second
)

type Config struct {
	Development bool
	Debug       bool
	Obs         []ConfigObs
}

type ConfigObs struct {
	Host      string
	Password  string
	DkTrackId int32
}

func Run(ctx context.Context, conf Config) error {
	// setup logger
	zapConf := zap.NewProductionConfig()
	if conf.Development {
		zapConf = zap.NewDevelopmentConfig()
	}
	zapConf.DisableStacktrace = true // due to output wrapped error in errorVerbose
	zapLogger, err := zapConf.Build()
	if err != nil {
		return err
	}
	logger := zapr.NewLogger(zapLogger).WithName(componentName)
	ctx = logr.NewContext(ctx, logger)

	mr := &sharedmem.Reader{UseStorageForDisableAutomation: true}

	eg, ctx := errgroup.WithContext(ctx)
	for _, obs := range conf.Obs {
		obswsClient, err := obsws.NewObsWebSocketClient(obs.Host, obs.Password)
		if err != nil {
			err := xerrors.Errorf("message: %w", err)
			logger.Error(err, "obsws.NewObsWebSocketClient() was failed")
			return err
		}
		eg.Go(watch(ctx, obs.DkTrackId, obswsClient, mr))
	}
	if err := eg.Wait(); err != nil {
		err := xerrors.Errorf("message: %w", err)
		logger.Error(err, "eg.Wait() was failed")
		return err
	}
	return nil
}

func watch(ctx context.Context, trackId int32,
	obswsClient obsws.ClientIface, mr sharedmem.ReaderIface,
) func() error {
	return func() error {
		logger := utils.GetLogger(ctx).WithValues("trackId", trackId)
		tick := time.NewTicker(syncPeriod)

		for {
			select {
			case <-ctx.Done():
				logger.Info("context was done.")
				return nil
			case <-tick.C:
				if err := procedure(ctx, trackId, obswsClient, mr); err != nil {
					return xerrors.Errorf("message: %w", err)
				}
			}
		}
	}
}

func procedure(ctx context.Context, trackId int32,
	obswsClient obsws.ClientIface, mr sharedmem.ReaderIface,
) error {
	logger := utils.GetLogger(ctx).WithValues("trackId", trackId)

	if disabled, err := mr.DisableAutomation(trackId); err != nil {
		logger.Error(xerrors.Errorf("message: %w", err), "mr.DisableAutomation() was failed")
		return nil
	} else if disabled {
		logger.Info("DisableAutomation was true, skipped")
		return nil
	}

	t, err := obswsClient.GetRemainingTimeOnCurrentScene(ctx)
	if err != nil {
		logger.Error(xerrors.Errorf("message: %w", err), "obswsClient.GetRemainingTimeOnCurrentScene() was failed")
		return nil
	}
	remainingMilliSecond := t.DurationMilliSecond - t.CursorMilliSecond

	if float64(startPreparetionPeriod/time.Millisecond) < remainingMilliSecond {
		logger.Info(fmt.Sprintf("remainingTime on current Scene's MediaInput is %ds: continue",
			startPreparetionPeriod/time.Second))
		return nil
	}
	logger.Info(fmt.Sprintf("remainingTime on current Scene's MediaInput is within %ds",
		startPreparetionPeriod/time.Second), "duration", t.DurationMilliSecond, "cursor", t.CursorMilliSecond)

	// sleep until MediaInput is finished
	time.Sleep(time.Duration(remainingMilliSecond) * time.Millisecond)
	if err := obswsClient.MoveSceneToNext(context.Background()); err != nil {
		logger.Error(xerrors.Errorf("message: %w", err), "obswsClient.MoveSceneToNext() on automated task was failed")
		return nil
	}
	logger.Info("automated task was completed. Scene should be to next.")
	return nil
}
