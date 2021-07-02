package repository

import "oauthServer/entity"

type ScopeRepository interface {
	Save(scope *entity.Scope) (*entity.Scope, map[string]string)
}
