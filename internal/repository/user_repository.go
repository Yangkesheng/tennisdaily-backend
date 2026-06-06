package repository

import (
	"errors"

	"tennisdaily-backend/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindOrCreateByOpenID(openid string) (*model.User, error) {
	var identity model.UserWechatIdentity
	appID := "legacy"
	err := r.db.Where("appid = ? AND openid = ? AND deleted_at IS NULL", appID, openid).First(&identity).Error
	if err == nil {
		var user model.User
		if err := r.db.Where("id = ? AND deleted_at IS NULL", identity.UserID).First(&user).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var user model.User
	err = r.db.Where("openid = ? AND deleted_at IS NULL", openid).First(&user).Error
	if err == nil {
		return &user, r.bindWechatIdentity(user.ID, appID, openid, "")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user = model.User{OpenID: openid}
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}
	if err := r.bindWechatIdentity(user.ID, appID, openid, ""); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id int64) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindOrCreateByPhoneAndWechat(phone string, appID string, openid string, unionid string) (*model.User, bool, error) {
	var isNewUser bool
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var identity model.UserWechatIdentity
		err := tx.Where("appid = ? AND openid = ? AND deleted_at IS NULL", appID, openid).First(&identity).Error
		if err == nil {
			var user model.User
			if err := tx.Where("id = ? AND deleted_at IS NULL", identity.UserID).First(&user).Error; err != nil {
				return err
			}
			updates := map[string]interface{}{"phone": phone}
			if user.OpenID == "" {
				updates["openid"] = openid
			}
			if err := tx.Model(&model.User{}).Where("id = ?", user.ID).Updates(updates).Error; err != nil {
				return err
			}
			if unionid != "" && identity.UnionID == "" {
				if err := tx.Model(&model.UserWechatIdentity{}).Where("id = ?", identity.ID).Update("unionid", unionid).Error; err != nil {
					return err
				}
			}
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		var user model.User
		err = tx.Where("phone = ? AND deleted_at IS NULL", phone).First(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = model.User{Phone: phone, OpenID: openid}
			if err := tx.Create(&user).Error; err != nil {
				return err
			}
			isNewUser = true
		} else if err != nil {
			return err
		} else if user.OpenID == "" {
			if err := tx.Model(&model.User{}).Where("id = ?", user.ID).Update("openid", openid).Error; err != nil {
				return err
			}
		}

		identity = model.UserWechatIdentity{UserID: user.ID, AppID: appID, OpenID: openid, UnionID: unionid}
		return tx.Create(&identity).Error
	})
	if err != nil {
		return nil, false, err
	}

	user, err := r.FindByPhone(phone)
	return user, isNewUser, err
}

func (r *UserRepository) FindByPhone(phone string) (*model.User, error) {
	var user model.User
	err := r.db.Where("phone = ? AND deleted_at IS NULL", phone).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateProfile(userID int64, nickname string, avatarURL string) (*model.User, error) {
	updates := map[string]interface{}{"nickname": nickname, "avatar_url": avatarURL}
	if err := r.db.Model(&model.User{}).Where("id = ? AND deleted_at IS NULL", userID).Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.FindByID(userID)
}

func (r *UserRepository) bindWechatIdentity(userID int64, appID string, openid string, unionid string) error {
	identity := model.UserWechatIdentity{UserID: userID, AppID: appID, OpenID: openid, UnionID: unionid}
	return r.db.Where("appid = ? AND openid = ?", appID, openid).FirstOrCreate(&identity).Error
}
