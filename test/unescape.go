package main

import "fmt"

import "github.com/EurasianMagpie/celadon/util"

func main()  {
	s := "Hello world"
	//s := "Help <&#8888;;> Julie to enroll in the prestigious cooking school &#39;Le Cookery&#39;."
	b, out := util.UnEscape(s)
	fmt.Println(b, out)
}