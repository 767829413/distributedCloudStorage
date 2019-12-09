package common

const (
	//静态文件地址
	StaticFileDir = "./static"
	//文件保存位置
	FileStoreTmp = "C:/Users/gogo/src/distributedCloudStorage/tmp/"

	//数据库文件状态
	FileStateAvailable = 1
	FileStateDisable   = 0
	FileStateDeleted   = 2

	//用户状态
	UserStateAvailable = 1
	UserStateDisable   = 0
	UserStateDeleted   = 2
	UserStatelocked    = 3

	//用户密码salt
	UserPwdSalt = "@#$%OKM"
	SecretKey   = "hello world"

	//用户过期时间,单位秒
	UserExpireTime = 3600
)
