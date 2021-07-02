package repository

import "oauthServer/entity"

type ClientRepository interface {
	Save(client *entity.Client) (*entity.Client, map[string]string)
}
