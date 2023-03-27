

# 字节跳动青训营项目-抖音后端







## 接口实现

-  基础接口
-  扩展接口-I
-  扩展接口-II

![image-20230326225532499](https://raw.githubusercontent.com/Doraemonkeys/picture/master/1/image-20230326225532499.png)

## 演示截图

![图叨叨_20230326_231558](https://raw.githubusercontent.com/Doraemonkeys/picture/master/1/%E5%9B%BE%E5%8F%A8%E5%8F%A8_20230326_231558.jpg)



![图叨叨_20230326_231844](https://raw.githubusercontent.com/Doraemonkeys/picture/master/1/%E5%9B%BE%E5%8F%A8%E5%8F%A8_20230326_231844.jpg)



![图叨叨_20230326_231912](https://raw.githubusercontent.com/Doraemonkeys/picture/master/1/%E5%9B%BE%E5%8F%A8%E5%8F%A8_20230326_231912.jpg)



## 技术栈

数据库：MySQL

缓存：LRU、ARC

GO框架：Gin、Gorm、logrus

消息队列： SimpleMQ

鉴权：JWT+AES

加密：bcrypt

对象存储：

获取封面：ffmpeg

## 使用方法

- Go编译器版本要求 ：1.20 及以上

```bash
git clone https://github.com/Doraemonkeys/douyin2.git
```

1. 安装mysql和ffmpeg（用于上传视频后的处理）。
2. 编译执行config/cmd下的main.go文件，将config/conf下的`example.yaml`改名为`config.yaml`。
3. 修改配置文件`config.yaml。`。
   domain配置项用于上传视频后生成的`play_url`与`cover_url` 注意将域名解析到后端所监听的IP。
   mysql相关配置只需要建立数据库并分配用户权限 数据表会在首次启动时自动生成。
4. 项目根目录执行go build即可生成可执行文件。





## 目录结构

```go
├─config
│  │  config.go
│  │  type.go
│  ├─cmd
│  │      main.go
│  └─conf
│          config.yaml
│          example.yaml
├─initiate
│      init.go
├─internal
│  ├─app
│  │  │  common.go
│  │  ├─handlers
│  │  │  ├─comment
│  │  │  │      comment.go
│  │  │  ├─favorite
│  │  │  │      favorite.go
│  │  │  ├─feed
│  │  │  │      vedio.go
│  │  │  ├─follow
│  │  │  │      follow.go
│  │  │  ├─publish
│  │  │  │      publish.go
│  │  │  ├─response
│  │  │  │      comment.go
│  │  │  │      common.go
│  │  │  │      favorite.go
│  │  │  │      feed.go
│  │  │  │      login.go
│  │  │  │      publish.go
│  │  │  │      register.go
│  │  │  │      user.go
│  │  │  └─user
│  │  │          register.go
│  │  │          user.go
│  │  ├─middleware
│  │  │      jwt.go
│  │  │      login.go
│  │  ├─models
│  │  │      collection.go
│  │  │      comment.go
│  │  │      follow.go
│  │  │      like.go
│  │  │      user.go
│  │  │      vedio.go
│  │  └─services
│  │          comment.go
│  │          follow.go
│  │          register.go
│  │          user.go
│  │          vedio.go
│  ├─database
│  │      cache.go
│  │      mysql.go
│  │      redis.go
│  │      storage.go
│  ├─msgQueue
│  │      comment.go
│  │      favorite.go
│  │      follow.go
│  ├─pkg
│  │  ├─cache
│  │  │      arc.go
│  │  │      cache.go
│  │  ├─messageQueue
│  │  │      simpleMQ.go
│  │  │      simpleMQ_test.go
│  │  │      type.go
│  │  └─storage
│  │          interface.go
│  │          loacal.go
│  └─server
│          server.go
├─monitor
│      system.go
├─pkg
│  ├─jwt
│  │      jwt.go
│  ├─log
│  │      formatter.go
│  │      log.go
└─utils
        crypto.go
        crypto_test.go
        file.go
        password.go
        password_test.go
        string.go
```





## 数据库表设计

### ER图

![2023-03-26_232448](https://raw.githubusercontent.com/Doraemonkeys/picture/master/1/2023-03-26_232448.png)





使用gorm自动建表

```go
func mirateTable() {
	db.SetupJoinTable(&models.UserModel{}, models.UserModelTable_FollowersSlice, &models.UserFollowerModel{})
	db.SetupJoinTable(&models.UserModel{}, models.UserModelTable_FansSlice, &models.UserFollowerModel{})
	db.SetupJoinTable(&models.UserModel{}, models.UserModelTable_LikesSlice, &models.UserLikeModel{})
	db.SetupJoinTable(&models.UserModel{}, models.UserModelTable_CollectionsSlice, &models.UserCollectionModel{})

	db.SetupJoinTable(&models.VideoModel{}, models.VideoModelTable_LikesSlice, &models.UserLikeModel{})
	db.SetupJoinTable(&models.VideoModel{}, models.VideoModelTable_CollectionsSlice, &models.UserCollectionModel{})

	db.AutoMigrate(
		&models.UserModel{},
		&models.VideoModel{},
		&models.CommentModel{},
		&models.UserFollowerModel{},
		&models.UserLikeModel{},
		&models.UserCollectionModel{},
	)
}
```



## 架构设计

项目通过引入消息队列和缓存，减轻了数据库的负载，提高系统的性能和可扩展性。其次，通过使用JWT鉴权，增强系统的安全性，并防止未经授权的访问。最后，通过分层设计实现了解耦，将不同的功能模块分离到不同的层中，可以使系统更易于维护和扩展。



### 处理流程

- 客户端 -> middleware -> handler -> service -> database
- 客户端 -> middleware -> handler -> service -> message queue -> database
- 客户端 -> middleware -> handler -> service -> cache -> database

![image-20230327014620713](https://raw.githubusercontent.com/Doraemonkeys/picture/master/1/image-20230327014620713.png)





### 鉴权

1. 颁发token

JWT 默认不加密，为了防止用户信息的泄露，本项目使用AES算法对JWT原始Token进行加密。

```go
func (j *CryptJWT) CreateToken(claims CustomClaims) (string, error) {
	jwTtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwTtoken.SignedString(j.signingKey)
	if err != nil {
		return "", err
	}
	if j.cryptoer != nil {
		return j.cryptoer.Encrypt(token)
	}
	return token, nil
}
```

2. 验证token

```go
// ParseToken parses the token.
func (j *CryptJWT) ParseToken(tokenString string) (*CustomClaims, error) {
	// 解密token
	if j.cryptoer != nil {
		var err error
		tokenString, err = j.cryptoer.Decrypt(tokenString)
		if err != nil {
			return nil, err
		}
	}
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return j.signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	claims, ok := token.Claims.(*CustomClaims)
	if ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, jwt.ErrInvalidType
}
```





### 缓存



对于查询次数较多的场景，如用户信息，视频信息等，本项目使用了缓存来减少数据库的查询次数，提高查询效率。



#### LRU

LRU(Least Recently Used) 最近最少使用缓存是一种常见的缓存策略，它会将最近最少使用的数据淘汰掉，从而保证缓存中的数据都是热点数据。



本项目在前期开发中使用了LRU缓存。LRU 缓存是使用双向链表和哈希表实现的。具体来说，它使用双向链表来维护缓存中的元素顺序，使用哈希表来实现快速查找元素。被查询的元素从链表中断开，移动到链表头节点，插入元素放到链表末尾，哈希表中存储key对应的链表节点，这样被查询多的元素总是留在链表的前面。



#### ARC

**在本项目中，我使用了 ARC 作为最终的缓存策略。**



传统LRU策略对最近访问的数据有很好的缓存效果。但有一些缺点，**LRU不能很好地处理突发请求**，当有大量新数据进入缓存时，LRU会将最近使用的数据替换掉，这可能会造成缓存污染导致缓存未命中率的急剧上升。



为了解决这些问题，ARC(Adaptive Replacement Cache) 缓存策略被提出。ARC 策略是基于 LRU 和 LFU(Least Frequently Used) 思想的组合，它会根据缓存中数据的访问情况动态地调整T1的大小，从而保证缓存中的数据都是热点数据。



ARC 策略的核心思想是将缓存分为两个部分：T1和T2,T1用来存放首次访问的数据，T2用来存放至少访问了两次的数据。此外ARC还保存了两条链的淘汰数据的key，B1和B2，保存着从T1和T2的历史淘汰信息，也称为ghost list。



本项目使用ARC缓存来减少数据库的查询次数，提高查询效率。相比于传统的LRU缓存，ARC缓存可以更好地适应不同的访问模式，从而提高缓存的命中率，进一步提高系统的性能。 






### 消息队列

本项目使用自己实现的简单高性能消息队列SimpleMQ对部分请求进行异步消峰，大大增强了项目的并发能力。



SimpleMQ的实现使用自己实现的可动态扩容的circularBuffer(底层为切片)[Doraemonkeys/arrayQueue](https://github.com/Doraemonkeys/arrayQueue)，相对于链表实现的队列，处理速度和空间利用率均有极大提升。



同时，消息队列的实现使用了泛型，使代码获得了类型检测，提高了代码复用能力，降低了心智负担和维护成本。

<img src="https://raw.githubusercontent.com/Doraemonkeys/picture/master/1/image-20230327191952425.png" alt="image-20230327191952425" style="zoom:80%;" />





[simpleMQ](https://github.com/Doraemonkeys/douyin2/blob/master/internal/pkg/messageQueue/simpleMQ_test.go) Benchmark测试的结果表明，**处理 10w 条并发数据的Push总共仅需 38ms**。

![image-20230327004502370](https://raw.githubusercontent.com/Doraemonkeys/picture/master/1/image-20230327004502370.png)





### 面向接口

本项目使用了面向接口的编程方式，将不同的功能模块分离到不同的层中，实现了解耦，使系统更易于维护和扩展。



在本项目中，定义了多个接口，如Cacher、Cryptoer、MQ、VideoStorageService等，通过这些接口，我们可以方便地实现缓存、加密、消息队列、视频对象存储等功能。同时，我们也可以通过实现这些接口来替换底层的具体实现，在后期替换新的技术栈如Redis，RabbitMQ时，可以做到无缝切换，从而实现更好的扩展性和灵活性。



- 缓存接口

```go
type Cacher[K comparable, T any] interface {
	// Get returns the value associated with the key.
	// Returns true if an eviction occurred.
	Get(key K) (T, bool)
	// Set sets the value associated with the key.
	// Returns true if the value was set.
	Set(key K, val T) bool
	// Delete deletes the value associated with the key.
	Delete(key K)
	// IsExist returns true if the key exists.
	IsExist(key K) bool
	// ClearAll clears all cache.
	ClearAll()
	// GetMulti returns the values associated with the keys.
	GetMulti(keys []K) map[K]T
	//PeekRandom returns a random value.
	PeekRandom() (T, error)
	// PeekRandomMulti returns random values.
	PeekRandomMulti(count int) ([]T, error)
	// SetMulti sets the values associated with the keys.
	SetMulti(kvs map[K]T) []bool
	// DeleteMulti deletes the values associated with the keys.
	DeleteMulti(keys []K)
	Len() int
	Cap() int
}
```



- 加解密接口

```go
type Cryptoer interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}
```



- 消息队列接口

```go
type MQ[T any] interface {
	// Push push a message to queue
	Push(T)
	// Len get the length of queue
	Len() int
}
```



- 对象存储接口

```go
type VideoStorageService[T any] interface {
	// Save object
	Save(T) (uint, error)
	// Delete object
	Delete(uint) error
	// Get object
	Get(uint) (T, error)
	// SaveUnique 保存视频，如果视频已经存在则返回已存在的视频ID和Error
	SaveUnique(T) (uint, error)
	GetURL(uint) (string, string, error)
}
```









### 日志

日志默认按日期分割，并将错误日志和普通日志的分离。当发生Panic等严重错误时，会单独创建文件对其保存。

![image-20230327010443049](https://raw.githubusercontent.com/Doraemonkeys/picture/master/1/image-20230327010443049.png)



本项目日志是基于logrus库的封装，实现了各种定制化配置。

```go
// 日志配置,可以为空
type LogConfig struct {
	//日志路径(可以为空)
	LogPath string
	//日志文件名后缀
	LogFileNameSuffix string
	//默认日志文件名(若按日期或大小分割日志，此项无效)
	DefaultLogName string
	//是否分离错误日志(Error级别以上)
	ErrSeparate bool
	//如果分离错误日志，普通日志文件是否仍然包含错误日志
	ErrInNormal bool
	//按日期分割日志(不能和按大小分割同时使用)
	DateSplit bool
	//取消日志输出到文件
	NoFile bool
	//取消日志输出到控制台
	NoConsole bool
	//取消时间戳Timestamp
	NoTimestamp bool
	//在控制台输出shortfile
	ShowShortFileInConsole bool
	//在控制台输出func
	ShowFuncInConsole bool
	//按大小分割日志,单位byte。(不能和按日期分割同时使用)
	MaxLogSize int64
	//日志扩展名(默认.log)
	LogExt string
	//panic,fatal,error,warn,info,debug,trace
	LogLevel string
	//时区
	TimeLocation *time.Location
	//在每条log末尾添加key-value
	key string
	//在每条log末尾添加key-value
	value interface{}
}
```







## 具体实现

本项目中**尽量避免了出现magic number和magic string**，增强了代码的可读性和可维护性。



### Init

<img src="https://raw.githubusercontent.com/Doraemonkeys/picture/master/1/run.png" alt="run" style="zoom:80%;" />



### 路由分发

```go
baseGroup := router.Group("/douyin")

// basic api
baseGroup.GET("/feed", middleware.JWTMiddleWare("/douyin/feed"), feed.FeedVideoListHandler)
baseGroup.POST("/user/register/", user.UserRegisterHandler)
baseGroup.POST("/user/login/", middleware.UserLoginHandler)
baseGroup.GET("/user/", middleware.JWTMiddleWare(), user.GetUserInfoHandler)
baseGroup.POST("/publish/action/", middleware.JWTMiddleWare(), publish.PublishVedioHandler)
baseGroup.GET("/publish/list/", middleware.JWTMiddleWare(), publish.QueryPublishListHandler)

//extend 1
baseGroup.POST("/favorite/action/", middleware.JWTMiddleWare(), favorite.PostFavorHandler)
baseGroup.GET("/favorite/list/", middleware.JWTMiddleWare(), favorite.QueryFavorVideoListHandler)
baseGroup.POST("/comment/action/", middleware.JWTMiddleWare(), comment.PostCommentHandler)
baseGroup.GET("/comment/list/", middleware.JWTMiddleWare(), comment.QueryCommentListHandler)

//extend 2
baseGroup.POST("/relation/action/", middleware.JWTMiddleWare(), follow.PostFollowActionHandler)
baseGroup.GET("/relation/follow/list/", middleware.JWTMiddleWare(), follow.QueryFollowListHandler)
baseGroup.GET("/relation/follower/list/", middleware.JWTMiddleWare(), follow.QueryFanListHandler)
```



### 视频流：/douyin/feed GET

<img src="https://raw.githubusercontent.com/Doraemonkeys/picture/master/1/feed2.png" alt="feed2" style="zoom: 67%;" />



### 用户信息与注册

![user](https://raw.githubusercontent.com/Doraemonkeys/picture/master/1/user.png)



### 登录

<img src="https://raw.githubusercontent.com/Doraemonkeys/picture/master/1/middleware.png" alt="middleware" style="zoom: 67%;" />



### 投稿与查询

<img src="https://raw.githubusercontent.com/Doraemonkeys/picture/master/1/publish.png" alt="publish" style="zoom: 67%;" />





### 喜欢列表

<img src="https://raw.githubusercontent.com/Doraemonkeys/picture/master/1/favorite.png" alt="favorite" style="zoom:67%;" />



### 评论列表

![comment](https://raw.githubusercontent.com/Doraemonkeys/picture/master/1/comment.png)

### 关注列表与粉丝列表

<img src="https://raw.githubusercontent.com/Doraemonkeys/picture/master/1/follow.png" alt="follow" style="zoom:67%;" />

### 点赞、关注、评论

<img src="https://raw.githubusercontent.com/Doraemonkeys/picture/master/1/msgQueue.png" alt="msgQueue" style="zoom:67%;" />



## 安全问题



### 用户密码的存储

用户密码使用bcrypt哈希函数取哈希值存入数据库，bcrypt是一种加盐的单向Hash加密算法，MD5加密时候，同一个密码经过hash的时候生成的是同一个hash值，在大数据的情况下，有些经过md5加密的方法将会被破解，而bcrypt能够很好的抵御彩虹表攻击。



### 重复注册

用户注册时会检查邮箱或username的唯一性，发现重复注册则返回错误。



### 权限检测

1. 除未登录用户获取视频流外，用户的所有操作均在JWT鉴权通过后处理。
2. 用户进行数据更改时，会检查数据的归属权是否为用户所有，杜绝了删除别人的评论，取消不存在的点赞等违规操作。
3. 评论的删除采用软删除的策略，以防意外情况发生。



### SQL 注入

1. 所有用户传入的参数均进行了合法性检查。
2. 避免使用SQL直接拼接，所有SQL语句均使用预处理语句进行预编译，彻底杜绝了SQL注入问题。



