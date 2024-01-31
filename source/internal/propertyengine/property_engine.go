package propertyengine

import (
	"time"

	"github.com/LearningMotors/nms-property-engine/source/engine"
	"github.com/LearningMotors/platform/service/platenv"
)

var propertyEngine engine.Engine

func Initialize() (err error) {
	const (
		envConfigRoot      = "NMS_PROPERTY_ENGINE_CONFIG_ROOT"
		envExpiry          = "NMS_PROPERTY_ENGINE_CACHE_EXPIRY_SECONDS"
		envCleanupInterval = "NMS_PROPERTY_ENGINE_CACHE_CLEANUP_SECONDS"

		defaultConfigRoot      = "nms-benthos-yamls/"
		defaultExpiry          = 1
		defaultCleanupInterval = 1
	)

	configRoot := platenv.GetEnvWithDefaultAsString(envConfigRoot, defaultConfigRoot)

	expiry := platenv.GetEnvWithDefaultAsInt(envExpiry, defaultExpiry)
	expiryDuration := time.Duration(expiry) * time.Second

	cleanupInterval := platenv.GetEnvWithDefaultAsInt(envCleanupInterval, defaultCleanupInterval)
	cleanupIntervalDuration := time.Duration(cleanupInterval) * time.Second

	propertyEngine, err = engine.NewCachedEngine(configRoot, expiryDuration, cleanupIntervalDuration)

	return
}

func Get() engine.Engine {
	return propertyEngine
}
