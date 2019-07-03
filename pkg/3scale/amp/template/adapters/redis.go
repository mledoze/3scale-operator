package adapters

import (
	"github.com/3scale/3scale-operator/pkg/3scale/amp/component"
	"github.com/3scale/3scale-operator/pkg/3scale/amp/product"
	"github.com/3scale/3scale-operator/pkg/common"
	templatev1 "github.com/openshift/api/template/v1"
)

type RedisAdapter struct {
}

func NewRedisAdapter(options []string) Adapter {
	return NewAppenderAdapter(&RedisAdapter{})
}

func (a *RedisAdapter) Parameters() []templatev1.Parameter {
	productVersion := product.CurrentProductVersion()
	imageProvider, err := product.NewImageProvider(productVersion)
	if err != nil {
		panic(err)
	}
	return []templatev1.Parameter{
		{
			Name:        "REDIS_IMAGE",
			Description: "Redis image to use",
			Required:    true,
			// We use backend-redis image because we have to choose one
			// but in templates there's no distinction between Backend Redis image
			// used and System Redis image. They are always the same
			Value: imageProvider.GetBackendRedisImage(),
		},
	}
}

func (r *RedisAdapter) Objects() ([]common.KubernetesObject, error) {
	redisOptions, err := r.options()
	if err != nil {
		return nil, err
	}
	redisComponent := component.NewRedis(redisOptions)
	return redisComponent.Objects(), nil
}

func (r *RedisAdapter) options() (*component.RedisOptions, error) {
	rob := component.RedisOptionsBuilder{}
	rob.AppLabel("${APP_LABEL}")

	return rob.Build()
}
