package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"runtime"
	"src/utill"
)


func main() {
	runtime.GOMAXPROCS(utill.NumOfCpu)
	app:=&cli.App{
		Name:                   "unAuthorizedCheck",
		Description:"Tools to scan unAuthorized services",
		UseShortOptionHandling:true,
		Flags: []cli.Flag{
			&cli.StringFlag{Name:"input",Value:"",Required:true,Aliases:[]string{"i"},Destination:&utill.InputFile,Usage:"input json `file` to scan",},
			&cli.BoolFlag{Name:"brute",Value:false,Aliases:[]string{"b"},Destination:&utill.IsBrute,Usage:"run password brute at last"},
			&cli.StringFlag{Name:"password",Value:"",Aliases:[]string{"p"},Destination:&utill.PassWordFile,Usage:"password `file` for brute(-b needed)"},
			&cli.IntFlag{Name:"thread",Value:5,Aliases:[]string{"t"},Destination:&utill.Thread,Usage:"brute process thread(-b needed)"},
		},
		Action:utill.FlagMain,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}