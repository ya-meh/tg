package telebot

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

type AlbumHandlerFunc func(cs []Context) error

type albumHandler struct {
	Group   *Group
	Func    AlbumHandlerFunc
	Timeout time.Duration

	data         map[string][]Context
	registerLock sync.Mutex
}

func (handler *albumHandler) mediaGroupToId(msg *Message) string {
	if msg.AlbumID != "" {
		return msg.AlbumID
	} else {
		return fmt.Sprintf("%d_%d", msg.Chat.ID, msg.ID)
	}
}

func (handler *albumHandler) Register(ctx Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	handler.registerLock.Lock()
	defer handler.registerLock.Unlock()

	id := handler.mediaGroupToId(ctx.Message())
	if _, contains := handler.data[id]; !contains {
		handler.data[id] = []Context{ctx}

		go handler.delayHandling(ctx, id)
	} else {
		handler.data[id] = append(handler.data[id], ctx)
	}

	return nil
}

func (handler *albumHandler) delayHandling(ctx Context, id string) {
	message := ctx.Message()
	defer func() {
		delete(handler.data, id)
		if r := recover(); r != nil {
			ctx.Bot().OnError(errors.New(fmt.Sprintf("%v", r)), ctx)
		}
	}()
	if message.AlbumID != "" { // no need to delay handling of single medias
		time.Sleep(handler.Timeout)
	}
	contexts := handler.data[handler.mediaGroupToId(message)]
	sort.Slice(contexts, func(i, j int) bool {
		return contexts[i].Message().ID < contexts[j].Message().ID
	})

	if err := handler.Func(contexts); err != nil {
		ctx.Bot().OnError(err, ctx)
	}
}

// HandleAlbum uses options to specify tg.OnPhoto/tg.OnVideo/etc, or time.Duration to set timeout for handling media.
// Default timeout is 0.333 sec. If you want to use MiddlewareFunc, create a Group.
func (b *Bot) HandleAlbum(handler AlbumHandlerFunc, options ...interface{}) {
	b.Group().HandleAlbum(handler, options...)
}

// HandleAlbum uses options to specify tg.OnPhoto/tg.OnVideo/etc, or time.Duration to set timeout for handling media
func (g *Group) HandleAlbum(handler AlbumHandlerFunc, options ...interface{}) {
	endpoints := []string{}
	timeout := time.Millisecond * 333
	if len(options) == 0 {
		endpoints = append(endpoints, OnMedia)
	} else {
		for _, option := range options {
			switch option.(type) {
			case string:
				endpoints = append(endpoints, option.(string))
			case time.Duration:
				timeout = option.(time.Duration)
			}
		}
	}

	albumHandler := albumHandler{
		Group:        g,
		Func:         handler,
		Timeout:      timeout,
		data:         map[string][]Context{},
		registerLock: sync.Mutex{},
	}

	for _, endpoint := range endpoints {
		albumHandler.Group.Handle(endpoint, func(ctx Context) error { return albumHandler.Register(ctx) })
	}
}
