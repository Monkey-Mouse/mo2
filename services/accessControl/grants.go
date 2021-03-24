package accessControl

import (
	"github.com/Monkey-Mouse/go-abac/abac"
	_ "github.com/Monkey-Mouse/go-abac/abac"
)

var Ctrl = abac.AccessControl{}

func init() {
	Ctrl.Grant(abac.GrantsType{"account": {
		"blog": {
			abac.ActionUpdate: []abac.RuleType{&AllowOwn{}},
			"create:any":      []abac.RuleType{},
			"read:own":        abac.RulesType{},
		},
	}})

}
