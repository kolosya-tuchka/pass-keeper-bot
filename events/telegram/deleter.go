package telegram

import (
	"github.com/hashicorp/go-set"
	"log"
	"pass-keeper-bot/clients/telegram"
	"sync"
	"time"
)

var cmp = func(msg1, msg2 MsgDel) int {
	return msg1.sendTime.Compare(msg2.sendTime)
}

type MsgDelSchedule struct {
	tg     *telegram.Client
	msgSet set.TreeSet[MsgDel, set.Compare[MsgDel]]
	sync.Mutex
}

func NewSchedule(tg *telegram.Client) *MsgDelSchedule {
	return &MsgDelSchedule{tg: tg, msgSet: *set.NewTreeSet[MsgDel, set.Compare[MsgDel]](cmp)}
}

func (mds *MsgDelSchedule) AddMsg(msg *telegram.BotMessage) {
	log.Printf("new message to delete: chat %d message %d", msg.Chat.ID, msg.ID)

	mds.Lock()
	mds.msgSet.Insert(MsgDel{msg: msg, sendTime: time.Now().Add(time.Minute * 5)})
	mds.Unlock()
}

func (mds *MsgDelSchedule) Init() {
	go func() {
		for {
			delTime, slice := time.Now(), mds.msgSet.Slice()

			for i := 0; i < len(slice); i++ {
				if slice[i].sendTime.Compare(delTime) < 0 {
					mds.Lock()
					mds.msgSet.Remove(slice[i])
					mds.Unlock()
					err := mds.tg.DeleteMessage(slice[i].msg.Chat.ID, slice[i].msg.ID)
					if err != nil {
						log.Printf("can't delete message: %s", err.Error())
						continue
					}
					log.Printf("message was deleted: chat %d message %d", slice[i].msg.Chat.ID, slice[i].msg.ID)
				} else {
					break
				}
			}

			time.Sleep(time.Second * 1)
		}
	}()
}

type MsgDel struct {
	msg      *telegram.BotMessage
	sendTime time.Time
}
