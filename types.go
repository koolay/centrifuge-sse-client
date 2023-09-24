package sseclient

import (
	"encoding/json"
	"fmt"
)

type APIMethod string

const (
	publishMethod  APIMethod = "publish"
	presenceMethod APIMethod = "presence"
)

type ErrBadStatus struct {
	Code int
	Body string
}

func (err *ErrBadStatus) Error() string {
	return fmt.Sprintf("bad status code: %d, text:%s", err.Code, err.Body)
}

type ErrLogic struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (err *ErrLogic) Error() string {
	return fmt.Sprintf("logic error, code: %d, msg: %s ", err.Code, err.Message)
}

// https://github.com/centrifugal/centrifugo/blob/5a7b9176b3be622af1a21e134b880e1cc5e4b4e4/internal/jwtverify/token_verifier_jwt.go#L377

type BoolValue struct {
	Value bool `json:"value,omitempty"`
}

type SubscribeOptionOverride struct {
	// Presence turns on participating in channel presence.
	Presence *BoolValue `json:"presence,omitempty"`
	// JoinLeave enables sending Join and Leave messages for this client in channel.
	JoinLeave *BoolValue `json:"join_leave,omitempty"`
	// ForcePushJoinLeave forces sending join/leave for this client.
	ForcePushJoinLeave *BoolValue `json:"force_push_join_leave,omitempty"`
	// ForcePositioning on says that client will additionally sync its position inside
	// a stream to prevent message loss. Make sure you are enabling ForcePositioning in channels
	// that maintain Publication history stream. When ForcePositioning is on  Centrifuge will
	// include StreamPosition information to subscribe response - for a client to be able
	// to manually track its position inside a stream.
	ForcePositioning *BoolValue `json:"force_positioning,omitempty"`
	// ForceRecovery turns on recovery option for a channel. In this case client will try to
	// recover missed messages automatically upon resubscribe to a channel after reconnect
	// to a server. This option also enables client position tracking inside a stream
	// (like ForcePositioning option) to prevent occasional message loss. Make sure you are using
	// ForceRecovery in channels that maintain Publication history stream.
	ForceRecovery *BoolValue `json:"force_recovery,omitempty"`
}

type SubscribeOptions struct {
	// Info defines custom channel information, zero value means no channel information.
	Info json.RawMessage `json:"info,omitempty"`
	// Base64Info is like Info but for binary.
	Base64Info string `json:"b64info,omitempty"`
	// Data to send to a client with Subscribe Push.
	Data json.RawMessage `json:"data,omitempty"`
	// Base64Data is like Data but for binary data.
	Base64Data string `json:"b64data,omitempty"`
	// Override channel options can contain channel options overrides.
	Override *SubscribeOptionOverride `json:"override,omitempty"`
}

// claimsSub define per-subscription options.
type claimsSub map[string]interface{}

type APIRequest struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type PublishParams struct {
	Channel     string          `json:"channel,omitempty"`
	Data        json.RawMessage `json:"data,omitempty"`
	SkipHistory bool            `json:"skip_history,omitempty"`
}

type PublishResponse struct {
	Error  ErrLogic `json:"error,omitempty"`
	Result struct {
		Offset int    `json:"offset,omitempty"`
		Epoch  string `json:"epoch,omitempty"`
	} `json:"result,omitempty"`
}

type EventData struct {
	Channel string `json:"channel,omitempty"`
	Pub     struct {
		Data json.RawMessage `json:"data,omitempty"`
	} `json:"pub,omitempty"`
}

// https://centrifugal.dev/docs/3/server/server_api#presence
type GetOnlineParams struct {
	Channel string `json:"channel,omitempty"`
}
type PresenceItem struct {
	Client string `json:"client,omitempty"`
	User   string `json:"user,omitempty"`
}

type GetOnlineResponse struct {
	Error  ErrLogic `json:"error,omitempty"`
	Result struct {
		Presence map[string]PresenceItem `json:"presence,omitempty"`
	} `json:"result,omitempty"`
}
