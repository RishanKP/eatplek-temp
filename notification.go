package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

func notification(){

    type Body struct{
        Type string `json:"type"`
        Time string `json:"time"`
        Guests int `json:"guests"`
    }

	config := &firebase.Config{ProjectID: os.Getenv("FIREBASE_PROJECT_NAME")}
	opt := option.WithCredentialsFile(os.Getenv("KEY_PATH"))
	app, err := firebase.NewApp(context.Background(), config, opt)

	if err != nil {
		fmt.Println(err)
	}


    //data := map[string]interface{}{"time":"5.00 pm","guests":10}

    u,err := json.Marshal(Body{Type:"Take Away", Time:"5:00pm",Guests: 2})

    if err != nil{
        fmt.Println(err)
    }

    ctx := context.Background()
	client, err := app.Messaging(ctx)

	if err != nil {
		fmt.Println(err)
	}

	message := &messaging.Message{
	    Notification: &messaging.Notification{
	                Title: "New Order Request",
	                Body: string(u),
		},
	//	Webpush: &messaging.WebpushConfig{
	//		Notification: &messaging.WebpushNotification{
	//			CustomData: data,
	//		},
	//	},
	    Token: "e0BLTB1tTEmTWekAbCuKhX:APA91bHGD10Rm5f16BFtJlYcg8Nat-3AEdsfrVKOUljkJGhUowaiylOTvgC8d6WaGBQ9QYfLXdyx85fcyV54T8rny9NVtSoq5-zD86sXPuAlc6mSg-tqQheZSmr_WYFlQswzHEPOTJUq",
	}

	response, err := client.Send(ctx, message)
	if err != nil{
		fmt.Println("error sending notification:",err)
	}
	fmt.Println("Successfully sent notification:", response)

}
