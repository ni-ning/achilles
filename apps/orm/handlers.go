package orm

import (
	"achilles/global"
	"achilles/internal/model"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func createHandler(c *gin.Context) {
	var age2 uint8 = 0
	user1 := model.User{Username: "linda", Age: 0, Age2: &age2}
	// 实际创建时 age=18，age2=0 写入到数据库，是个坑
	// INSERT INTO `users` (`username`,`age`,`age2`,`created_at`,`updated_at`,`is_deleted`) VALUES ('linda',18,0,'2022-12-09 17:48:22.147','2022-12-09 17:48:22.147',false)
	result1 := global.DBEngine.Debug().Create(&user1)          // 通过数据的指针来创建
	fmt.Println(user1.ID)                                      // 返回插入数据的主键
	fmt.Println("result1.Error:", result1.Error)               // 返回 error
	fmt.Println("result1.RowsAffected:", result1.RowsAffected) // 返回插入记录的条数

	user2 := model.User{Username: "tom", Age: 18}
	// INSERT INTO `users` (`username`,`age`,`created_at`,`updated_at`) VALUES ('tom',18,'2022-12-09 17:25:03.652','2022-12-09 17:25:03.652')
	result2 := global.DBEngine.Debug().Select("Username", "Age").Create(&user2) // 用指定的字段创建记录
	fmt.Println("result2.Error:", result2.Error)                                // 返回 error
	fmt.Println("result2.RowsAffected:", result2.RowsAffected)                  // 返回插入记录的条数

	// 创建时不包含 username 和 age 字段
	// models.AchillesDB.Omit("Username", "Age").Create(&user)

	// 批量创建
	var users = []model.User{
		{Username: "linda001", Age: 21},
		{Username: "linda002", Age: 22},
		{Username: "linda003", Age: 23},
	}
	// 指定每批的数量
	result3 := global.DBEngine.CreateInBatches(users, 100)
	fmt.Println("result3.Error:", result3.Error)               // 返回 error
	fmt.Println("result3.RowsAffected:", result3.RowsAffected) // 返回插入记录的条数

	// 根据 Map 创建，用Model(&models.User{})来识别具体的库表
	result4 := global.DBEngine.Model(&model.User{}).Create(map[string]interface{}{
		"Username": "catherine", "Age": 18,
	})
	fmt.Println("result4.Error:", result4.Error)               // 返回 error
	fmt.Println("result4.RowsAffected:", result4.RowsAffected) // 返回插入记录的条数

	c.JSON(http.StatusOK, map[string]interface{}{
		"user1": user1, // 程序级别定义的模型示例，和数据库返回值没有联系，创建时是吧，该user1仅仅是穿件的原始数据
		"user2": user2,
	})
}

func selectHandler(c *gin.Context) {
	// 定义变量去接收数据，无数据时零值
	var user1 model.User
	var user2 model.User
	var user3 model.User
	var user4 model.User
	global.DBEngine.Model(&model.User{}).First(&user1)
	global.DBEngine.Model(&model.User{}).Last(&user2)
	global.DBEngine.Model(&model.User{}).Take(&user3)
	global.DBEngine.Model(&model.User{}).Find(&user4)

	result := map[string]interface{}{}
	ret := global.DBEngine.Model(&model.User{}).Take(&result)
	fmt.Println(ret.RowsAffected) // 返回找到的记录数
	fmt.Println(ret.Error)        // returns error or nil
	ok := errors.Is(ret.Error, gorm.ErrRecordNotFound)
	fmt.Println(ok)

	c.JSON(http.StatusOK, map[string]interface{}{
		"result":      result,
		"First:user1": user1,
		"Last:user2":  user2,
		"Take:user3":  user3,
		"Find:user4":  user4,
	})
}
