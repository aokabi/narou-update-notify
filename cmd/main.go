package main

import (
	"fmt"

	"github.com/aokabi/narou-update-notify/api"
)

func main() {
	response, err := api.GetNovelInfo()
	if err != nil {
		panic(err)
	}
	// 0番目はよくわからんサマリーが入っている
	fmt.Println(response[1].Title)
	fmt.Println(response[1].GeneralLastup)

}
