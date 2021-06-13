package itermctl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"mrz.io/itermctl/auth"
	"mrz.io/itermctl/env"
	"mrz.io/itermctl/internal/seq"
	"mrz.io/itermctl/iterm2"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

const (
	AllSessions = "all"
)

var (
	Socket         = "~/Library/Application Support/iTerm2/private/socket"
	Subprotocol    = "api.iterm2.com"
	AppName        = "itermctl"
	LibraryVersion = "itermctl 0.0.3"
	Origin         = "ws://localhost/"
	Url            = url.URL{Scheme: "ws", Host: "localhost:1912"}

	ErrClosed                        = fmt.Errorf("connection is closed")
	ErrSessionNotFound               = fmt.Errorf("NotificationResponse_SESSION_NOT_FOUND")
	ErrRequestMalformed              = fmt.Errorf("NotificationResponse_REQUEST_MALFORMED")
	ErrNotSubscribed                 = fmt.Errorf("NotificationResponse_NOT_SUBSCRIBED")
	ErrAlreadySubscribed             = fmt.Errorf("NotificationResponse_ALREADY_SUBSCRIBED")
	ErrDuplicatedServerOriginatedRpc = fmt.Errorf("NotificationResponse_DUPLICATE_SERVER_ORIGINATED_RPC")
	ErrInvalidIdentifier             = fmt.Errorf("NotificationResponse_INVALID_IDENTIFIER")
)

func init() {
	var level log.Level
	logLevel := os.Getenv("ITERMCTL_LOG_LEVEL")

	if logLevel != "" {
		level, _ = log.ParseLevel(logLevel)
	}

	log.SetLevel(level)
}

// Receiver
type Receiver struct {
	name       string
	ch         chan *iterm2.ServerOriginatedMessage
	acceptFunc AcceptFunc
	mx         *sync.Mutex
}

func NewReceiver(name string, f AcceptFunc) *Receiver {
	r := &Receiver{name: name, mx: &sync.Mutex{}, ch: make(chan *iterm2.ServerOriginatedMessage, 100)}
	r.SetAcceptFunc(f)
	return r
}

func (r *Receiver) Ch() <-chan *iterm2.ServerOriginatedMessage {
	return r.ch
}

func (r *Receiver) Name() string {
	r.mx.Lock()
	defer r.mx.Unlock()

	return r.name
}

func (r *Receiver) SetName(n string) {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.name = n
}

func (r *Receiver) SetAcceptFunc(acceptFunc AcceptFunc) {
	r.mx.Lock()
	defer r.mx.Unlock()
	if acceptFunc == nil {
		acceptFunc = func(message *iterm2.ServerOriginatedMessage) bool {
			return true
		}
	}
	r.acceptFunc = acceptFunc
}

func (r *Receiver) Accept(msg *iterm2.ServerOriginatedMessage) bool {
	r.mx.Lock()
	defer r.mx.Unlock()
	return r.acceptFunc(msg)
}

type Receivers []*Receiver

func (r *Receivers) Close() {
	for _, recv := range *r {
		close(recv.ch)
	}

	*r = []*Receiver{}
}

func (r *Receivers) Send(msg *iterm2.ServerOriginatedMessage) {
	if len(*r) == 0 {
		log.Warnf("message ID %d: lost, no receivers registered", msg.GetId())
		return
	}

	for _, recv := range *r {
		if !recv.Accept(msg) {
			log.Debugf("message ID %d: not accepted by %q", msg.GetId(), recv.Name())
			continue
		}

		select {
		case recv.ch <- msg:
			log.Debugf("message ID %d: sent to %q", msg.GetId(), recv.Name())
		case <-time.After(1 * time.Second):
			log.Errorf("message ID %d: time out sending to %q", msg.GetId(), recv.Name())
		}
	}
}

func (r *Receivers) Add(recv *Receiver) {
	*r = append(*r, recv)
}

func (r *Receivers) Delete(other *Receiver) {
	var tmp []*Receiver

	for _, recv := range *r {
		if recv == other {
			close(recv.ch)
		} else {
			tmp = append(tmp, recv)
		}
	}

	*r = tmp
}

// AcceptFunc is the function given to Connection.Receiver() to filter out uninteresting ServerOriginatedMessages.
type AcceptFunc func(msg *iterm2.ServerOriginatedMessage) bool

// AcceptNotificationType filters ServerOriginatedMessages whose submessage is a Notification of the given type.
func AcceptNotificationType(t iterm2.NotificationType) AcceptFunc {
	return func(msg *iterm2.ServerOriginatedMessage) bool {
		actualType := getNotificationType(msg.GetNotification())
		return actualType == t
	}
}

func acceptMessageId(msgId int64) AcceptFunc {
	return func(msg *iterm2.ServerOriginatedMessage) bool {
		return msg.GetId() == msgId
	}
}

type Transaction struct {
	cancelFunc context.CancelFunc
	errCh      chan error
}

func newTransaction(ctx context.Context, conn *Connection) *Transaction {
	ctx, cancel := context.WithCancel(ctx)

	tx := &Transaction{
		cancelFunc: cancel,
		errCh:      make(chan error),
	}

	go func() {
		<-ctx.Done()
		begin := false
		endMessage := &iterm2.ClientOriginatedMessage{
			Submessage: &iterm2.ClientOriginatedMessage_TransactionRequest{
				TransactionRequest: &iterm2.TransactionRequest{Begin: &begin},
			},
		}

		if _, err := conn.GetResponse(context.Background(), endMessage); err != nil {
			tx.errCh <- fmt.Errorf("end transaction: %w", err)
		}

		close(tx.errCh)
	}()

	return tx
}

func (t *Transaction) End() error {
	t.cancelFunc()
	return <-t.errCh
}

// GetCredentialsAndConnect checks if iTerm2 is configured to require authentication, retrieves the cookie and key if
// necessary, and then establishes the connection to iTerm2's websocket.
func GetCredentialsAndConnect(appName string, active bool) (*Connection, error) {
	if appName == "" {
		appName = AppName
	}

	var cookie, key string
	var err error

	if cookie, key, err = env.CookieAndKey(); err != nil {
		if err = auth.Disabled(); err != nil {
			cookie, key, err = auth.RequestCookieAndKey(appName, active)
			if err != nil {
				return nil, err
			}
		}
	}

	return Connect(appName, cookie, key)
}

// Connect connects to iTerm2's websocket using the optional credentials. AppName is used as a default app name if none
// is given.
func Connect(appName, cookie, key string) (*Connection, error) {
	if appName == "" {
		appName = AppName
	}

	socket, err := homedir.Expand(Socket)
	if err != nil {
		return nil, fmt.Errorf("connect: cannot expand %s: %w", Socket, err)
	}

	var headers = map[string][]string{
		"Origin":                   {Origin},
		"x-iterm2-disable-auth-ui": {"false"},
		"x-iterm2-advisory-name":   {appName},
		"x-iterm2-library-version": {LibraryVersion},
	}

	if cookie != "" {
		headers["x-iterm2-cookie"] = []string{cookie}
	}
	if key != "" {
		headers["x-iterm2-key"] = []string{key}
	}

	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5 * time.Second,
		Subprotocols:     []string{Subprotocol},
		NetDial: func(network, addr string) (net.Conn, error) {
			return net.Dial("unix", socket)
		},
	}

	ws, _, err := dialer.Dial(Url.String(), headers)

	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	return NewConnection(ws), err
}

// Connection represents a connection to iTerm2, providing basic methods to Send and read messages.
type Connection struct {
	outgoingMessages chan *iterm2.ClientOriginatedMessage
	incomingMessages <-chan *iterm2.ServerOriginatedMessage
	addReceivers     chan *Receiver
	deleteReceivers  chan *Receiver
	closed           bool
	closedLock       *sync.Mutex
	closeCtx         context.Context
	closeFunc        context.CancelFunc

	websocket *websocket.Conn
}

// NewConnection creates a new Connection wrapping around a *websocket.Conn.
func NewConnection(ws *websocket.Conn) *Connection {
	closeCtx, closeFunc := context.WithCancel(context.Background())
	conn := &Connection{
		addReceivers:     make(chan *Receiver),
		deleteReceivers:  make(chan *Receiver),
		outgoingMessages: make(chan *iterm2.ClientOriginatedMessage),
		closed:           false,
		closedLock:       &sync.Mutex{},
		closeCtx:         closeCtx,
		closeFunc:        closeFunc,

		websocket: ws,
	}

	conn.incomingMessages = conn.read()

	go func() {
		var receivers Receivers

		for {
			select {
			case <-conn.closeCtx.Done():
				goto shutdown
			case recv := <-conn.addReceivers:
				receivers.Add(recv)
			case recv := <-conn.deleteReceivers:
				receivers.Delete(recv)
			case msg, ok := <-conn.incomingMessages:
				if !ok {
					goto shutdown
				}

				receivers.Send(msg)
			case msg := <-conn.outgoingMessages:
				if msg.GetId() == 0 {
					msg.Id = seq.MessageId.Next()
				}

				if err := conn.write(msg); err != nil {
					log.Error(err)
				}
			}
		}

	shutdown:
		conn.closedLock.Lock()
		defer conn.closedLock.Unlock()
		conn.closed = true

		close(conn.addReceivers)
		close(conn.deleteReceivers)
		close(conn.outgoingMessages)

		receivers.Close()

		if err := conn.websocket.Close(); err != nil {
			log.Errorf("sendCloseRequest: %s", err)
		}
	}()

	return conn
}

func (conn *Connection) read() <-chan *iterm2.ServerOriginatedMessage {
	messages := make(chan *iterm2.ServerOriginatedMessage, 1000)

	go func() {
		for {
			_, data, err := conn.websocket.ReadMessage()
			if err != nil {
				log.Error()
				break
			}

			msg := &iterm2.ServerOriginatedMessage{}
			if err := proto.Unmarshal(data, msg); err != nil {
				log.Errorf("read: could not unmarshal message: %s", err)
				continue
			}

			messages <- msg
		}
		close(messages)
	}()

	return messages
}

func (conn *Connection) write(msg *iterm2.ClientOriginatedMessage) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("write: could not marshal message: %w", err)
	}

	if err := conn.websocket.WriteMessage(websocket.BinaryMessage, data); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

// Wait blocks until the conn's shuts down.
func (conn *Connection) Wait() {
	<-conn.closeCtx.Done()
}

// Close initiates the connection's shutdown, causing all the receivers channels to be closed. It will also sendCloseRequest the
// underlying websocket.
func (conn *Connection) Close() {
	conn.closeFunc()
}

// Receiver returns a receiver for ServerOriginatedMessages. Messages can be read from the receiver's Ch() until the
// Connection is closed or the context is canceled. A context should be given only to interrupt receiving before the
// Connection is closed, and should not be the same as the one used to cancel the Connection. The receiver will receive
// a copy of any ServerOriginatedMessage being shipped on the Connection, but an AcceptFunc can be given to exclude
// uninteresting messages.
func (conn *Connection) Receiver(ctx context.Context, name string, f AcceptFunc) (*Receiver, error) {
	conn.closedLock.Lock()
	defer conn.closedLock.Unlock()
	if conn.closed {
		return nil, ErrClosed
	}

	recv := NewReceiver(name, f)

	if ctx != nil {
		go func() {
			<-ctx.Done()
			conn.closedLock.Lock()
			defer conn.closedLock.Unlock()
			if !conn.closed {
				conn.deleteReceivers <- recv
			}
		}()
	}

	conn.addReceivers <- recv
	return recv, nil
}

// Send sends a message to iTerm2, without waiting for a response. ErrClosed is returned when Send is called after
// the connection was closed.
func (conn *Connection) Send(msg *iterm2.ClientOriginatedMessage) error {
	conn.closedLock.Lock()
	defer conn.closedLock.Unlock()

	if conn.closed {
		return ErrClosed
	}

	conn.outgoingMessages <- msg
	return nil
}

// GetResponse sends a message to iTerm2, and waits for a message to be read from src and returns it. If the message is
// an error from iTerm2, a nil message and an error are returned. A nil message with an error will also be returned if
// the context is canceled before the response is received.
func (conn *Connection) GetResponse(ctx context.Context, req *iterm2.ClientOriginatedMessage) (*iterm2.ServerOriginatedMessage, error) {
	src, err := conn.request(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get response: %w", err)
	}

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("get response: %w", ctx.Err())
	case resp := <-src:
		if resp != nil {
			if resp.GetError() != "" {
				return nil, fmt.Errorf("get response: %s", resp.GetError())
			}
		}
		return resp, nil
	}
}

func (conn *Connection) request(ctx context.Context, req *iterm2.ClientOriginatedMessage) (<-chan *iterm2.ServerOriginatedMessage, error) {
	if req.Id == nil {
		req.Id = seq.MessageId.Next()
	}

	respCh := make(chan *iterm2.ServerOriginatedMessage)
	recvCtx, cancelRecv := context.WithCancel(ctx)

	recv, err := conn.Receiver(
		recvCtx,
		fmt.Sprintf("receive Message ID %d", req.GetId()),
		acceptMessageId(req.GetId()),
	)

	if err != nil {
		cancelRecv()
		return nil, err
	}

	err = conn.Send(req)
	if err != nil {
		cancelRecv()
		return nil, err
	}

	go func() {
		defer cancelRecv()
		select {
		case <-recvCtx.Done():
			break
		case resp := <-recv.Ch():
			if resp != nil {
				respCh <- resp
			}
			break
		}
		close(respCh)
	}()

	return respCh, nil
}

// InvokeFunction invokes an RPC function and unmarshalls the result into target. If iTerm2's response to the invocation
// is an error, target is left untouched and an error is returned.
func (conn *Connection) InvokeFunction(invocation string, target interface{}) error {
	req := &iterm2.ClientOriginatedMessage{
		Submessage: &iterm2.ClientOriginatedMessage_InvokeFunctionRequest{
			InvokeFunctionRequest: &iterm2.InvokeFunctionRequest{
				Context:    &iterm2.InvokeFunctionRequest_App_{},
				Invocation: &invocation,
			},
		},
	}

	resp, err := conn.GetResponse(context.Background(), req)
	if resp == nil {
		return err
	}

	if invocationErr := resp.GetInvokeFunctionResponse().GetError(); invocationErr != nil {
		return fmt.Errorf("%s: %s", invocationErr.GetStatus(), invocationErr.GetErrorReason())
	}

	jsonResult := resp.GetInvokeFunctionResponse().GetSuccess().GetJsonResult()
	if err := json.Unmarshal([]byte(jsonResult), &target); err != nil {
		return fmt.Errorf("could not unmarshal invocation target: %w", err)
	}

	return nil
}

// Subscribe uses the given NotificationRequest to subscribe with iTerm2, and returns a channel from which notifications
// of requested type can be read. The NotificationRequest will be modified to ensure the Subscribe field is set to true.
// The subscription will be canceled automatically as soon as the context is canceled. The subscription lasts until the
// give context is canceled or the conn connection is closed.
func (conn *Connection) Subscribe(ctx context.Context, req *iterm2.NotificationRequest) (*Receiver, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	recv, err := conn.Receiver(ctx,
		fmt.Sprintf("receive %s", req.NotificationType.String()),
		AcceptNotificationType(req.GetNotificationType()),
	)

	if err != nil {
		return nil, fmt.Errorf("subscribe: %w", err)
	}

	subscribe := true
	req.Subscribe = &subscribe

	msg := &iterm2.ClientOriginatedMessage{}
	msg.Submessage = &iterm2.ClientOriginatedMessage_NotificationRequest{
		NotificationRequest: req,
	}

	resp, err := conn.GetResponse(context.Background(), msg)
	if err != nil {
		return nil, fmt.Errorf("subscribe: %w", err)
	}

	subscriptionErr := getSubscriptionStatusError(resp)
	if subscriptionErr != nil {
		if subscriptionErr == ErrAlreadySubscribed {
			return recv, nil
		}

		return nil, fmt.Errorf("subscribe: %w", subscriptionErr)
	}

	go func() {
		<-ctx.Done()

		unsubReq := NewNotificationRequest(false, req.GetNotificationType(), req.GetSession())
		unsubReq.Arguments = req.Arguments

		unsubErr := conn.unsubscribe(unsubReq)
		if unsubErr != nil {
			log.Errorf("unsubscribe %s: %s", req.GetNotificationType(), unsubErr)
		} else {
			log.Debugf("unsubscribe %s: successful", req.GetNotificationType())
		}
	}()

	return recv, nil
}

func (conn *Connection) unsubscribe(req *iterm2.NotificationRequest) error {
	msg := &iterm2.ClientOriginatedMessage{
		Submessage: &iterm2.ClientOriginatedMessage_NotificationRequest{
			NotificationRequest: req,
		},
	}

	resp, err := conn.GetResponse(context.Background(), msg)
	if err != nil {
		return err
	}

	return getSubscriptionStatusError(resp)
}

// Transaction start a transaction, a sequence of API calls can occur without anything else happening in between.
// Note that this effectively freezes iTerm2 until Transaction.End is called.
// See https://iterm2.com/python-api/transaction.html.
func (conn *Connection) Transaction() (*Transaction, error) {
	begin := true
	beginMessage := &iterm2.ClientOriginatedMessage{
		Submessage: &iterm2.ClientOriginatedMessage_TransactionRequest{
			TransactionRequest: &iterm2.TransactionRequest{Begin: &begin},
		},
	}

	ctx := context.Background()

	resp, err := conn.GetResponse(ctx, beginMessage)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	if resp.GetTransactionResponse().GetStatus() != iterm2.TransactionResponse_OK {
		return nil, fmt.Errorf("begin transaction: %s", resp.GetTransactionResponse().GetStatus().String())
	}

	return newTransaction(ctx, conn), nil
}

// NewNotificationRequest creates a notification request to subscribe or unsubscribe for the given notification
// type. If an empty sessionId is given, the subscription is created for all sessions.
func NewNotificationRequest(subscribe bool, nt iterm2.NotificationType, sessionId string) *iterm2.NotificationRequest {
	if sessionId == "" {
		sessionId = AllSessions
	}

	return &iterm2.NotificationRequest{
		Session:          &sessionId,
		Subscribe:        &subscribe,
		NotificationType: &nt,
	}
}

func getNotificationType(n *iterm2.Notification) iterm2.NotificationType {
	if n.GetTerminateSessionNotification() != nil {
		return iterm2.NotificationType_NOTIFY_ON_TERMINATE_SESSION
	} else if n.GetNewSessionNotification() != nil {
		return iterm2.NotificationType_NOTIFY_ON_NEW_SESSION
	} else if n.GetCustomEscapeSequenceNotification() != nil {
		return iterm2.NotificationType_NOTIFY_ON_CUSTOM_ESCAPE_SEQUENCE
	} else if n.GetServerOriginatedRpcNotification() != nil {
		return iterm2.NotificationType_NOTIFY_ON_SERVER_ORIGINATED_RPC
	} else if n.GetBroadcastDomainsChanged() != nil {
		return iterm2.NotificationType_NOTIFY_ON_BROADCAST_CHANGE
	} else if n.GetFocusChangedNotification() != nil {
		return iterm2.NotificationType_NOTIFY_ON_FOCUS_CHANGE
	} else if n.GetKeystrokeNotification() != nil {
		return iterm2.NotificationType_NOTIFY_ON_KEYSTROKE
	} else if n.GetProfileChangedNotification() != nil {
		return iterm2.NotificationType_NOTIFY_ON_PROFILE_CHANGE
	} else if n.GetScreenUpdateNotification() != nil {
		return iterm2.NotificationType_NOTIFY_ON_SCREEN_UPDATE
	} else if n.GetLayoutChangedNotification() != nil {
		return iterm2.NotificationType_NOTIFY_ON_LAYOUT_CHANGE
	} else if n.GetVariableChangedNotification() != nil {
		return iterm2.NotificationType_NOTIFY_ON_VARIABLE_CHANGE
	} else if n.GetPromptNotification() != nil {
		return iterm2.NotificationType_NOTIFY_ON_PROMPT
	}
	return 0
}

func getSubscriptionStatusError(resp *iterm2.ServerOriginatedMessage) error {
	switch resp.GetNotificationResponse().GetStatus() {
	case iterm2.NotificationResponse_SESSION_NOT_FOUND:
		return ErrSessionNotFound
	case iterm2.NotificationResponse_REQUEST_MALFORMED:
		return ErrRequestMalformed
	case iterm2.NotificationResponse_DUPLICATE_SERVER_ORIGINATED_RPC:
		return ErrDuplicatedServerOriginatedRpc
	case iterm2.NotificationResponse_INVALID_IDENTIFIER:
		return ErrInvalidIdentifier
	case iterm2.NotificationResponse_ALREADY_SUBSCRIBED:
		return ErrAlreadySubscribed
	case iterm2.NotificationResponse_NOT_SUBSCRIBED:
		return ErrNotSubscribed
	}
	return nil
}
