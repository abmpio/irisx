package controllerx

import (
	"github.com/abmpio/entity"
	"github.com/abmpio/webserver/controller"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SearchInput struct {
	controller.Pagination
	Filter map[string]interface{} `json:",inline"`

	SortInput
}

// 获取指定的key的primitive.ObjectID值
func (i *SearchInput) GetFilterValueAsObjectId(key string) (*primitive.ObjectID, error) {
	if len(i.Filter) <= 0 {
		return nil, nil
	}
	v, ok := i.Filter[key].(string)
	if !ok {
		return nil, nil
	}
	if len(v) <= 0 {
		return nil, nil
	}
	oid, err := primitive.ObjectIDFromHex(v)
	if err != nil {
		return nil, err
	}
	return &oid, nil
}

// 获取指定的key的string值
func (i *SearchInput) GetFilterValueAsString(key string) string {
	if len(i.Filter) <= 0 {
		return ""
	}
	v, ok := i.Filter[key].(string)
	if ok {
		return v
	}
	return ""
}

// 确保Filter属性不为nil
func (i *SearchInput) EnsureFilterNotNil() {
	if i.Filter == nil {
		i.Filter = make(map[string]interface{})
	}
}

func (i *SearchInput) ReadFromQuery(ctx iris.Context) error {
	err := ctx.ReadQuery(&i.Pagination)
	if err != nil {
		return err
	}
	err = ctx.ReadQuery(&i.SortInput)
	if err != nil {
		return err
	}
	return nil
}

type SortInput struct {
	Sorts []entity.Sort `json:"sorts" url:"sorts"`
}
