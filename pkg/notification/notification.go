package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type Body struct{
    Type string `json:"type"`
    Time string `json:"time"`
    Guests int `json:"guests"`
    Id string `json:"id"`
}

func NotifyRestaurant(ordertype,time string,guests int, token []string,webtoken []string,cartid string){
	config := &firebase.Config{ProjectID: os.Getenv("FIREBASE_PROJECT_NAME")}
	opt := option.WithCredentialsFile(os.Getenv("KEY_PATH"))
	app, err := firebase.NewApp(context.Background(), config, opt)

	if err != nil {
		fmt.Println(err)
	}


	ctx := context.Background()
	client, err := app.Messaging(ctx)

	if err != nil {
		fmt.Println(err)
	}

    title := "New Order Request"
	body := "Looks like somebody is waiting to have food from your restaurant."


	dataBody := map[string]string{"ClickAction":"notification-screen",}
    u,err := json.Marshal(Body{Type: ordertype, Time: time, Guests: guests, Id: cartid })
    if err != nil{
            fmt.Println(err)
    }


    messages := &messaging.MulticastMessage{
            Notification: &messaging.Notification{
                Title: title,
                Body: body,
            },
            Android: &messaging.AndroidConfig{
                Notification: &messaging.AndroidNotification{
                    Priority:5,
                },
            },
            Data:dataBody,
            Tokens: token,
    }
       // {
       //     Notification: &messaging.Notification{
       //         Title: "New Order Request",
       //         Body: string(u),
       //     },
       //     Token: token[1],
       // },
    //}

  //  if token[0] == ""{
  //      messages = []*messaging.Message{
  //          {
  //              Notification: &messaging.Notification{
  //                  Title: "New Order Request",
  //                  Body: string(u),
  //              },
  //              Token: token[1],
  //          },
  //      }   
  //  }

  //  if token[1] == ""{
  //      messages = []*messaging.Message{
  //          {
  //              Notification: &messaging.Notification{
  //                  Title: title,
  //                  Body: body,
  //              },
  //              Android: &messaging.AndroidConfig{
  //                  Notification: &messaging.AndroidNotification{
  //                      Priority:5,
  //                  },
  //              },
  //              Data:dataBody,
  //              Token: token[0],
  //          },
  //      }
  //  }


	response, err := client.SendMulticast(ctx, messages)
	if err != nil{
		fmt.Println("error sending notification:",err)
	}
	fmt.Println("Successfully sent mobile notification:", response)
    
    messages = &messaging.MulticastMessage{
            Notification: &messaging.Notification{
                Title: title,
                Body: string(u),
            },
            Tokens: webtoken,
    }

    response, err = client.SendMulticast(ctx,messages)
    if err != nil {
		fmt.Println("error sending notification:",err)
    }

    fmt.Printf("Web notification sent successfully\n", response)
}

func NotifyUser(title,body,token string){
	config := &firebase.Config{ProjectID: os.Getenv("FIREBASE_USER_PROJECT_NAME")}
	opt := option.WithCredentialsFile(os.Getenv("USER_KEY_PATH"))
	app, err := firebase.NewApp(context.Background(), config, opt)

	if err != nil {
		fmt.Println(err)
	}


	ctx := context.Background()
	client, err := app.Messaging(ctx)

	if err != nil {
		fmt.Println(err)
	}

	dataBody := map[string]string{"ClickAction":"notification-screen",}

	message := &messaging.Message{
	        Notification: &messaging.Notification{
	                Title: title,
	                Body: body,
		},
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				Priority:5,
			},
		},
		Data:dataBody,
	        Token: token,
	}

	response, err := client.Send(ctx, message)
	if err != nil{
		fmt.Println("error sending notification:",err)
	}
	fmt.Println("Successfully sent notification:", response)

}
