package controller

import (
	"sync"

	"abmp.cc/app/pkg/app"
	"github.com/abmpio/entity"
	"github.com/abmpio/mongodbr"
)

var (
	_serviceFactoryInstanceOnce sync.Once
)

func getEntityService[T mongodbr.IEntity]() entity.IEntityService[T] {
	return app.Context.GetInstance(new(entity.IEntityService[T])).(entity.IEntityService[T])
}
