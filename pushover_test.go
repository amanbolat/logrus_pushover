package logrusPushover

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"time"
)

// get pushoverUserToken, pushoverAPIToken from ENV
func getTokensFromEnv() (pushoverUserToken, pushoverAPIToken string, err error) {
	pushoverUserToken = os.Getenv("PUSHOVER_USER_TOKEN")
	pushoverAPIToken = os.Getenv("PUSHOVER_API_TOKEN")
	if pushoverUserToken == "" || pushoverAPIToken == "" {
		err = errors.New("set env var PUSHOVER_API_TOKEN and PUSHOVER_USER_TOKEN")
	}
	return
}

func TestSync(t *testing.T) {
	pushoverUserToken, pushoverAPIToken, err := getTokensFromEnv()
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}

	hook := NewPushoverAsyncHook(pushoverAPIToken, pushoverUserToken)

	msgOne := "message one"
	//msgTwo := "message two"
	//msgThree := "message three"

	log := logrus.New()
	log.Out = ioutil.Discard
	log.Hooks.Add(hook)
	log.WithFields(logrus.Fields{"withField": "1", "filterMe": "1"}).Error(msgOne)
	//log.WithFields(logrus.Fields{"withField": "1", "filterMe": "1"}).Error(msgTwo)
	//log.WithFields(logrus.Fields{"withField": "1", "filterMe": "1"}).Error(msgThree)
	time.Sleep(time.Second * 5)

	t.Log()
}

//func TestAsync(t *testing.T) {
//	userToken, appToken, err := getTokensFromEnv()
//	if err != nil {
//		println(err.Error())
//		os.Exit(0)
//	}
//
//	hook := NewPushoverAsyncHook(appToken, userToken, time.Second*3)
//
//	msgOne := "message one"
//	msgTwo := "message two"
//	msgThree := "message three"
//
//	log := logrus.New()
//	log.Out = ioutil.Discard
//	log.Hooks.Add(hook)
//	log.WithFields(logrus.Fields{"withField": "1", "filterMe": "1"}).Error(msgOne)
//	log.WithFields(logrus.Fields{"withField": "1", "filterMe": "1"}).Error(msgTwo)
//	log.WithFields(logrus.Fields{"withField": "1", "filterMe": "1"}).Error(msgThree)
//	time.Sleep(time.Second * 10)
//
//}

//func TestSetDuration(t *testing.T) {
//	hook, err := getNewHook()
//	if err != nil {
//		t.Error("expected err == nil, got", err)
//	}
//	err = hook.SetMuteDelay("blabla")
//	if err == nil {
//		t.Error("expected err != nil, got", err)
//	}
//	err = hook.SetMuteDelay("15m")
//	if err != nil {
//		t.Error("expected err == nil, got", err)
//	}
//}
