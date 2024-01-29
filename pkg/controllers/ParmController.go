package controllers

import (
	//	"encoding/json"
	"errors"
	"fmt"

	"github.com/Tibirlayn/GoAdmin/pkg/config"
	"github.com/gofiber/fiber/v2"
)

type ItemResult struct {
	Boss      string
	NameGroup string
	NameItem  string
	ItemDesc  string
	Event     string
}

type SpecificItem struct {
	ItemIID  int
	ItemName string
	X        int64
	Y        int64
	Z        int64
}

//type Gift struct {
//	MSysID           int64     // BIGINT = 196491, /* — ID сообщения от администратора */
//	MSvrNo           int16     // SMALLINT = 9991, /* — Номер вашего сервера */
//	MItemID          int       // INT = 409, /* — Номер предмета (подарка) */
//	MCnt             int       // INT = 1000, /* — Количество */
//	MAvailablePeriod int       // INT = 0, /* — Доступный период (сколько будет лежать в подароках) */
//	MPracticalPeriod int       // INT = 0, /* — Практический период (количество времени которое будет у предмета после получения)*/
//	MBindingType     uint8     // TINYINT = 0, /* — Под замком предмет или нет (Нет = 0 | Да = 1) */
//	MLimitedDate     time.Time // SMALLDATETIME = '2079-06-06', /* — Ограниченная дата */
//	MItemStatus      uint8     // TINYINT = 1; /* — Статус предмета */
//}

func GetInfoBossDrop(c *fiber.Ctx) error {
	mid := c.Query("MID")
	name := c.Query("MName")
	if mid == "" && name == "" {
		return errors.New("Введите id монстра или название монстра! Поле не может быть пустым!")
	} else if mid != "" && name != "" {
		return errors.New("Введите что-то одно: id монстра или название монстра!")
	}
	if DB, err := config.ParmConfiguration(); err != nil {
		return err
	} else {
		results := []ItemResult{}

		query := `
			SELECT DISTINCT
				mo.MName AS 'Boss',
				tdg.DName AS 'NameGroup',
				it.IName AS 'NameItem',
				tpis.mDesc AS 'ItemDesc',
				CASE
					WHEN evo.mObjID IS NULL THEN 'Not event item'
					ELSE 'Event item'
				END AS 'Event'
				FROM DT_Monster AS mo
				INNER JOIN DT_MonsterDrop AS md ON mo.MID = md.MID
				INNER JOIN DT_DropGroup AS dg ON md.DGroup = dg.DGroup
				INNER JOIN TP_DropGroup AS tdg ON md.DGroup = tdg.DGroup
				INNER JOIN DT_DropItem AS di ON dg.DDrop = di.DDrop
				INNER JOIN TP_ItemStatus AS tpis ON di.DStatus = tpis.mStatus
				INNER JOIN DT_Item AS it ON di.DItem = it.IID
				LEFT JOIN TblEventObj AS evo ON it.IID = evo.mObjID
				WHERE
					mo.MID = ? OR mo.MName = ? AND
					di.DIsEvent = 0
		`
		if err := DB.Raw(query, mid, name).Scan(&results).Error; err != nil {
			// Handle error
			return err
		}

		// проверка для себя
		if false {
			for _, item := range results {
				fmt.Println(item)
			}
		}
		return c.JSON(results)

	}
}

func GetSpecificProcItem(c *fiber.Ctx) error {
	name := c.Query("NameItem")
	fmt.Println(name)

	if DB, err := config.ParmConfiguration(); err != nil {
		return err
	} else {
		result := []SpecificItem{}

		query := `
		SELECT
		a.mIID as 'ItemIID',
		b.IName as 'ItemName',
		a.mAParam as 'X',
		a.mBParam as 'Y',
		a.mCParam as 'Z'
		FROM TblSpecificProcItem as a
		INNER JOIN DT_Item as b on (a.mIID = b.IID)
		WHERE b.IName LIKE ?
		`

		if err := DB.Raw(query, "%"+name+"%").Scan(&result).Error; err != nil {
			return err
		}

		return c.JSON(result)
	}
}

// Просмотр DT_Refine > сделать поиск по названию
func GetRefine(c *fiber.Ctx, pageNumber int, limitCnt int) error {

	ParmDB, err := config.ParmConfiguration()
	if err != nil {
		return err
	}

	// Рассчитываем смещение (offset) и лимит записей на основе номера страницы
	limit := limitCnt
	offset := (pageNumber - 1) * limit

	var results []struct {
		ID               int     `gorm:"column:RID"`
		ReceivedItem     int     `gorm:"column:RItemID0"`
		ReceivedItemName string  `gorm:"column:IName"`
		RecipeItemID     int     `gorm:"column:RItemID"`
		RecipeItemName   string  `gorm:"column:IName"`
		SuccessChance    float64 `gorm:"column:RSuccess"`
	}

	if err := ParmDB.Table("DT_Refine a").
		Select("a.RID as ID, a.RItemID0, b.IName as ReceivedItemName, c.RItemID as RecipeItemID, b1.IName as RecipeItemName, a.RSuccess as SuccessChance").
		Joins("INNER JOIN DT_Item as b ON a.RItemID0 = b.IID").
		Joins("INNER JOIN DT_RefineMaterial as c ON a.RID = c.RID").
		Joins("INNER JOIN DT_Item as b1 ON c.RItemID = b1.IID").
		Offset(offset).
		Limit(limit).
		Scan(&results).Error; err != nil {
		return err
	}

	return c.JSON(results)
}

func GetRefineByName(c *fiber.Ctx, pageNumber int, limitCnt int) error {
	ParmBD, err := config.ParmConfiguration()
	if err != nil {
		return err
	}

	limit := limitCnt
	offset := (pageNumber -1) * limit

	var results []struct {
		ID               int     `gorm:"column:RID"`
		ReceivedItem     int     `gorm:"column:RItemID0"`
		ReceivedItemName string  `gorm:"column:IName"`
		RecipeItemName   string  `gorm:"column:IName"`
		SuccessChance    float64 `gorm:"column:RSuccess"`
	}

	if err := ParmBD.Table("DT_Refine a").
	Select("a.RID, a.RItemID0 as ReceivedItem, b.IName as ReceivedItemName, STRING_AGG(b1.IName, ', ') as RecipeItemName, a.RSuccess as SuccessChance").
	Joins("INNER JOIN DT_Item b ON a.RItemID0 = b.IID").
	Joins("INNER JOIN DT_RefineMaterial c ON a.RID = c.RID").
	Joins("INNER JOIN DT_Item b1 ON c.RItemID = b1.IID").
	Group("a.RID, a.RItemID0, b.IName, a.RSuccess").
	Order("RID").
	Offset(offset).
	Limit(limit).
	Scan(&results).Error; err != nil {
	return err
	}

	return c.JSON(results)
}

/* 
SELECT
a.RID, 
a.RItemID0 as 'Получаемый предмет',
b.IName as 'Название получаемого предмета',
STRING_AGG(b1.IName, ', ') as 'Названия предметов из рецепта',
a.RSuccess as 'Шанс проточки'
FROM DT_Refine as a
INNER JOIN DT_Item as b on (a.RItemID0 = b.IID)
INNER JOIN DT_RefineMaterial as c on (a.RID = c.RID)
INNER JOIN DT_Item as b1 on (c.RItemID = b1.IID)
GROUP BY a.RID, a.RItemID0, b.IName, a.RSuccess
ORDER BY [RID];
 */

/*  
SELECT
a.RID,
a.RItemID0 as 'Получаемый предмет',
b.IName as 'Навание получаемого предмета',
c.RItemID as 'Рецепт',
b1.IName as 'Название предмета из рецепта',
a.RSuccess as 'Шанс проточки'
FROM DT_Refine as a
INNER JOIN DT_Item as b on (a.RItemID0 = b.IID)
INNER JOIN DT_RefineMaterial as c on (a.RID = c.RID)
INNER JOIN DT_Item as b1 on (c.RItemID = b1.IID)
 */