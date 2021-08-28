package repository

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"oauthServer/entity"
	"oauthServer/pkg/logger"
	"sync"
	"time"
)

var (
	address      = "root:123456@tcp(127.0.0.1:3306)/oauth-server?charset=utf8&parseTime=True&loc=Local"
	addr         = "localhost:6379"
	once         sync.Once
	repositories *Repositories
)

// Repositories 持久化
type Repositories struct {
	//mysql 实现
	UserRepository       UserRepository
	RoleRepository       RoleRepository
	PermissionRepository PermissionRepository
	//redis 实现
}

// NewRepositories 生成 Go 持久化api
func NewRepositories() *Repositories {
	once.Do(func() {
		//sqlDB.LogMode(true)
		db := newGormClient()
		//rdb := newRedisClient()
		repositories = &Repositories{
			UserRepository:       NewUserRepository(db),
			RoleRepository:       NewRoleRepository(db),
			PermissionRepository: NewPermissionRepository(db),
			//UserCacheRepository: repository.NewUserCacheRepository(rdb),
			//db: db,
		}
	})
	return repositories
}

func newGormClient() *gorm.DB {
	db, err := gorm.Open(mysql.Open(address), &gorm.Config{
		Logger: logger.GormConfig(),
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix: "t_",   // 表名前缀，`User`表为`t_users`
			SingularTable: true, // 使用单数表名，启用该选项后，`User` 表将是`user`
			//NameReplacer: strings.NewReplacer("CID", "Cid"), // 在转为数据库名称之前，使用NameReplacer更改结构/字段名称。
		},
	})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	err = db.AutoMigrate(
		entity.User{},
		entity.Permission{},
		entity.Client{},
		entity.Role{},
		entity.Scope{},
	)

	if err != nil {
		panic("failed to migrate the schema")
	}

	sqlDB, err := db.DB()

	if err != nil {
		panic(err)
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	//db.WithContext(ctx)

	return db
	// Output: key value
	// key2 does not exist
}

func newRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       14, // use default DB
	})
	//
	//mycache := cache.New(&cache.Options{
	//	Redis: rdb,
	//})

	return rdb
	//err := rdb.Set(ctx, "key", entity.User{}, 0).Err()
	//if err != nil {
	//	panic(err)
	//}
	//
	//val, err := rdb.Get(ctx, "key").Result()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("key", val)
	//
	//val2, err := rdb.Get(ctx, "key2").Result()
	//if err == redis.Nil {
	//	fmt.Println("key2 does not exist")
	//} else if err != nil {
	//	panic(err)
	//} else {
	//	fmt.Println("key2", val2)
	//}
	// Output: key value
	// key2 does not exist
}
