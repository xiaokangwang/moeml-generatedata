package main

import (
	"flag"
	"github.com/xiaokangwang/moeml-generatedata/challangegen"
)

func main() {
	var InputF,InputB,Out string
	var sum,n int
	flag.StringVar(&InputF,"f","F","")
	flag.StringVar(&InputB,"b","B","")
	flag.StringVar(&Out,"o","O","")
	flag.IntVar(&sum,"s",3,"")
	flag.IntVar(&n,"n",80,"")
	flag.Parse()
	o:=challangegen.Generator{InputB,InputF,sum,Out}
	for i:= 0;i<n;i++ {
		o.ComposeAll(i)
	}
}
