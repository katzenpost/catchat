package main

import (
	"fmt"

	"github.com/therecipe/qt/core"
)

// AccountBridge makes an account available in QML
type AccountBridge struct {
	core.QObject

	_ string `property:"nickname"`
	_ string `property:"avatar"`
	_ string `property:"status"`
	_ string `property:"error"`
	_ string `property:"recipient"`
	_ bool   `property:"keyExchanged"`

	_ func(passphrase string, nickname string) bool `slot:"addContact"`
	_ func(contact string) bool                     `slot:"removeContact"`
	_ func(contact string)                          `slot:"loadConversation"`
	_ func(recipient string, message string)        `slot:"sendMessage"`
	_ func(contact string, url string)              `slot:"loadAvatar"`
	_ func(passphrase string)                       `slot:"loadCatshadow"`

	_ *core.QAbstractListModel `property:"contactListModel"`
	_ *core.QAbstractListModel `property:"conversationModel"`
}

// ConfigBridge allows QML to access the app's config
type ConfigBridge struct {
	core.QObject

	_ bool   `property:"firstRun"`
	_ string `property:"authURL"`
	_ string `property:"redirectURL"`
	_ string `property:"theme"`
	_ string `property:"style"`
	_ string `property:"notification"`
	_ int    `property:"positionX"`
	_ int    `property:"positionY"`
	_ int    `property:"width"`
	_ int    `property:"height"`
}

var (
	accountBridge *AccountBridge
	configBridge  *ConfigBridge
)

// setupQmlBridges initializes the QML bridges
func setupQmlBridges() {
	accountBridge = NewAccountBridge(nil)
	configBridge = NewConfigBridge(nil)

	accountBridge.ConnectAddContact(func(passphrase string, nickname string) bool {
		return addContact(catshadowClient, nickname, passphrase)
	})
	accountBridge.ConnectRemoveContact(func(nickname string) bool {
		return removeContact(catshadowClient, nickname)
	})
	accountBridge.ConnectLoadConversation(func(nickname string) {
		loadConversation(catshadowClient, nickname)
	})
	accountBridge.ConnectSendMessage(func(recipient string, message string) {
		sendMessage(catshadowClient, recipient, message)
	})
	accountBridge.ConnectLoadAvatar(func(nickname string, iu string) {
		fmt.Println("Loading avatar:", iu)
		contactListModel.updateAvatar(nickname, iu)
	})
	accountBridge.ConnectLoadCatshadow(func(password string) {
		go setupCatShadow(catshadowCfg, []byte(password))
	})
}
