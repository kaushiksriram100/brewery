package main 


import (
"fmt"
machinery "github.com/RichardKnop/machinery/v1"
machineryconfig "github.com/RichardKnop/machinery/v1/config"
breweryutils "github.com/kaushiksriram100/brewery/brewery-utils"
tasks "github.com/RichardKnop/machinery/v1/tasks"
"time"
"flag"
"github.com/marpaia/graphite-golang"
)

func SendToKafka(message interface {}, brokers, topic string ) {


	
}


func main() {

	configfile := flag.String("c", "/var/tmp/brewery-conf.json", "- json format config file with broker end points")
	concurrency := flag.Int("w", 3, "- dispatcher workers")
	flag.Parse()
	brewerycnf, err := breweryutils.LoadServerConfig(*configfile)  //cnf is of type *config.Config
	if err != nil {
		fmt.Printf("Some issues loading the config- %s", err)
		return
	}

	//Now broker endpoints are unmarshalled and put in the Config Struct. Let's add this to the brokerconfig struct we created and return
	var brokerconfig = &machineryconfig.Config{} //this config as mandated in machinery code. Don't confuse with the above config that is meant to parse the json config file

	(*brokerconfig).Broker = brewerycnf.Broker.MessageBroker+"://"+brewerycnf.Broker.BrokerHost
	(*brokerconfig).ResultBackend = brewerycnf.Broker.MessageBroker+"://"+brewerycnf.Broker.BrokerHost
	(*brokerconfig).DefaultQueue = brewerycnf.Broker.DefaultQueue

	//create a new server object. Even for a client, machinery requires a server to be instantiated before launching the worker. We will use our custom utils to do all that. 

	server, err := machinery.NewServer(brokerconfig)
	if err != nil {
		fmt.Printf("Can't Create Server.. exiting")
		return
	}

	//create a slice of *tasks.Signature

	alltasks := make([]*tasks.Signature,0,2048)



//Send each check element in the inputs array to the queue

	for _,v := range brewerycnf.Inputs {

		alltasks = append(alltasks, &tasks.Signature{
		Name: "CommandExecutor",
		Args: []tasks.Arg{tasks.Arg{Type: "string",Value: v.CheckCommand}},
		})

	}



	//create a dispatcher worker pool before sending metrics to graphite. we can use an empty struct with no values just as a placeholder

	wrkrpool := make(chan struct{}, *concurrency)

  go func() {

  		if (*concurrency) > 0 {

		for i := 0; i < *concurrency; i++ {
			wrkrpool <- struct{}{}  //each empty struct value represents a pool. this is the power of empty struct
		}

	} else {
		wrkrpool <- struct{}{}  //if concurrency is 0 then we will put just one value for now as belw we need to consume from pool to send metrics to graphite
	}
  }()


// do a for loop to send tasks every n minutes.. right now hardcoding sleep time as 60 sec


for {

//alltasks has all the tasks. We can create a group of tasks. 
	group := tasks.NewGroup(alltasks...)
	asyncResults, err := server.SendGroup(group,1000)

	if err != nil {
		fmt.Println("failed to send tasks")
		return
	}

//Receive the outputs of all the checks and start processing the output.. PENDING

	//Create a graphite object and send it over. 

	Graphite, err := graphite.NewGraphite(brewerycnf.Outputs.GraphiteOutput.Graphiteendpoints, brewerycnf.Outputs.GraphiteOutput.Graphiteport)

	if err != nil {
		fmt.Println("sorry. could not connect to graphite. Can't process metrics for now")
		return
	}

	for _, asyncResult := range asyncResults {
  		results, err := asyncResult.Get(time.Duration(time.Millisecond * 10))
		  if err != nil {
			   fmt.Println("failed to process results for-:", asyncResult.Signature.Name, asyncResult.Signature.UUID)
  				 continue
 		  }

 		  <-wrkrpool  //we will take one pool out. This will block if no pools are available until one is available(graphitehandler.go will return a value after done). So that we don't end up exhausting goroutines
 		  go breweryutils.SendToGraphite(results,Graphite, wrkrpool)
  		
	}

	duration := time.Second * 60
  	time.Sleep(duration)

}


}