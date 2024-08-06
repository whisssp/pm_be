package files

import (
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"math"
	"os"
	"pm/domain/entity"
	"pm/domain/repository/files"
	"pm/infrastructure/implementations/mailer"
	"pm/infrastructure/persistences/base"
)

type FileRepository struct {
	p *base.Persistence
}

func (fr FileRepository) ExportExcelProductReport() error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			zap.S().Errorw("error closing file ExportExcelProductReport: ", err.Error())
		}

	}()

	sheet1, err := f.NewSheet("Sheet1")
	if err != nil {
		zap.S().Errorw("error creating new sheet ExportExcelProductReport: ", err.Error())
		return err
	}

	headers := []string{"ID", "Name", "Created at", "Updated At", "Total order", "Total cost"}
	for i, h := range headers {
		cell := fmt.Sprintf("%s%v", string(rune(97+i)), 1)
		err := f.SetCellValue("Sheet1", cell, h)
		if err != nil {
			return err
		}
	}

	prods := make([]entity.Product, 0)
	if err := fr.p.GormDB.Model(&entity.Product{}).Find(&prods).Error; err != nil {
		zap.S().Errorw("error creating new sheet ExportExcelProductReport: ", err.Error())
		return err
	}
	productMap := make(map[int64]entity.Product)
	pIDs := make([]int64, 0)
	for _, p := range prods {
		productMap[int64(p.ID)] = p
		pIDs = append(pIDs, int64(p.ID))
	}
	//orders := make([]entity.Order, 0)
	//if err := fr.p.GormDB.Model(&entity.Order{}).Preload("OrderItems", "product_id IN (?)", pIDs).Find(&orders).Error; err != nil {
	//	zap.S().Errorw("error creating new sheet ExportExcelProductReport: ", err.Error())
	//	return err
	//}
	for index, id := range pIDs {
		//headers := []string{"ID", "Name", "Created at", "Updated At", "Total order", "Total cost"}
		cell := fmt.Sprintf("%s%v", "A", index+2)
		err := f.SetCellValue("Sheet1", cell, id)
		if err != nil {
			return err
		}
		cell = fmt.Sprintf("%s%v", "B", index+2)
		err = f.SetCellValue("Sheet1", cell, productMap[id].Name)
		if err != nil {
			return err
		}
		cell = fmt.Sprintf("%s%v", "C", index+2)
		err = f.SetCellValue("Sheet1", cell, productMap[id].CreatedAt)
		if err != nil {
			return err
		}
		cell = fmt.Sprintf("%s%v", "D", index+2)
		err = f.SetCellValue("Sheet1", cell, productMap[id].UpdatedAt)
		if err != nil {
			return err
		}

		cell = fmt.Sprintf("%s%v", "E", index+2)
		//var totalOrder int64 = 0
		orders := make([]entity.Order, 0)
		if err := fr.p.GormDB.Model(&entity.Order{}).Preload("OrderItems", "product_id = (?)", id).Find(&orders).Error; err != nil {
			zap.S().Errorw("error creating new sheet ExportExcelProductReport: ", err.Error())
			return err
		}
		err = f.SetCellValue("Sheet1", cell, len(orders))
		if err != nil {
			return err
		}

		var totalCost float64 = 0
		for _, o := range orders {
			for _, oi := range o.OrderItems {
				totalCost += float64(oi.Quantity) * productMap[int64(oi.ProductID)].Price
			}
		}
		err = f.SetCellValue("Sheet1", cell, math.Round(totalCost*100/100))
		if err != nil {
			return err
		}
	}
	f.SetActiveSheet(sheet1)

	if err := f.SaveAs("Product_Report.xlsx"); err != nil {
		return err
	}

	file, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	mailRepo := mailer.NewMailerRepository(fr.p)
	contentType, err := mimetype.DetectReader(file)
	if err != nil {
		return err
	}
	err = mailRepo.SendEmailWithAttachment("Product Report", "Demo Product Report", contentType.String(), []string{"nghia14802@gmail.com"}, file.Name())
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			zap.S().Errorw("error closing file ExportExcelProductReport: ", err.Error())
		}

	}()
	return nil
}

func NewFileRepository(p *base.Persistence) files.FileRepository {
	return &FileRepository{p: p}
}