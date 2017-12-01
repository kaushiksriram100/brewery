# brewery
a distributed worker setup to run collect metrics of remote clients and send to graphite. Ideal to collect remote check metrics (like check_http,ping metrics or some app's metrics by polling their apis or like resulting count of a search query to elasticsearch etc). Our workers are distributed so metrics will be collected in HA mode even if some nodes go down. 

Dependency: 

We use Redis for queue, so please install redis env(https://redis.io/topics/quickstart)


1. Create a config file with all the checks that return metrics. The check output must conform to the graphite standards.

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
                        "check_command": "/Users/skaush1/Documents/test_plugins/check_ping_metrics -u www.google.com -p 443"
                },
                {
                        "check_name": "check_ping",
                        "check_command": "/Users/skaush1/Documents/test_plugins/test.sh"
                }
        ],
        "outputs": {
                "graphite": {
                        "endpoint": "dfw-mon12.prod.walmart.com",
                        "port": 2003
                }

        }
}
```

plugin outputs must be like this else it will get dropped. 

```

```

2. Provide the redis server endpoints and port as mentioned above. 

3. Start the dispatcher like this. (use -w to set the max number of pools to allow while processing the results. We will not spawn more than those many goroutines). -w 0 means results are processed serially. Recommended to set -w > 0. Runs on one master machine. 
```
	go build dispatcher.go
	./dispatcher -c /var/tmp/brewery-client.json -w 300
```

4. Build the worker code (brewery-client.go)

```
	go build -o metric_worker
	
```

5. Start workers (on multiple nodes). Start this on several nodes. These will consume from redis. 
```
./metric_worker -c /var/tmp/brewery-client.json

```



NOTE: please teach me if something is not right in the code. I am new to this. PRs welcome. Need to add logging features. I will do that. 
