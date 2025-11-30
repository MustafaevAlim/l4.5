package controllers

import "l4.5/internal/repository"

// В принципе можно использовать только кеш, бд на будущее пусть будет
type Controller struct {
	DB    *repository.Storage
	Cache *repository.LRUcache
}
