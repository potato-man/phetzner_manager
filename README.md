# Potato hetzner manager

A simple telegram bot to manage your pool of servers.  
At the moment, server power management is in place.

-----


##### Requirments:
[golang](https://go.dev/dl/)  
[telebot](https://github.com/tucnak/telebot)  
[hcloud](https://pkg.go.dev/github.com/hetznercloud/hcloud-go/v2/hcloud)


-----

##### setup build run:
First, you will need to add the environment variables for your Telegram token (TELEGRAM_TOKEN) and Hetzner API (HCLOUD_TOKEN).  
You will also need to specify the ID of your Telegram account in the code (Allowed UserID).  
This is necessary to restrict access to the bot to authorized users only.  
Once you have entered all the necessary information, you can start building the program. 
```
go build phetzbot.go
```
