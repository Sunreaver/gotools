package database

// MgoConfiger 数据库配置
type MgoConfiger struct {
	DB                  string
	Colection           string
	SocketTimeoutSecond uint
}

// GetDBName 获取数据库名字
func (c *MgoConfiger) GetDBName() string {
	if c == nil {
		return ""
	}
	return c.DB
}

// GetCollectionName 获取数据库表名字
func (c *MgoConfiger) GetCollectionName() string {
	if c == nil {
		return ""
	}
	return c.Colection
}

// GetSocketTimeoutSecond 获取数据库连接超时时间
func (c *MgoConfiger) GetSocketTimeoutSecond() uint {
	if c == nil || c.SocketTimeoutSecond == 0 {
		// 默认1分钟
		return 60
	}
	return c.SocketTimeoutSecond
}
