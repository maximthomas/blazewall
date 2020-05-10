package authmodules

import (
	"github.com/gin-gonic/gin"
	"github.com/maximthomas/blazewall/auth-service/pkg/auth"
	"github.com/maximthomas/blazewall/auth-service/pkg/models"
	"reflect"
)

const (
	keyAdditionalFileds = "additionalFileds"
)

type Registration struct {
	BaseAuthModule
	afs []AdditionalFiled
}

type AdditionalFiled struct {
	dataStore string
	prompt    string
}

func (rm *Registration) Process(lss *auth.LoginSessionState, c *gin.Context) (ms auth.ModuleState, cbs []models.Callback, err error) {
	return auth.InProgress, rm.callbacks, err
}

func (rm *Registration) ProcessCallbacks(inCbs []models.Callback, lss *auth.LoginSessionState, c *gin.Context) (ms auth.ModuleState, cbs []models.Callback, err error) {
	var username string
	var password string

	for _, cb := range inCbs {
		switch cb.Name {
		case "login":
			username = cb.Value
			break
		case "password":
			password = cb.Value
		}
	}

	valid := rm.r.UserRepo.ValidatePassword(username, password)
	if valid {
		lss.UserId = username
		return auth.Pass, cbs, err
	} else {
		cbs = rm.callbacks
		(&cbs[0]).Error = "Invalid username or password"
		return auth.InProgress, cbs, err
	}
}

func (rm *Registration) ValidateCallbacks(cbs []models.Callback) error {
	return rm.BaseAuthModule.ValidateCallbacks(cbs)
}

func NewRegistrationModule(base BaseAuthModule) *Registration {
	rm := &Registration{
		base,
		nil,
	}

	if af, ok := base.properties[keyAdditionalFileds]; ok {
		afObj := reflect.ValueOf(af)
		afs := make([]AdditionalFiled, afObj.Len())
		for i := 0; i < afObj.Len(); i++ {
			adf := afObj.Index(i).Interface().(AdditionalFiled)
			afs[i] = adf
		}
		rm.afs = afs
	}
	adcbs := make([]models.Callback, len(rm.afs)+2)
	if rm.afs != nil {
		for i, af := range rm.afs {
			adcbs[i+1] = models.Callback{
				Name:   af.dataStore,
				Type:   "text",
				Value:  "",
				Prompt: af.prompt,
			}
		}
	}
	adcbs[0] = models.Callback{
		Name:   "login",
		Type:   "text",
		Prompt: "Login",
		Value:  "",
	}
	adcbs[len(rm.afs)+1] = models.Callback{
		Name:   "password",
		Type:   "password",
		Prompt: "Password",
		Value:  "",
	}

	(&rm.BaseAuthModule).callbacks = adcbs
	return rm
}
