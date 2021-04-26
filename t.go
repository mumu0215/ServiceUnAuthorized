package main

import (
	"fmt"
	"strings"
)

func main(){
	a:=[]string{
		"213",
		" sdfsd ",
		" dfsfff",
	}
	fmt.Println(a)
	b:= func() []string{
		var temp []string
		for _,n :=range a{
			temp=append(temp,strings.TrimSpace(n))
		}
		return temp
	}()
	fmt.Println(b)
}
