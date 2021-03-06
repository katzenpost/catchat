package main

import (
	"fmt"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/therecipe/qt/core"
)

// ConversationModel holds a collection of messages
type ConversationModel struct {
	core.QAbstractListModel

	_ func() `constructor:"init"`

	_ map[int]*core.QByteArray `property:"roles"`
	_ []*Message               `property:"messages"`

	_ func(*Message) `slot:"addMessage"`
	_ func(*Message) `slot:"appendMessage"`
	_ func(int)      `slot:"removeMessage"`
	_ func()         `slot:"clear"`
}

func (m *ConversationModel) init() {
	m.SetRoles(map[int]*core.QByteArray{
		RoleMessageID: core.NewQByteArray2("messageid", -1),
		RoleNickname:  core.NewQByteArray2("nickname", -1),
		RoleAvatar:    core.NewQByteArray2("avatar", -1),
		RoleMessage:   core.NewQByteArray2("message", -1),
		RoleStatus:    core.NewQByteArray2("status", -1),
		RoleOutbound:  core.NewQByteArray2("outbound", -1),
		RoleTimestamp: core.NewQByteArray2("timestamp", -1),
	})

	m.ConnectData(m.data)
	m.ConnectRowCount(m.rowCount)
	m.ConnectColumnCount(m.columnCount)
	m.ConnectRoleNames(m.roleNames)

	m.ConnectAddMessage(m.addMessage)
	m.ConnectAppendMessage(m.appendMessage)
	m.ConnectRemoveMessage(m.removeMessage)
	m.ConnectClear(m.clear)

	go func() {
		for {
			time.Sleep(time.Minute)
			m.updateMessageTime()
		}
	}()
}

func (m *ConversationModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if !index.IsValid() {
		return core.NewQVariant()
	}
	if index.Row() >= len(m.Messages()) {
		return core.NewQVariant()
	}

	var p = m.Messages()[index.Row()]
	if p == nil {
		return core.NewQVariant()
	}

	switch role {
	case RoleMessageID:
		{
			return core.NewQVariant1(p.MessageID)
		}
	case RoleNickname:
		{
			return core.NewQVariant1(p.Nickname)
		}
	case RoleAvatar:
		{
			return core.NewQVariant1(p.Avatar)
		}
	case RoleMessage:
		{
			return core.NewQVariant1(p.Message)
		}
	case RoleStatus:
		{
			return core.NewQVariant1(p.Status)
		}
	case RoleOutbound:
		{
			return core.NewQVariant1(p.Outbound)
		}
	case RoleTimestamp:
		{
			if time.Since(p.Timestamp) < time.Minute {
				return core.NewQVariant1("now")
			}
			return core.NewQVariant1(humanize.Time(p.Timestamp))
		}

	default:
		{
			return core.NewQVariant()
		}
	}
}

func (m *ConversationModel) rowCount(parent *core.QModelIndex) int {
	return len(m.Messages())
}

func (m *ConversationModel) columnCount(parent *core.QModelIndex) int {
	return 1
}

func (m *ConversationModel) roleNames() map[int]*core.QByteArray {
	return m.Roles()
}

func (m *ConversationModel) clear() {
	m.BeginResetModel()
	m.SetMessages([]*Message{})
	m.EndResetModel()
}

func (m *ConversationModel) addMessage(p *Message) {
	m.BeginInsertRows(core.NewQModelIndex(), 0, 0)
	m.SetMessages(append([]*Message{p}, m.Messages()...))
	m.EndInsertRows()
}

func (m *ConversationModel) appendMessage(p *Message) {
	m.BeginInsertRows(core.NewQModelIndex(), len(m.Messages()), len(m.Messages()))
	m.SetMessages(append(m.Messages(), p))
	m.EndInsertRows()
}

func (m *ConversationModel) removeMessage(row int) {
	m.BeginRemoveRows(core.NewQModelIndex(), row, row)
	m.SetMessages(append(m.Messages()[:row], m.Messages()[row+1:]...))
	m.EndRemoveRows()
}

func (m *ConversationModel) updateMessageStatus(id string, status int) {
	// fmt.Println("Updating message status:", id, status)
	for _, v := range m.Messages() {
		if v.MessageID == id {
			v.Status = status

			/*
				var index = m.Index(i, 0, core.NewQModelIndex())
				m.DataChanged(index, index, []int{RoleStatus})
			*/

			var fIndex = m.Index(0, 0, core.NewQModelIndex())
			var lIndex = m.Index(len(m.Messages())-1, 0, core.NewQModelIndex())
			m.DataChanged(fIndex, lIndex, []int{RoleStatus})

			return
		}
	}
}

func (m *ConversationModel) updateMessageTime() {
	fmt.Println("Updating message timestamps...")
	if len(m.Messages()) > 0 {
		var fIndex = m.Index(0, 0, core.NewQModelIndex())
		var lIndex = m.Index(len(m.Messages())-1, 0, core.NewQModelIndex())
		m.DataChanged(fIndex, lIndex, []int{RoleTimestamp})
	}
}

func init() {
	ConversationModel_QRegisterMetaType()
	Message_QRegisterMetaType()
}
