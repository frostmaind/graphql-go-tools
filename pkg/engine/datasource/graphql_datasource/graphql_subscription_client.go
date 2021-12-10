package graphql_datasource

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/buger/jsonparser"
	"github.com/cespare/xxhash/v2"
	"github.com/jensneuse/abstractlogger"
	"nhooyr.io/websocket"
)

const (
	connectionInitMessage = `{"type":"connection_init","payload":%s}`
	startMessage          = `{"type":"start","id":"%s","payload":%s}`
	stopMessage           = `{"type":"stop","id":"%s"}`
	internalError         = `{"errors":[{"message":"connection error"}]}`
	connectionError       = `{"errors":[{"message":"connection error"}]}`
)

const connectionInitTimeout = 3 * time.Second

// WebSocketGraphQLSubscriptionClient is a WebSocket client that allows running multiple subscriptions via the same WebSocket Connection
// It takes care of de-duplicating WebSocket connections to the same origin under certain circumstances
// If Hash(URL,Body,Headers) result in the same result, an existing WS connection is re-used
type WebSocketGraphQLSubscriptionClient struct {
	httpClient *http.Client
	ctx        context.Context
	log        abstractlogger.Logger
	hashPool   sync.Pool
	handlers   map[uint64]*connectionHandler
	handlersMu sync.Mutex

	readTimeout time.Duration
	readLimit   int64
}

type Options func(options *opts)

func WithLogger(log abstractlogger.Logger) Options {
	return func(options *opts) {
		options.log = log
	}
}

func WithReadTimeout(timeout time.Duration) Options {
	return func(options *opts) {
		options.readTimeout = timeout
	}
}

func WithReadLimit(limit int64) Options {
	return func(options *opts) {
		options.readLimit = limit
	}
}

type opts struct {
	readTimeout time.Duration
	log         abstractlogger.Logger
	readLimit   int64
}

func NewWebSocketGraphQLSubscriptionClient(httpClient *http.Client, ctx context.Context, options ...Options) *WebSocketGraphQLSubscriptionClient {
	op := &opts{
		readTimeout: time.Second,
		log:         abstractlogger.NoopLogger,
	}
	for _, option := range options {
		option(op)
	}
	return &WebSocketGraphQLSubscriptionClient{
		httpClient:  httpClient,
		ctx:         ctx,
		handlers:    map[uint64]*connectionHandler{},
		log:         op.log,
		readTimeout: op.readTimeout,
		readLimit:   op.readLimit,
		hashPool: sync.Pool{
			New: func() interface{} {
				return xxhash.New()
			},
		},
	}
}

// Subscribe initiates a new GraphQL Subscription with the origin
// Each WebSocket (WS) to an origin is uniquely identified by the Hash(URL,Headers,Body)
// If an existing WS with the same ID (Hash) exists, it is being re-used
// If no connection exists, the client initiates a new one and sends the "init" and "connection ack" messages
func (c *WebSocketGraphQLSubscriptionClient) Subscribe(ctx context.Context, options GraphQLSubscriptionOptions, next chan<- []byte) error {

	handlerID, err := c.generateHandlerIDHash(options)
	if err != nil {
		return err
	}

	sub := subscription{
		ctx:     ctx,
		options: options,
		next:    next,
	}

	c.handlersMu.Lock()
	defer c.handlersMu.Unlock()
	handler, exists := c.handlers[handlerID]
	if exists {
		select {
		case handler.subscribeCh <- sub:
		case <-handler.connectionDoneCh:
			c.handleClosedConnectionHandler(ctx, next)
		case <-ctx.Done():
			close(next)
		}
		return nil
	}

	if options.Header == nil {
		options.Header = http.Header{}
	}

	initialPayload, err := json.Marshal(options.Header)
	if err != nil {
		return err
	}

	options.Header.Set("Sec-WebSocket-Protocol", "graphql-ws")
	options.Header.Set("Sec-WebSocket-Version", "13")

	connectionInitCtx, cancelConnectionInitCtx := context.WithTimeout(ctx, connectionInitTimeout)
	defer cancelConnectionInitCtx()

	conn, upgradeResponse, err := websocket.Dial(connectionInitCtx, options.URL, &websocket.DialOptions{
		HTTPClient:      c.httpClient,
		HTTPHeader:      options.Header,
		CompressionMode: websocket.CompressionDisabled,
		Subprotocols:    []string{"graphql-ws"},
	})
	if err != nil {
		return err
	}
	if upgradeResponse.StatusCode != http.StatusSwitchingProtocols {
		return fmt.Errorf("upgrade unsuccessful")
	}

	if c.readLimit != 0 {
		conn.SetReadLimit(c.readLimit)
	}

	// init + ack
	initialMessage := fmt.Sprintf(connectionInitMessage, string(initialPayload))
	err = conn.Write(connectionInitCtx, websocket.MessageText, []byte(initialMessage))
	if err != nil {
		return err
	}
	msgType, connectionAckMsg, err := conn.Read(connectionInitCtx)
	if err != nil {
		return err
	}
	if msgType != websocket.MessageText {
		return fmt.Errorf("unexpected msg type")
	}
	connectionAck, err := jsonparser.GetString(connectionAckMsg, "type")
	if err != nil {
		return err
	}
	if connectionAck != "connection_ack" {
		return fmt.Errorf("expected connection_ack, got: %s", connectionAck)
	}

	handler = newConnectionHandler(c.ctx, conn, c.readTimeout, c.log)
	c.handlers[handlerID] = handler

	go func(handlerID uint64) {
		handler.startBlocking(sub)
		c.handlersMu.Lock()
		delete(c.handlers, handlerID)
		c.handlersMu.Unlock()
	}(handlerID)

	return nil
}

// generateHandlerIDHash generates a Hash based on: URL and Headers to uniquely identify Upgrade Requests
func (c *WebSocketGraphQLSubscriptionClient) generateHandlerIDHash(options GraphQLSubscriptionOptions) (uint64, error) {
	var (
		err error
	)
	xxh := c.hashPool.Get().(*xxhash.Digest)
	defer c.hashPool.Put(xxh)
	xxh.Reset()

	_, err = xxh.WriteString(options.URL)
	if err != nil {
		return 0, err
	}
	err = options.Header.Write(xxh)
	if err != nil {
		return 0, err
	}

	return xxh.Sum64(), nil
}

func (c *WebSocketGraphQLSubscriptionClient) handleClosedConnectionHandler(ctx context.Context, next chan<- []byte) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	select {
	case next <- []byte(connectionError):
	case <-ctx.Done():
	}

	close(next)
}

func newConnectionHandler(ctx context.Context, conn *websocket.Conn, readTimeout time.Duration, log abstractlogger.Logger) *connectionHandler {
	return &connectionHandler{
		conn:               conn,
		ctx:                ctx,
		log:                log,
		subscribeCh:        make(chan subscription),
		nextSubscriptionID: 0,
		subscriptions:      map[string]subscription{},
		readTimeout:        readTimeout,
		connectionDoneCh:   make(chan struct{}),
	}
}

// connectionHandler is responsible for handling a connection to an origin
// it is responsible for managing all subscriptions using the underlying WebSocket connection
// if all Subscriptions are complete or cancelled/unsubscribed the handler will terminate
type connectionHandler struct {
	conn               *websocket.Conn
	ctx                context.Context
	log                abstractlogger.Logger
	subscribeCh        chan subscription
	nextSubscriptionID int
	subscriptions      map[string]subscription
	readTimeout        time.Duration
	connectionDoneCh   chan struct{}
}

type subscription struct {
	ctx     context.Context
	options GraphQLSubscriptionOptions
	next    chan<- []byte
}

// startBlocking starts the single threaded event loop of the handler
// if the global context returns or the websocket connection is terminated, it will stop
func (h *connectionHandler) startBlocking(sub subscription) {
	readCtx, cancel := context.WithCancel(h.ctx)
	defer func() {
		h.unsubscribeAllAndCloseConn()
		cancel()
		close(h.connectionDoneCh)
	}()
	h.subscribe(sub)
	dataCh := make(chan []byte)
	errCh := make(chan error, 1)

	go func() { errCh <- h.readBlocking(readCtx, dataCh) }()

	for {
		if h.ctx.Err() != nil {
			return
		}
		hasActiveSubscriptions := h.checkActiveSubscriptions()
		if !hasActiveSubscriptions {
			return
		}
		select {
		case <-time.After(h.readTimeout):
			continue
		case readErr := <-errCh:
			if readErr != nil {
				h.log.Info("Got an error from WS reader", abstractlogger.String("message", readErr.Error()))
				h.handleMessageTypeConnectionError()
				return
			}
		case sub = <-h.subscribeCh:
			h.subscribe(sub)
		case data := <-dataCh:
			messageType, err := jsonparser.GetString(data, "type")
			if err != nil {
				continue
			}
			switch messageType {
			case "data":
				h.handleMessageTypeData(data)
			case "complete":
				h.handleMessageTypeComplete(data)
			case "connection_error":
				h.handleMessageTypeConnectionError()
				return
			case "error":
				h.handleMessageTypeError(data)
				continue
			default:
				continue
			}
		}
	}
}

// readBlocking is a dedicated loop running in a separate goroutine
// because the library "nhooyr.io/websocket" doesn't allow reading with a context with Timeout
// we'll block forever on reading until the context of the connectionHandler stops
func (h *connectionHandler) readBlocking(ctx context.Context, dataCh chan []byte) error {
	for {
		msgType, data, err := h.conn.Read(ctx)
		if ctx.Err() != nil {
			return nil
		}
		if err != nil {
			return err
		}
		if msgType != websocket.MessageText {
			continue
		}
		select {
		case dataCh <- data:
		case <-ctx.Done():
			return nil
		}
	}
}

func (h *connectionHandler) unsubscribeAllAndCloseConn() {
	for id := range h.subscriptions {
		h.unsubscribe(id)
	}
	_ = h.conn.Close(websocket.StatusNormalClosure, "")
}

// subscribe adds a new subscription to the connectionHandler and sends the startMessage to the origin
func (h *connectionHandler) subscribe(sub subscription) {
	graphQLBody, err := json.Marshal(sub.options.Body)
	if err != nil {
		return
	}

	h.nextSubscriptionID++

	subscriptionID := strconv.Itoa(h.nextSubscriptionID)

	startRequest := fmt.Sprintf(startMessage, subscriptionID, string(graphQLBody))
	err = h.conn.Write(h.ctx, websocket.MessageText, []byte(startRequest))
	if err != nil {
		return
	}

	h.subscriptions[subscriptionID] = sub
}

func (h *connectionHandler) handleMessageTypeData(data []byte) {
	id, err := jsonparser.GetString(data, "id")
	if err != nil {
		return
	}
	sub, ok := h.subscriptions[id]
	if !ok {
		return
	}
	payload, _, _, err := jsonparser.Get(data, "payload")
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(h.ctx, time.Second*5)
	defer cancel()

	select {
	case <-ctx.Done():
	case sub.next <- payload:
	case <-sub.ctx.Done():
	}
}

func (h *connectionHandler) handleMessageTypeConnectionError() {
	for _, sub := range h.subscriptions {
		ctx, cancel := context.WithTimeout(h.ctx, time.Second*5)
		select {
		case sub.next <- []byte(connectionError):
			cancel()
			continue
		case <-ctx.Done():
			cancel()
			continue
		}
	}
}

func (h *connectionHandler) handleMessageTypeComplete(data []byte) {
	id, err := jsonparser.GetString(data, "id")
	if err != nil {
		return
	}
	sub, ok := h.subscriptions[id]
	if !ok {
		return
	}
	close(sub.next)
	delete(h.subscriptions, id)
}

func (h *connectionHandler) handleMessageTypeError(data []byte) {
	id, err := jsonparser.GetString(data, "id")
	if err != nil {
		return
	}
	sub, ok := h.subscriptions[id]
	if !ok {
		return
	}
	value, valueType, _, err := jsonparser.Get(data, "payload")
	if err != nil {
		sub.next <- []byte(internalError)
		return
	}
	switch valueType {
	case jsonparser.Array:
		response := []byte(`{}`)
		response, err = jsonparser.Set(response, value, "errors")
		if err != nil {
			sub.next <- []byte(internalError)
			return
		}
		sub.next <- response
	case jsonparser.Object:
		response := []byte(`{"errors":[]}`)
		response, err = jsonparser.Set(response, value, "errors", "[0]")
		if err != nil {
			sub.next <- []byte(internalError)
			return
		}
		sub.next <- response
	default:
		sub.next <- []byte(internalError)
	}
}

func (h *connectionHandler) unsubscribe(subscriptionID string) {
	sub, ok := h.subscriptions[subscriptionID]
	if !ok {
		return
	}
	close(sub.next)
	delete(h.subscriptions, subscriptionID)
	stopRequest := fmt.Sprintf(stopMessage, subscriptionID)
	_ = h.conn.Write(h.ctx, websocket.MessageText, []byte(stopRequest))
}

func (h *connectionHandler) checkActiveSubscriptions() (hasActiveSubscriptions bool) {
	for id, sub := range h.subscriptions {
		if sub.ctx.Err() != nil {
			h.unsubscribe(id)
		}
	}
	return len(h.subscriptions) != 0
}
