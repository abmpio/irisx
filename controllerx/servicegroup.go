package controllerx

import (
	"abmp.cc/app/pkg/app"
	"github.com/abmpio/entity"
	"github.com/abmpio/mongodbr"
)

func GetEntityService[T mongodbr.IEntity]() entity.IEntityService[T] {
	return app.Context.GetInstance(new(entity.IEntityService[T])).(entity.IEntityService[T])
}
