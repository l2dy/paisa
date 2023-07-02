package price

import (
	"time"

	"gorm.io/gorm"

	"github.com/google/btree"
	log "github.com/sirupsen/logrus"
)

type CommodityType string

const (
	MutualFund CommodityType = "mutualfund"
	NPS        CommodityType = "nps"
	Stock      CommodityType = "stock"
	Unknown    CommodityType = "unknown"
)

type Price struct {
	ID            uint          `gorm:"primaryKey" json:"id"`
	Date          time.Time     `json:"date"`
	CommodityType CommodityType `json:"commodity_type"`
	CommodityID   string        `json:"commodity_id"`
	CommodityName string        `json:"commodity_name"`
	Value         float64       `json:"value"`
}

func (p Price) Less(o btree.Item) bool {
	return p.Date.Before(o.(Price).Date)
}

func UpsertAllByTypeAndID(db *gorm.DB, commodityType CommodityType, commodityID string, prices []*Price) {
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(&Price{}, "commodity_type = ? and commodity_id = ?", commodityType, commodityID).Error
		if err != nil {
			return err
		}
		for _, price := range prices {
			err := tx.Create(price).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func UpsertAllByType(db *gorm.DB, commodityType CommodityType, prices []Price) {
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(&Price{}, "commodity_type = ?", commodityType).Error
		if err != nil {
			return err
		}
		for _, price := range prices {
			err := tx.Create(&price).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
