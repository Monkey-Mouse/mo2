# 设计手册


## [directory](directory.go)
``` go
// file: directory.go
// DirectoryInfo 目录信息
type DirectoryInfo struct {
	Description string `json:"description,omitempty" example:"course materials" bson:"description,omitempty"`
	Cover       string `json:"cover,omitempty" example:"https://www.motwo.cn/cover" bson:"cover,omitempty"`
}

// Directory 目录
type Directory struct {
	ID       primitive.ObjectID   `json:"id,omitempty" example:"xxxxxxxxxxxxxx==" bson:"_id,omitempty"`
	ParentID primitive.ObjectID   `json:"parent_id,omitempty" example:"xxxxxxxxxxxxxx==" bson:"parent_id,omitempty"`
	Name     string               `json:"name,omitempty" example:"records" bson:"name,omitempty"`
	Info     DirectoryInfo        `json:"info,omitempty" bson:"info,omitempty"`
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
1. 首先，用户初始化category功能，生成一个category的对象（定义名称为`root`）
   - parent_id为创建用户id
   - 此时的ownerIDs也添加用户id
2. 若为category新增子category，对子归档(subCategory)进行修改：
   - parent_id为父category的id
   - 此时的ownerIDs也添加用户id
    
这种设计的优势有：
- 易得某目录c的子目录s(c.ID==s.ParentID)
- 易得子目录s的父目录p(p.ID==s.ParentID)
- 易得某用户u的根目录r(u.ID==r.ParentID)
- 易得某用户u的所有目录cs(u.ID in cs.OwnerIDs)
    
### api设计
由于新增了归档这一实体，其内部之间的关系、归档与博客之间的关系、归档与用户之间的关系都需要设定与查询。
与此相类，之后复用目录结构的实体也将面对这一需求，因此，设计泛化通用的api势在必行   
遵循REST API的设计理念，将api的路径对应相应资源，方法代表操作，于是引入以下实体：


#### [dto/directory.go](../../dto/directory.go)
```go
// file: dto/directory.go

// RelateEntity2Entity 将单实体关联到单实体dto
type RelateEntity2Entity struct {
	RelatedID  primitive.ObjectID `json:"related_id,omitempty"`
	RelateToID primitive.ObjectID `json:"relateTo_id,omitempty"`
}

// RelateEntitySet2EntitySet 关联两个实体集dto
type RelateEntitySet2EntitySet struct {
	RelatedIDs  []primitive.ObjectID `json:"related_ids,omitempty"`
	RelateToIDs []primitive.ObjectID `json:"relateTo_ids,omitempty"`
}

// RelateEntity2EntitySet 关联单实体到多实体集dto
type RelateEntity2EntitySet struct {
	RelatedID   primitive.ObjectID   `json:"related_id,omitempty"`
	RelateToIDs []primitive.ObjectID `json:"relateTo_ids,omitempty"`
}

// RelateEntitySet2Entity 关联实体集到单实体dto
type RelateEntitySet2Entity struct {
	RelatedIDs []primitive.ObjectID `json:"related_ids,omitempty"`
	RelateToID primitive.ObjectID   `json:"relateTo_id,omitempty"`
}
```

分别对应实现：
- 将单实体关联到单实体
- 关联两个实体集
- 关联单实体到多实体集
- 关联实体集到单实体

### api使用
[controller.go](../controller/controller.go)
```go
    api := middleware.H.Group("/api")
    {
    	//...
        relation := api.Group("relation", model.OrdinaryUser)
        {
            relation.Post("categories/:type", c.RelateCategories2Entity)
            relation.Post("category/:type", c.RelateCategory2Entity)
            relation.Get("category/:type/:ID", c.FindCategoriesByType)
        }
        //...
    }
       
```
* 未实现前加(x)
- 基础
    - 访问控制未完善
    - api/blog/category [post]
        - upsert增改category信息
    - api/blog/category [get]
        - 查询category信息
  - (x) api/blog/category [delete] (important)
      - 删除category信息
      - 解除与该category相关的所有联系：blog字段寻找相关字段并删去

 
- relation部分
    - categories
      - 多对一
        - api/relation/categories/:type [post]
          - 建立categories与type之间的联系，目前可选：
          - category:将多实体集categories的父categoryid均设为单实体的id
          - blog:将多实体集categories的ids添加到blog的categories列表中去
        
      - (x)多对多
    - category
      - 一对一
        - api/relation/category/:type [post]
          - 建立category与type之间的联系，目前可选：
          - user:将user的id添加到category的ownerIDs列表中
          - userMain:将user的id设定为category的parent_id
          - category:将related id的parent_id设定为relateTo的category的id
          - blog:将category的id添加到blog的categories列表中去
      - 一对一/多  
        - api/relation/category/:type/:id [get]
          - 获取categories与type之间的联系，目前可选：
          - user:将多实体集categories的父categoryid均设为单实体的id
          - sub:将多实体集categories添加到blog的categories列表中去
      
    
    
#### 删除部分
因为删除某directory涉及到冗余关联数据的删除，因此需要：
- 对blog端：删除所有blog中categories列表里的相关id即可
- 对category端：所有删除category(`dCat`)的子category加入到其上一级(`id=dCat.parent_id`)中
- 增加鉴权，只有操作用户id in owner_ids匹配成功的可以进行删除操作
    - 新思路，增加过滤器，在请求的id列表中过滤出可以进行操作的id列表






