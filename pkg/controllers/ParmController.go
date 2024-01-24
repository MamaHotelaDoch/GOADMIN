package controllers

import (
	//	"encoding/json"
	"errors"
	"fmt"
	"time"

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

type Gift struct {
	MSysID           int64     // BIGINT = 196491, /* — ID сообщения от администратора */
	MSvrNo           int16     // SMALLINT = 9991, /* — Номер вашего сервера */
	MItemID          int       // INT = 409, /* — Номер предмета (подарка) */
	MCnt             int       // INT = 1000, /* — Количество */
	MAvailablePeriod int       // INT = 0, /* — Доступный период (сколько будет лежать в подароках) */
	MPracticalPeriod int       // INT = 0, /* — Практический период (количество времени которое будет у предмета после получения)*/
	MBindingType     uint8     // TINYINT = 0, /* — Под замком предмет или нет (Нет = 0 | Да = 1) */
	MLimitedDate     time.Time // SMALLDATETIME = '2079-06-06', /* — Ограниченная дата */
	MItemStatus      uint8     // TINYINT = 1; /* — Статус предмета */

}

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

func PostGift(c *fiber.Ctx) error {
	nameItem: c.Query("NameItem")
	var itemID int
	// Получение времени из параметра запроса с помощью c.Query
	limitedDateStr := c.Query("LimitedDate")

	// Парсинг строки времени в значение time.Time
	limitedDate, err := time.Parse("2006-01-02T15:04:05", limitedDateStr)
	if err != nil {
		return err
	}

	gift := Gift {
		MSysID: 196491,
		MSvrNo: 9991,
		MItemID: itemID,
		MCnt: c.QueryInt("Cnt"),
		MAvailablePeriod: c.QueryInt("AvailablePeriod"),
		MPracticalPeriod: c.QueryInt("PracticalPeriod"),
		MBindingType: uint8(c.QueryInt("BindingType")),
		MLimitedDate: limitedDate,
		MItemStatus: 1,
	}

	ParmDB, err := config.ParmConfiguration()
	if err != nil { return err }

	searchItemName = `
	SELECT
	b.IID as 'itemID'
	FROM DT_Item as b
	WHERE b.IName LIKE ?
	`

	if err := ParmDB.Raw(searchItemName, ""+nameItem+"").Scan(&itemID).Error; err != nil {
		return err
	}

	
}

/*
func GetInfoBossDrop(c *fiber.Ctx) error {
	var bossDrop BossDrop
	if DB, err := config.ParmConfiguration(); err != nil {
		return err
	} else {
		DB.Debug().Preload("Group").Preload("Group.GroupTP").Preload("DDrops").First(&bossDrop.MonsterDrop)
		JsonData, err := json.MarshalIndent(bossDrop.MonsterDrop, "", "  ")
		if err != nil {
			return err
		}
		return c.JSON(string(JsonData))
	}
}
*/
