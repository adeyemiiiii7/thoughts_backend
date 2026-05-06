package messages

import (
	"thoughts_backend_api/models"

	"gorm.io/gorm"
)

func canUsersMessage(db *gorm.DB, userAID, userBID uint) (bool, error) {
	var followCount int64

	if err := db.Model(&models.Follow{}).
		Where("follower_id = ? AND following_id = ?", userAID, userBID).
		Count(&followCount).Error; err != nil {
		return false, err
	}
	if followCount == 0 {
		return false, nil
	}

	if err := db.Model(&models.Follow{}).
		Where("follower_id = ? AND following_id = ?", userBID, userAID).
		Count(&followCount).Error; err != nil {
		return false, err
	}

	return followCount > 0, nil
}
