package authmodules

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/maximthomas/blazewall/auth-service/pkg/auth"
	"github.com/maximthomas/blazewall/auth-service/pkg/models"
	"reflect"
)

const (
	keyAdditionalFields = "additionalFields"
)

type Registration struct {
	BaseAuthModule
	afs []AdditionalFiled
}

type AdditionalFiled struct {
	dataStore string
	prompt    string
	required  bool
}

func (rm *Registration) Process(lss *auth.LoginSessionState, c *gin.Context) (ms auth.ModuleState, cbs []models.Callback, err error) {
	return auth.InProgress, rm.callbacks, err
}

func (rm *Registration) ProcessCallbacks(inCbs []models.Callback, lss *auth.LoginSessionState, c *gin.Context) (ms auth.ModuleState, cbs []models.Callback, err error) {
	if inCbs == nil {
		return auth.Fail, cbs, errors.New("Callbacks cant be nil")
	}

	callbacksValid := true
	errCbs := make([]models.Callback, len(rm.callbacks))
	copy(errCbs, rm.callbacks)

	for i, cb := range inCbs {
		if cb.Value == "" && errCbs[i].Required {
			(&errCbs[i]).Error = (&errCbs[i]).Prompt + " required"
			callbacksValid = false
		} else if i == 0 {
			_, exists := rm.r.UserRepo.GetUser(cb.Value)
			if exists {
				(&errCbs[i]).Error = "User exists"
				callbacksValid = false
			}
		}
	}

	if !callbacksValid {
		return auth.InProgress, errCbs, nil
	}

	var username string
	var password string

	fields := make(map[string]string, len(inCbs)-2)

	for _, cb := range inCbs {
		switch cb.Name {
		case "login":
			username = cb.Value
			break
		case "password":
			password = cb.Value
		default:
			fields[cb.Name] = cb.Value
		}
	}

	user := models.User{
		ID:         username,
		Properties: fields,
	}

	_, err = rm.r.UserRepo.CreateUser(user)
	if err != nil {
		return auth.Fail, cbs, err
	}

	err = rm.r.UserRepo.SetPassword(user.ID, password)
	if err != nil {
		return auth.Fail, cbs, err
	}

	return auth.Pass, rm.callbacks, err
}

func (rm *Registration) ValidateCallbacks(cbs []models.Callback) error {
	return rm.BaseAuthModule.ValidateCallbacks(cbs)
}

func NewRegistrationModule(base BaseAuthModule) *Registration {
	rm := &Registration{
		base,
		nil,
	}

	if af, ok := base.properties[keyAdditionalFields]; ok {
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
				Name:     af.dataStore,
				Type:     "text",
				Value:    "",
				Prompt:   af.prompt,
				Required: af.required,
			}
		}
	}
	adcbs[0] = models.Callback{
		Name:     "login",
		Type:     "text",
		Prompt:   "Login",
		Value:    "",
		Required: true,
	}
	adcbs[len(rm.afs)+1] = models.Callback{
		Name:     "password",
		Type:     "password",
		Prompt:   "Password",
		Value:    "",
		Required: true,
	}

	(&rm.BaseAuthModule).callbacks = adcbs
	return rm
}
