package utils

//import (
//	"github.com/0987363/mgo/bson"
//	"time"
//)
//
//type Platform struct {
//	ID          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
//	CreatedAt   *time.Time    `json:"created_at,omitempty" bson:"created_at,omitempty"`
//	UpdatedAt   *time.Time    `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
//	Name        string        `json:"name,omitempty" bson:"name,omitempty"`               // 名称
//	TypeID      uint32        `json:"type_id,omitempty" bson:"type_id,omitempty"`         // 类别id
//	ClientId    uint32        `json:"client_id,omitempty" bson:"client_id,omitempty"`     // 客户id
//	Secret      string        `json:"secret,omitempty" bson:"secret,omitempty"`           // 密钥
//	Status      uint32        `json:"status" bson:"status"`                               // 运行状态
//	Description string        `json:"description,omitempty" bson:"description,omitempty"` // 备注
//}
//
//// 产品
//type Product struct {
//	Status           uint32          `json:"status,omitempty"  bson:"status,omitempty"`    // 产品状态：0-EOL,1-开发中，2-公开
//	LimitPlatformsID []bson.ObjectId `json:"limit_platforms_id" bson:"limit_platforms_id"` // 平台限制
//}
