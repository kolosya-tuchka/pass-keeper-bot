package telegram

import (
	"context"
	"fmt"
	"log"
	"pass-keeper-bot/storage"
	"strings"
)

var (
	userContexts = map[int64]interface{}{}
	userStates   = map[int64]userState{}
	userCommands = map[userState]func(ctx context.Context) error{
		Unknown:    cmdBegin,
		Begin:      cmdBegin,
		Get:        cmdGet,
		Del:        cmdDel,
		SetService: cmdSetService,
		SetLogin:   cmdSetLogin,
		SetPass:    cmdSetPass,
	}
)

const (
	SetCmd   = "/set"
	GetCmd   = "/get"
	DelCmd   = "/del"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

type userState uint

const (
	Unknown    userState = 0000
	Begin      userState = 0001
	SetService userState = 1000
	SetLogin   userState = 1001
	SetPass    userState = 1002
	Get        userState = 2000
	Del        userState = 3000
)

type userContext struct {
	text   string
	id     int64
	chatID int64
	p      *Processor
}

func (p *Processor) doCmd(text string, chatID int64, username string, userID int64) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'(id: %d)", text, username, userID)

	userState := userStates[userID]
	cmd := userCommands[userState]

	if cmd == nil {
		log.Printf("got some error with the command")
		return p.tg.SendMessage(chatID, msgCommandError)
	}

	err := cmd(context.WithValue(context.TODO(), "user", userContext{text, userID, chatID, p}))

	if err != nil {
		userStates[userID] = Begin
		p.tg.SendMessage(chatID, msgBotError)
		return err
	}

	return nil
}

func cmdBegin(ctx context.Context) error {
	userCtx := ctx.Value("user").(userContext)

	delete(userContexts, userCtx.id)

	switch userCtx.text {
	case StartCmd:
		userStates[userCtx.id] = Begin
		return userCtx.p.tg.SendMessage(userCtx.chatID, msgHello)
	case HelpCmd:
		userStates[userCtx.id] = Begin
		return userCtx.p.tg.SendMessage(userCtx.chatID, msgHelp)
	case DelCmd:
		userStates[userCtx.id] = Del
		return userCtx.p.tg.SendMessage(userCtx.chatID, msgDelBegin)
	case GetCmd:
		userStates[userCtx.id] = Get
		return userCtx.p.tg.SendMessage(userCtx.chatID, msgGetBegin)
	case SetCmd:
		userStates[userCtx.id] = SetService
		return userCtx.p.tg.SendMessage(userCtx.chatID, msgSetService)
	default:
		return userCtx.p.tg.SendMessage(userCtx.chatID, msgUnknownCommand)
	}
}

func cmdGet(ctx context.Context) error {
	userCtx := ctx.Value("user").(userContext)

	p := userCtx.p

	userService, err := p.storage.Pick(ctx, userCtx.id, userCtx.text)

	if err == storage.ErrNoSuchService {
		userStates[userCtx.id] = Begin
		return p.tg.SendMessage(userCtx.chatID, msgNoSuchService)
	}

	if err != nil {
		return err
	}

	userStates[userCtx.id] = Begin
	msg, err := p.tg.SendMessageWithResponse(userCtx.chatID, fmt.Sprintf(msgGet, userService.Service, userService.Login, userService.Password))
	p.deleter.AddMsg(msg)

	return err
}

func cmdDel(ctx context.Context) error {
	userCtx := ctx.Value("user").(userContext)

	p := userCtx.p

	err := p.storage.Remove(ctx, userCtx.id, userCtx.text)

	if err == storage.ErrNoSuchService {
		userStates[userCtx.id] = Begin
		return p.tg.SendMessage(userCtx.chatID, msgNoSuchService)
	}

	if err != nil {
		return err
	}

	userStates[userCtx.id] = Begin
	return p.tg.SendMessage(userCtx.chatID, msgDel)
}

func cmdSetService(ctx context.Context) error {
	userCtx := ctx.Value("user").(userContext)

	userContexts[userCtx.id] = &storage.UserService{UserID: userCtx.id, Service: userCtx.text}

	userStates[userCtx.id] = SetLogin
	return userCtx.p.tg.SendMessage(userCtx.chatID, msgSetLogin)
}

func cmdSetLogin(ctx context.Context) error {
	userCtx := ctx.Value("user").(userContext)

	userContexts[userCtx.id].(*storage.UserService).Login = userCtx.text

	userStates[userCtx.id] = SetPass
	return userCtx.p.tg.SendMessage(userCtx.chatID, msgSetPass)
}

func cmdSetPass(ctx context.Context) error {
	userCtx := ctx.Value("user").(userContext)

	userContexts[userCtx.id].(*storage.UserService).Password = userCtx.text
	userStates[userCtx.id] = Begin

	if err := userCtx.p.storage.Save(ctx, userContexts[userCtx.id].(*storage.UserService)); err != nil {
		return err
	}
	return userCtx.p.tg.SendMessage(userCtx.chatID, msgSaved)
}
