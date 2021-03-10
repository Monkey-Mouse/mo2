# 设计手册

## directory
``` go
// Directory 目录
type Directory struct {
	ID       primitive.ObjectID   `json:"id,omitempty" example:"xxxxxxxxxxxxxx==" bson:"_id,omitempty"`
	ParentID primitive.ObjectID   `json:"parent_id,omitempty" example:"xxxxxxxxxxxxxx==" bson:"parent_id,omitempty"`
	Name     string               `json:"name,omitempty" example:"records" bson:"name,omitempty"`
	OwnerIDs []primitive.ObjectID `json:"owner_ids,omitempty"  bson:"owner_ids,omitempty"`
}
```
### 字段含义
- ParentID:父目录的id
- Name：名称
- OwnerIDs：归属者的id列表，用于访问控制

### 设计思路
需要层级目录的对象都可以复用，比如收藏夹，群组共享等  
将应用于不同对象的数据存放在不同的表中，下面均以第一个实现的category为例：  
首先，用户初始化category功能，生成一个category的对象
