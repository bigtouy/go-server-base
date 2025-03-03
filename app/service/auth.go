package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go-server-base/app/dto"
	"go-server-base/constant"
	"go-server-base/global"
	"go-server-base/utils/encrypt"
	"go-server-base/utils/jwt"
	"go-server-base/utils/mfa"
)

type AuthService struct{}

type IAuthService interface {
	SafetyStatus(c *gin.Context) error
	CheckIsFirst() bool
	InitUser(c *gin.Context, req dto.InitUser) error
	VerifyCode(code string) (bool, error)
	SafeEntrance(c *gin.Context, code string) error
	Login(c *gin.Context, info dto.Login) (*dto.UserLoginInfo, error)
	LogOut(c *gin.Context) error
	MFALogin(c *gin.Context, info dto.MFALogin) (*dto.UserLoginInfo, error)
}

func NewIAuthService() IAuthService {
	return &AuthService{}
}

func (u *AuthService) SafeEntrance(c *gin.Context, code string) error {
	codeWithMD5 := encrypt.Md5(code)
	cookieValue, _ := encrypt.StringEncrypt(codeWithMD5)
	c.SetCookie(codeWithMD5, cookieValue, 604800, "", "", false, false)

	expiredSetting, err := settingRepo.Get(settingRepo.WithByKey("ExpirationDays"))
	if err != nil {
		return err
	}
	timeout, _ := strconv.Atoi(expiredSetting.Value)
	if err := settingRepo.Update("ExpirationTime", time.Now().AddDate(0, 0, timeout).Format("2006-01-02 15:04:05")); err != nil {
		return err
	}
	return nil
}

func (u *AuthService) Login(c *gin.Context, info dto.Login) (*dto.UserLoginInfo, error) {
	nameSetting, err := settingRepo.Get(settingRepo.WithByKey("UserName"))
	if err != nil {
		return nil, errors.WithMessage(constant.ErrRecordNotFound, err.Error())
	}
	passwrodSetting, err := settingRepo.Get(settingRepo.WithByKey("Password"))
	if err != nil {
		return nil, errors.WithMessage(constant.ErrRecordNotFound, err.Error())
	}
	pass, err := encrypt.StringDecrypt(passwrodSetting.Value)
	if err != nil {
		return nil, constant.ErrAuth
	}
	if info.Password != pass || nameSetting.Value != info.Name {
		return nil, constant.ErrAuth
	}
	mfa, err := settingRepo.Get(settingRepo.WithByKey("MFAStatus"))
	if err != nil {
		return nil, err
	}
	if mfa.Value == "enable" {
		return &dto.UserLoginInfo{Name: nameSetting.Value, MfaStatus: mfa.Value}, nil
	}

	return u.generateSession(c, info.Name, info.AuthMethod)
}

func (u *AuthService) MFALogin(c *gin.Context, info dto.MFALogin) (*dto.UserLoginInfo, error) {
	nameSetting, err := settingRepo.Get(settingRepo.WithByKey("UserName"))
	if err != nil {
		return nil, errors.WithMessage(constant.ErrRecordNotFound, err.Error())
	}
	passwrodSetting, err := settingRepo.Get(settingRepo.WithByKey("Password"))
	if err != nil {
		return nil, errors.WithMessage(constant.ErrRecordNotFound, err.Error())
	}
	pass, err := encrypt.StringDecrypt(passwrodSetting.Value)
	if err != nil {
		return nil, err
	}
	if info.Password != pass || nameSetting.Value != info.Name {
		return nil, constant.ErrAuth
	}

	mfaSecret, err := settingRepo.Get(settingRepo.WithByKey("MFASecret"))
	if err != nil {
		return nil, err
	}
	success := mfa.ValidCode(info.Code, mfaSecret.Value)
	if !success {
		return nil, constant.ErrAuth
	}

	return u.generateSession(c, info.Name, info.AuthMethod)
}

func (u *AuthService) generateSession(c *gin.Context, name, authMethod string) (*dto.UserLoginInfo, error) {
	setting, err := settingRepo.Get(settingRepo.WithByKey("SessionTimeout"))
	if err != nil {
		return nil, err
	}
	lifeTime, err := strconv.Atoi(setting.Value)
	if err != nil {
		return nil, err
	}

	if authMethod == constant.AuthMethodJWT {
		j := jwt.NewJWT()
		claims := j.CreateClaims(jwt.BaseClaims{
			Name: name,
		})
		token, err := j.CreateToken(claims)
		if err != nil {
			return nil, err
		}
		return &dto.UserLoginInfo{Name: name, Token: token}, nil
	}
	sID, _ := c.Cookie(constant.SessionName)
	sessionUser, err := global.SESSION.Get(sID)
	if err != nil {
		sID = uuid.New().String()
		c.SetCookie(constant.SessionName, sID, 604800, "", "", false, false)
		err := global.SESSION.Set(sID, sessionUser, lifeTime)
		if err != nil {
			return nil, err
		}
		return &dto.UserLoginInfo{Name: name}, nil
	}
	if err := global.SESSION.Set(sID, sessionUser, lifeTime); err != nil {
		return nil, err
	}

	return &dto.UserLoginInfo{Name: name}, nil
}

func (u *AuthService) LogOut(c *gin.Context) error {
	sID, _ := c.Cookie(constant.SessionName)
	if sID != "" {
		c.SetCookie(constant.SessionName, sID, -1, "", "", false, false)
		err := global.SESSION.Delete(sID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *AuthService) VerifyCode(code string) (bool, error) {
	setting, err := settingRepo.Get(settingRepo.WithByKey("SecurityEntrance"))
	if err != nil {
		return false, err
	}
	return setting.Value == code, nil
}

func (u *AuthService) SafetyStatus(c *gin.Context) error {
	setting, err := settingRepo.Get(settingRepo.WithByKey("SecurityEntrance"))
	if err != nil {
		return err
	}
	codeWithEcrypt, err := c.Cookie(encrypt.Md5(setting.Value))
	if err != nil {
		return err
	}
	code, err := encrypt.StringDecrypt(codeWithEcrypt)
	if err != nil {
		return err
	}
	if code != encrypt.Md5(setting.Value) {
		return errors.New("code not match")
	}
	return nil
}

func (u *AuthService) CheckIsFirst() bool {
	user, _ := settingRepo.Get(settingRepo.WithByKey("UserName"))
	pass, _ := settingRepo.Get(settingRepo.WithByKey("Password"))
	return len(user.Value) == 0 || len(pass.Value) == 0
}

func (u *AuthService) InitUser(c *gin.Context, req dto.InitUser) error {
	user, _ := settingRepo.Get(settingRepo.WithByKey("UserName"))
	pass, _ := settingRepo.Get(settingRepo.WithByKey("Password"))
	if len(user.Value) == 0 || len(pass.Value) == 0 {
		newPass, err := encrypt.StringEncrypt(req.Password)
		if err != nil {
			return err
		}
		if err := settingRepo.Update("UserName", req.Name); err != nil {
			return err
		}
		if err := settingRepo.Update("Password", newPass); err != nil {
			return err
		}
		expiredSetting, err := settingRepo.Get(settingRepo.WithByKey("ExpirationDays"))
		if err != nil {
			return err
		}
		timeout, _ := strconv.Atoi(expiredSetting.Value)
		if timeout != 0 {
			if err := settingRepo.Update("ExpirationTime", time.Now().AddDate(0, 0, timeout).Format("2006-01-02 15:04:05")); err != nil {
				return err
			}
		}
		return nil
	}

	return fmt.Errorf("can't init user because user %s is in system", user.Value)
}
