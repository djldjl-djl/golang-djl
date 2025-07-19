package sql

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"k8s.io/klog/v2"
)

const (
	username = "djl"
	password = "123456"
	ip       = "192.168.85.90"
	dbname   = "golang"
)

type User struct {
	gorm.Model        // 包含 ID、CreatedAt、UpdatedAt、DeletedAt
	Username   string `gorm:"type:varchar(100);not null;unique"`
	Password   string `gorm:"type:varchar(255)"`
	Email      string `gorm:"type:varchar(100)"`
}

type DBdjl struct {
	Db *gorm.DB
}

func Lianjie() *DBdjl {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, ip, dbname)
	d := DBdjl{}
	var err error
	d.Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		klog.Errorf("数据库链接失败:%v", err)
	}
	return &d
}
func (d *DBdjl) Linajiesql() error {
	err := d.Db.AutoMigrate(&User{})
	if err != nil {
		klog.Errorf("数据表创建失败: %v", err)
		return err
	}

	fmt.Println("✅ 数据表创建成功")
	return nil
}
func (d *DBdjl) Insertuser(name, pass, email string) error {
	user := User{Username: name, Password: pass, Email: email}
	result := d.Db.Create(&user) // 通过数据的指针来创建
	if result.Error != nil {
		return fmt.Errorf("插入失败:%v", result.Error)
	}
	fmt.Printf("插入成功，用户ID:%v\n", user.ID)
	return nil
}
