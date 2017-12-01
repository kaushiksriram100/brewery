package utils

import (
"fmt"
"runtime"
"github.com/marpaia/graphite-golang"
"strings"
"strconv"
"reflect"
)


func SendToGraphite (value []reflect.Value,Graphite *graphite.Graphite, wrkrpool chan struct{}) {

	var metrics graphite.Metric

	var tmp interface{}

	//we get only one return value (other than err) from the CommandExecutor task. So instead of iterating through the slice, we can take the 0th element. 

	tmp = value[0].Interface() //convert reflect.Value to interface

 	result, ok := tmp.(string)

 	if ok != true {
 		fmt.Println("error")
 		wrkrpool <- struct{}{}  //give back the pool so that other tasks can be executed in dispatcher.go
 		runtime.Goexit()
 	}

 	slice1 := strings.Split(strings.Trim(result, "\n"), "\n")  //break the full result into chunks of slices. trim leading and trailing new lines

 	for _,v := range slice1 {
 		breakmetrics := strings.Split(strings.TrimSpace(v)," ") //for each graphite metric now split into chunks of slices and assign to our graphite function
 		
 		if len(breakmetrics) != 3 { //check if the metrics align to graphite format, if not drop it. 
 			wrkrpool <- struct{}{}  //give back the pool so that other tasks can be executed in dispatcher.go
 			runtime.Goexit()
 		}
 		metrics.Name = breakmetrics[0]
 		metrics.Value = breakmetrics[1]

 		time_tmp, err := strconv.ParseInt(breakmetrics[2],10,64)
 		if err != nil {
 			fmt.Println("can't get correct timestamp")
 			continue
 		}
 		metrics.Timestamp = time_tmp

 		err = Graphite.SendMetric(metrics)

 		if err != nil {
 			fmt.Println("some error sending metrics")
 			continue
 		}
 	}

 	wrkrpool <- struct{}{}
 	runtime.Goexit()

}