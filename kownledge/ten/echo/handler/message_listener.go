package handler

import (
	"time"
)

import (
	"github.com/apache/dubbo-getty"
)

type MessageHandler struct {
	SessionOnOpen func(session getty.Session)
}

func (h *MessageHandler) OnOpen(session getty.Session) error {
	Log.Infof("OnOpen session{%s} open", session.Stat())
	if h.SessionOnOpen != nil {
		h.SessionOnOpen(session)
	}
	return nil
}

func (h *MessageHandler) OnError(session getty.Session, err error) {
	Log.Infof("OnError session{%s} got error{%v}, will be closed.", session.Stat(), err)
}

func (h *MessageHandler) OnClose(session getty.Session) {
	Log.Infof("OnClose session{%s} is closing......", session.Stat())
}

func (h *MessageHandler) OnMessage(session getty.Session, pkg interface{}) {
	s, ok := pkg.(string)
	if !ok {
		Log.Infof("illegal package{%#v}", pkg)
		return
	}

	Log.Infof("OnMessage: %s", s)
}

func (h *MessageHandler) OnCron(session getty.Session) {
	active := session.GetActive()
	if CronPeriod < time.Since(active).Nanoseconds() {
		Log.Infof("OnCorn session{%s} timeout{%s}", session.Stat(), time.Since(active).String())
		session.Close()
	}
}
