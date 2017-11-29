# brewery
a distributed worker master setup to run remote monitoring checks

1. Create a config file - 

```
m-c02jn0m1f1g4:parallelping skaush1$ cat /var/tmp/brewery-client.json 
{
	"broker": {
		"brokertype": "redis",
		"hostname": "localhost:6379",
		"queue": "brewery_tasks"
	},
	"inputs": [
		{
			"check_name": "check_mac_disk",
			"check_command": "/Users/skaush1/Documents/test_plugins/check_disk -c 50 -w 30 -d /"
		},
		{
			"check_name": "check_ping",
			"check_command": "/Users/skaush1/Documents/test_plugins/pping -hosts www.google.com -pingcount 2 -interval 5s"
		}
	]
}
m-c02jn0m1f1g4:parallelping skaush1$ 
```
