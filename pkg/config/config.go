package config

import "time"

type (
	// Config stores the configuration settings.
	Config struct {

		// Disable filesystem polling, used by hub
		DisableDirtyPolling bool `default:"false" envconfig:"DISABLE_DIRTY_POLLING"`
		DisableFlexVolume   bool `default:"false" envconfig:"DISABLE_FLEXVOLUME"`

		DotmeshUpgradesURL string

		PollDirty struct {
			SuccessTimeout time.Duration `default:"1s" envconfig:"POLL_DIRTY_SUCCESS_TIMEOUT"`
			ErrorTimeout   time.Duration `default:"1s" envconfig:"POLL_DIRTY_ERROR_TIMEOUT"`
		}

		Upgrades struct {
			URL             string `envconfig:"DOTMESH_UPGRADES_URL"`
			IntervalSeconds int    `default:"300" envconfig:"DOTMESH_UPGRADES_INTERVAL_SECONDS"`
		}
	}
)
