package utilites

import "time"
import "color"


//  LogStamp() returns current time in 00:00:00
//  in format in format for logging (yellow and \t at the end)
func LogStamp() string {
    t := time.Now()
    return color.Yellow(t.Format("15:04:05")+"\t")
}

func LogStampErr() string {
    t := time.Now()
    return color.Red(t.Format("15:04:05")+"\t")
}