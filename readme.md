# Telegram news bot
## What does this bot do
The bot collects news from sources using RSS, and the field sends news to the telegram channel.
## Bot commands
**Basic:**  
_/help_ and _/start_ for help information.  
**Admin commands:**  
(only admins of channel can use it)  
/add NAME, URL - add new source  
/list - get all sources  
/delete ID - delete source   
## How to start bot
- Create new bot in BotFather
- In BotFather fro your bot set "Group privacy" - OFF 
- Create new channel
- Add bot to channel
- Apply migartions (internal/storage/migrations) using [goose](https://github.com/pressly/goose) or using other method
- Create config.hcl or config.local.hcl
- Fill config like this:
```
telegram_bot_token = "token from botfather"
telegram_channel_id = channel id 
database_dsn = "postgres://..."
fetch_interval = "5m" # how often bot will collect information
notification_interval = "30m" # how often bot will send news to channel   
```
- Build cmd/main.go
- Start main.exe
- Try your bot!
## TODO:
- Dynamic source priority
- Summary usin neural networks
- New article type (video, audio)
- More types of resources, not only RSS
