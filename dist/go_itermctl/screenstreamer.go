package itermctl

import (
	"context"
	"mrz.io/itermctl/iterm2"
)

// MonitorScreenUpdates subscribes to ScreenUpdateNotification and forwards each one to the returned channel.
// Subscription lasts until the given context is canceled or the conn's connection is closed. Use methods such as
// App.ScreenContents to retrieve the screen's contents.
func MonitorScreenUpdates(ctx context.Context, conn *Connection, sessionId string) (<-chan *iterm2.ScreenUpdateNotification, error) {
	notifications := make(chan *iterm2.ScreenUpdateNotification)

	req := NewNotificationRequest(true, iterm2.NotificationType_NOTIFY_ON_SCREEN_UPDATE, sessionId)
	recv, err := conn.Subscribe(ctx, req)

	if err != nil {
		return nil, err
	}

	go func() {
		for msg := range recv.Ch() {
			if msg.GetNotification().GetScreenUpdateNotification() != nil {
				notifications <- msg.GetNotification().GetScreenUpdateNotification()
			}
		}
	}()

	return notifications, nil
}
