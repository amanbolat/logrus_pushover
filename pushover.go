package logrusPushover

import (
	"time"

	"fmt"
	"github.com/gregdel/pushover"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strings"
)

const (
	defaultDelay = 5
)

// PushoverHook sends log via Pushover (https://pushover.net/)
type PushoverHook struct {
	async  bool
	worker *worker
}

type worker struct {
	app       *pushover.Pushover
	recipient *pushover.Recipient
	muteDelay time.Duration
	messages  chan message
}

func (w *worker) run() {
	go func() {
		for {
			select {
			case m := <-w.messages:
				w.send(m)
				time.Sleep(w.muteDelay)
			}
		}
	}()
}

type message struct {
	title string
	body  string
}

// NewPushoverHook init & returns a new PushoverHook
func NewPushoverHook(appToken string, userToken string) *PushoverHook {
	return newPushoverHook(appToken, userToken, false)
}

// NewPushoverAsyncHook init & returns a new async PushoverHook
func NewPushoverAsyncHook(appToken string, userToken string) *PushoverHook {
	return newPushoverHook(appToken, userToken, true)
}

// newPushoverHook init & returns a new PushoverHook
func newPushoverHook(appToken string, userToken string, async bool) *PushoverHook {

	w := &worker{
		muteDelay: defaultDelay,
		messages:  make(chan message, 100),
		app:       pushover.New(appToken),
		recipient: pushover.NewRecipient(userToken),
	}

	p := PushoverHook{
		async:  async,
		worker: w,
	}

	if p.async {
		p.worker.run()
	}

	return &p
}

// Levels returns the available logging levels.
func (hook *PushoverHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	}
}

// Fire is called when a log event is fired.
func (hook *PushoverHook) Fire(entry *logrus.Entry) error {
	var messageBody string
	for k, v := range entry.Data {
		messageBody = messageBody + fmt.Sprintf("%s: %v\n", k, v)
	}

	messageBody = messageBody + fmt.Sprintf("<font color=\"#a0a0a0\">%s</font>", entry.Time.Format("02.01.2006 15:04:05"))

	title := fmt.Sprintf("%s: %s", strings.ToUpper(entry.Level.String()), entry.Message)

	if len(title) > pushover.MessageTitleMaxLength {
		title = title[0:pushover.MessageTitleMaxLength]
	}

	if len(messageBody) > pushover.MessageMaxLength {
		messageBody = messageBody[0:pushover.MessageMaxLength]
	}

	m := message{title, messageBody}

	if hook.async {
		hook.worker.messages <- m
		return nil
	}

	err := hook.worker.send(m)
	if err != nil {
		return err
	}

	return nil
}

func (w *worker) send(message message) error {
	m := pushover.NewMessageWithTitle(message.body, message.title)
	m.HTML = true
	_, err := w.app.SendMessage(m, w.recipient)
	if err != nil {
		return err
	}
	return nil
}

func (hook *PushoverHook) SetDelay(duration time.Duration) error {
	if duration < time.Second*1 {
		return errors.New("delay can't be less than 1 second")
	}

	hook.worker.muteDelay = duration

	return nil
}
