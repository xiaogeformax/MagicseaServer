package mysqldb

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type DBClient struct {
	ormDB orm.Ormer
}

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
}

//可以传多个models进去
func RegDBModule(models ...interface{}) {
	orm.RegisterModel(models...)
}

//Start 注册和连接db	dataSource连接信息  dbName数据库别名（不一定是数据库名）
func ConnectDB(dataSource string, dbName string) (*DBClient, error) {
	err := orm.RegisterDataBase(dbName, "mysql", dataSource)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	err2 := o.Using(dbName)
	if err2 != nil {
		return nil, err
	}
	return &DBClient{o}, nil
}

//Insert
func (client *DBClient) Insert(obj interface{}) (int64, error) {
	return client.ormDB.Insert(obj);
}

// Update
func (client *DBClient) Update(md interface{}, cols ...string) (int64, error) {
	return client.ormDB.Update(md, cols...)
}

//Insert Or Update
func (client *DBClient) InsertOrUpdate(md interface{}, colConflitAndArgs ...string) (int64, error) {
	return client.ormDB.InsertOrUpdate(md, colConflitAndArgs...)
}

//Delete 删除数据,condCols删除条件，默认Id字段
func (client *DBClient) Delete(md interface{}, condCols ...string) (int64, error) {
	return client.ormDB.Delete(md, condCols...)
}

//Read 有Cols用Cols做条件，没有，默认使用Id字段
// norow
func (client *DBClient) Read(md interface{}, cols ...string) (norow bool, e error) {
	e = client.ormDB.Read(md, cols...)
	norow = IsNoRow(e)
	return
}

//读取或者创建一行
func (client *DBClient) ReadOrCreate(md interface{}, col1 string, cols ...string) (bool, int64, error) {
	return client.ormDB.ReadOrCreate(md, col1, cols...)
}

//批量插入
func (client *DBClient) InsertMulti(bulk int, mds interface{}) (int64, error) {
	return client.ormDB.InsertMulti(bulk, mds)
}

//raw1
func (client *DBClient) Raw(query string, args ...interface{}) orm.RawSeter {
	return client.ormDB.Raw(query, args)
}

//报错解析
func IsNoRow(e error) bool {
	return e == orm.ErrNoRows
}

//事物相关
func (client *DBClient) Begin() error {
	return client.ormDB.Begin()
}

func (client *DBClient) Rollback() error {
	return client.ormDB.Rollback()
}

func (client *DBClient) Commit() error {
	return client.ormDB.Commit()
}
