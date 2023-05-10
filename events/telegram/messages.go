package telegram

const msgHelp = `👩‍🏫 I can keep login and password of your services
					
					/set - set new service data

					/get - get service data

					/del - delete service data
`

const msgHello = "👋 Hi there!\n\n" + msgHelp

const (
	msgCommandError   = "🤕 There is some error with the command"
	msgUnknownCommand = "🤔 Unknown command"
	msgNoSuchService  = "🙊 You don't have such service"
	msgSaved          = "💾 Saved!"
	msgGetBegin       = "✏️ Enter a service"
	msgDelBegin       = "✏️ Enter a service"
	msgGet            = "🎉 Here is your info about %s\n\n login: %s\n password: %s"
	msgDel            = "💥 Deleted!"
	msgSetService     = "1️⃣ Enter a service"
	msgSetLogin       = "2️⃣ Enter a login"
	msgSetPass        = "3️⃣ Enter a password"
	msgBotError       = "🤖 I've got some error. Reset..."
)
