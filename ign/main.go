package ign

// package main

// import (
// 	"fmt"
// 	"net/url"
// 	"strings"
// )

// func main() {
// 	params := url.Values{
// 		"lat":       {"3.2;3.9"},
// 		"lon":       {"3.2;3.9"},
// 		"resource":  {"test"},
// 		"delimiter": {";"},
// 		"measures":  {"false"},
// 		"zonly":     {"true"},
// 	}

// 	fmt.Printf("%+v\n", params)
// 	fmt.Println("=====")
// 	fmt.Println()

// 	var out string
// 	for k := range params {
// 		out = out + "&" + k + "=" + strings.Join(params[k], "")
// 		fmt.Printf("%+v\n", k)
// 	}

// 	fmt.Println("out", out)
// }
