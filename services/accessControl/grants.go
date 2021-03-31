package accessControl

import (
	"github.com/Monkey-Mouse/go-abac/abac"
	_ "github.com/Monkey-Mouse/go-abac/abac"
)

var Ctrl = abac.AccessControl{}

func init() {
	Ctrl.Grant(abac.GrantsType{"account": {
		"blog": {
			abac.ActionUpdate: []abac.RuleType{&AllowOwn{}, &AccessFilter{}},
			"create:any":      []abac.RuleType{},
			"read:own":        abac.RulesType{},
		},
		"group": {
			abac.ActionCreate: []abac.RuleType{},
			abac.ActionDelete: []abac.RuleType{&AllowOwn{}, &AccessFilter{}},
			abac.ActionUpdate: []abac.RuleType{&AllowOwn{}, &AccessFilter{}},
			abac.ActionRead:   []abac.RuleType{},
		},
	}})

}
