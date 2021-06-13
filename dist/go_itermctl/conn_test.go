package itermctl_test

import (
	"fmt"
	"mrz.io/itermctl"
	"mrz.io/itermctl/iterm2"
	"sync"
	"testing"
)

func TestAcceptNotificationType(t *testing.T) {
	examples := []struct {
		msg      *iterm2.ServerOriginatedMessage
		expected bool
	}{
		{
			msg: &iterm2.ServerOriginatedMessage{
				Submessage: &iterm2.ServerOriginatedMessage_Notification{Notification: &iterm2.Notification{
					NewSessionNotification: &iterm2.NewSessionNotification{},
				}},
			},
			expected: false,
		},
		{
			msg: &iterm2.ServerOriginatedMessage{
				Submessage: &iterm2.ServerOriginatedMessage_Notification{Notification: &iterm2.Notification{
					CustomEscapeSequenceNotification: &iterm2.CustomEscapeSequenceNotification{},
				}},
			},
			expected: true,
		},
		{
			msg: &iterm2.ServerOriginatedMessage{
				Submessage: &iterm2.ServerOriginatedMessage_GetPropertyResponse{},
			},
			expected: false,
		},
	}

	for i, example := range examples {
		t.Run(fmt.Sprintf("AcceptNotificationType %d", i), func(t *testing.T) {
			f := itermctl.AcceptNotificationType(iterm2.NotificationType_NOTIFY_ON_CUSTOM_ESCAPE_SEQUENCE)
			if v := f(example.msg); v != example.expected {
				t.Fatalf("expected %t, got %t", example.expected, v)
			}
		})
	}

}

func TestNewNotificationRequest(t *testing.T) {
	examples := []struct {
		t               iterm2.NotificationType
		subscribe       bool
		session         string
		expectedSession string
	}{
		{
			subscribe:       true,
			t:               iterm2.NotificationType_NOTIFY_ON_KEYSTROKE,
			session:         itermctl.AllSessions,
			expectedSession: itermctl.AllSessions,
		},
		{
			subscribe:       false,
			t:               iterm2.NotificationType_NOTIFY_ON_PROMPT,
			session:         "",
			expectedSession: itermctl.AllSessions,
		},
		{
			subscribe:       false,
			t:               iterm2.NotificationType_NOTIFY_ON_PROMPT,
			session:         "f45b0496-d1a8-4d6e-b540-a2f3af4796ac",
			expectedSession: "f45b0496-d1a8-4d6e-b540-a2f3af4796ac",
		},
	}

	for i, example := range examples {
		t.Run(fmt.Sprintf("NewNotificationRequest %d", i), func(t *testing.T) {
			req := itermctl.NewNotificationRequest(example.subscribe, example.t, example.session)

			if req.GetNotificationType() != example.t {
				t.Fatalf("expected %s, got %s", example.t, req.GetNotificationType())
			}

			if req.GetSubscribe() != example.subscribe {
				t.Fatalf("expected %t, got %t", example.subscribe, req.GetSubscribe())
			}

			if req.GetSession() != example.expectedSession {
				t.Fatalf("expected %s, got %s", example.session, req.GetSession())
			}
		})
	}
}

func TestReceivers(t *testing.T) {
	recv1 := itermctl.NewReceiver("test1", nil)
	recv2 := itermctl.NewReceiver("test2", nil)

	collector := collectMessages([]*itermctl.Receiver{recv1, recv2})
	var receivedMessages []*iterm2.ServerOriginatedMessage
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for msg := range collector {
			receivedMessages = append(receivedMessages, msg)
		}
		wg.Done()
	}()

	recvList := &itermctl.Receivers{}
	recvList.Add(recv1)
	recvList.Add(recv2)

	id := int64(42)
	recvList.Send(&iterm2.ServerOriginatedMessage{Id: &id})

	recvList.Delete(recv1)

	id = int64(43)
	recvList.Send(&iterm2.ServerOriginatedMessage{Id: &id})

	recvList.Close()

	wg.Wait()

	if len(receivedMessages) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(receivedMessages))
	}
}

func collectMessages(receivers []*itermctl.Receiver) <-chan *iterm2.ServerOriginatedMessage {
	collector := make(chan *iterm2.ServerOriginatedMessage)
	wg := &sync.WaitGroup{}
	wg.Add(len(receivers))

	go func() {
		wg.Wait()
		close(collector)
	}()

	for _, recv := range receivers {
		go func(dst chan<- *iterm2.ServerOriginatedMessage, src <-chan *iterm2.ServerOriginatedMessage) {
			for msg := range src {
				dst <- msg
			}
			wg.Done()
		}(collector, recv.Ch())
	}

	return collector
}
