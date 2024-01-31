package controllers

import (
	"fmt"
	"strconv"
	"time"
	"github.com/Tibirlayn/GoAdmin/pkg/config"
	"github.com/Tibirlayn/GoAdmin/pkg/models/billing"
	"github.com/Tibirlayn/GoAdmin/pkg/models/game"
	"github.com/gofiber/fiber/v2"
)

func DeleteAllGift(c *fiber.Ctx) error {
	BillingDB, err := config.BillingConfiguration()
	if err != nil {
		return err
	}

	tx := BillingDB.Begin()
	if err := BillingDB.Unscoped().Where("1 = 1").Delete(&billing.SysOrderList{}).Error; err != nil { 
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "data delete",
	})
}

// выдать только один подарок на сервер
func PostGift(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	id, _ := strconv.ParseInt(data["id"], 10, 64)                        /* — ID сообщения от администратора */
	svr, _ := strconv.ParseInt(data["svr"], 10, 16)                      /* — Номер вашего сервера */
	itemid, _ := strconv.Atoi(data["itemid"])                            /* — Номер предмета (подарка) */
	cnt, _ := strconv.Atoi(data["cnt"])                                  /* — Количество */
	aperiod, _ := strconv.Atoi(data["aperiod"])                          /* — Доступный период (сколько будет лежать в подароках) */
	pperiod, _ := strconv.Atoi(data["pperiod"])                          /* — Практический период (количество времени которое будет у предмета после получения)*/
	binding, _ := strconv.ParseUint(data["binding"], 10, 8)              /* — Под замком предмет или нет (Нет = 0 | Да = 1) */
	limitedDate, _ := time.Parse("2006-01-02 15:04:05", data["limited"]) /* — Ограниченная дата */
	status, _ := strconv.ParseUint(data["status"], 10, 8)                /* — Статус предмета */

	BillingDB, err := config.BillingConfiguration()
	if err != nil {
		return err
	}
	GameDB, err := config.GameConfiguration()
	if err != nil {
		return err
	}

	var owners []int

	if err := GameDB.Model(&game.Pc{}).Distinct().Pluck("mOwner", &owners).Error; err != nil {
		return err
	}

	tx := BillingDB.Begin()
	for _, owner := range owners {
		giftPc := billing.SysOrderList{
			MSysID:           id,
			MUserNo:          owner,
			MSvrNo:           int16(svr),
			MItemID:          itemid,
			MCnt:             cnt,
			MAvailablePeriod: aperiod,
			MPracticalPeriod: pperiod,
			MBindingType:     uint8(binding),
			MLimitedDate:     limitedDate,
			MItemStatus:      uint8(status),
		}
		if err := BillingDB.Omit("mRegDate", "mReceiptDate", "mReceiptPcNo", "mRecepitPcNm").Create(&giftPc).Error; err != nil {
			tx.Rollback() // Откатить транзакцию при возникновении ошибки
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Gifts added",
	})

}

// выдать подарок всем персонажам
func PostGiftAll(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	id, _ := strconv.ParseInt(data["id"], 10, 64)                        /* — ID сообщения от администратора */
	svr, _ := strconv.ParseInt(data["svr"], 10, 16)                      /* — Номер вашего сервера */
	itemid, _ := strconv.Atoi(data["itemid"])                            /* — Номер предмета (подарка) */
	cnt, _ := strconv.Atoi(data["cnt"])                                  /* — Количество */
	aperiod, _ := strconv.Atoi(data["aperiod"])                          /* — Доступный период (сколько будет лежать в подароках) */
	pperiod, _ := strconv.Atoi(data["pperiod"])                          /* — Практический период (количество времени которое будет у предмета после получения)*/
	binding, _ := strconv.ParseUint(data["binding"], 10, 8)              /* — Под замком предмет или нет (Нет = 0 | Да = 1) */
	limitedDate, _ := time.Parse("2006-01-02 15:04:05", data["limited"]) /* — Ограниченная дата */
	status, _ := strconv.ParseUint(data["status"], 10, 8)                /* — Статус предмета */

	BillingDB, err := config.BillingConfiguration()
	if err != nil {
		return err
	}
	GameDB, err := config.GameConfiguration()
	if err != nil {
		return err
	}

	var owners []int

	if err := GameDB.Model(&game.Pc{}).Where("mDelDate IS NULL").Pluck("mOwner", &owners).Error; err != nil {
		fmt.Println("Error GameDB.Model(&game.Pc{})")
		return err
	}

	tx := BillingDB.Begin()
	for _, owner := range owners {
		giftPc := billing.SysOrderList{
			MSysID:           id,
			MUserNo:          owner,
			MSvrNo:           int16(svr),
			MItemID:          itemid,
			MCnt:             cnt,
			MAvailablePeriod: aperiod,
			MPracticalPeriod: pperiod,
			MBindingType:     uint8(binding),
			MLimitedDate:     limitedDate,
			MItemStatus:      uint8(status),
		}
		if err := BillingDB.Omit("mRegDate", "mReceiptDate", "mReceiptPcNo", "mRecepitPcNm").Create(&giftPc).Error; err != nil {
			tx.Rollback() // Откатить транзакцию при возникновении ошибки
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Gifts added",
	})
}

// выдать подарок персонажу
func PostGiftPcName(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	id, _ := strconv.ParseInt(data["id"], 10, 64)                        /* — ID сообщения от администратора */
	svr, _ := strconv.ParseInt(data["svr"], 10, 16)                      /* — Номер вашего сервера */
	itemid, _ := strconv.Atoi(data["itemid"])                            /* — Номер предмета (подарка) */
	mNm, _ := data["mNm"]                                                /* — Имя персонажа */
	cnt, _ := strconv.Atoi(data["cnt"])                                  /* — Количество */
	aperiod, _ := strconv.Atoi(data["aperiod"])                          /* — Доступный период (сколько будет лежать в подароках) */
	pperiod, _ := strconv.Atoi(data["pperiod"])                          /* — Практический период (количество времени которое будет у предмета после получения)*/
	binding, _ := strconv.ParseUint(data["binding"], 10, 8)              /* — Под замком предмет или нет (Нет = 0 | Да = 1) */
	limitedDate, _ := time.Parse("2006-01-02 15:04:05", data["limited"]) /* — Ограниченная дата */
	status, _ := strconv.ParseUint(data["status"], 10, 8)                /* — Статус предмета */

	BillingDB, err := config.BillingConfiguration()
	if err != nil {
		return err
	}
	GameDB, err := config.GameConfiguration()
	if err != nil {
		return err
	}

	var owners []int

	if err := GameDB.Model(&game.Pc{}).Where("mNm = ?", mNm).Pluck("mOwner", &owners).Error; err != nil {
		fmt.Println("Error GameDB.Model(&game.Pc{})")
		return err
	}

	tx := BillingDB.Begin()
	for _, owner := range owners {
		giftPc := billing.SysOrderList{
			MSysID:           id,
			MUserNo:          owner,
			MSvrNo:           int16(svr),
			MItemID:          itemid,
			MCnt:             cnt,
			MAvailablePeriod: aperiod,
			MPracticalPeriod: pperiod,
			MBindingType:     uint8(binding),
			MLimitedDate:     limitedDate,
			MItemStatus:      uint8(status),
		}
		if err := BillingDB.Omit("mRegDate", "mReceiptDate", "mReceiptPcNo", "mRecepitPcNm").Create(&giftPc).Error; err != nil {
			tx.Rollback() // Откатить транзакцию при возникновении ошибки
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Gifts added",
	})
}

//SQL Запрос. Добавить предмет в ШОП
func PostAddShopItem(c *fiber.Ctx) error {


	return c.JSON(fiber.Map{
		"status": "добавлен новый шоп предмет",
	})
}

// DECLARE @ItemID INT = 8531 /* ID предмета*/
// DECLARE @ICount INT = 1 /* Количество предметов */
// DECLARE @IName VARCHAR(40) = 'Особое Зельe Жизни' /* Название предмета */
// DECLARE @IDesc VARCHAR(500) = 'Средство, восстанавливающее большое количество здоровья.' /* Описание предмета */
// DECLARE @IPrice INT = 300 /* Цена */
// DECLARE @Istatus INT = 1 /* Статус 0 проклятый, 1-обычный, 2-благой */
// DECLARE @ICat INT = 3 /* Вкладка шопа 1,2,3,4 */
// DECLARE @IDay INT = 30 /* Время на предмете в днях*/
// DECLARE @IHour INT = 0 /* Время эффекта предмета в днях */
// DECLARE @SvrNo INT = 1164 /* Номер сервера */

// -- НИЖЕ НИЧЕГО НЕ ТРОГАТЬ --
// IF NOT EXISTS (SELECT 1 FROM FNLParm.dbo.DT_Item WHERE IID = @ItemID AND IIsCharge = 1)
// UPDATE FNLParm.dbo.DT_Item SET IIsCharge = 1 WHERE IID = @ItemID

// DECLARE @Date DATETIME SET @Date = GETDATE() /* Сегодняшняя дата */
// DECLARE @GIid INT = (SELECT MAX(GoldItemID) FROM TBLGoldItem) + 1
// DECLARE @GIOrder INT = (SELECT MAX(OrderNO) FROM TBLCategoryAssign) + 1

// IF @GIid IS NULL SET @GIid = 1
// IF @GIOrder IS NULL SET @GIOrder = 1
// INSERT INTO TBLGoldItem (GoldItemID, IID, ItemName, ItemDesc, OriginalGoldPrice, GoldPrice, ItemCategory, IsPackage, Status, AvailablePeriod, Count, PracticalPeriod, RegistAdmin)
// VALUES (@GIid, @ItemID, @IName, @IDesc, @IPrice, @IPrice, @ICat, 0, @Istatus, @IDay, @ICount, @IHour, 'R2Genius')
// INSERT INTO TBLCategoryAssign (CategoryID, GoldItemID, Status, OrderNO, RegistDate, RegistAdmin)
// VALUES (@ICat, @GIid, 1, @GIOrder, @Date, 'R2Genius')
// INSERT INTO TBLGoldItemSupportSvr(GoldItemID, mSvrNo)
// VALUES (@GIid, @SvrNo)