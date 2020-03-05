//
// Copyright (C) 2020 white <white@Whites-Mac-Air.local>
//
// Distributed under terms of the MIT license.
//

package main

import (
    "flag"
    "encoding/base64"
    "encoding/json"
    "time"
    "errors"
    //"log"
    //"net/http"
    "os"
    //"reflect"
    "fmt"
    "crypto/tls"

    "github.com/Nerzal/gocloak/v4"
    "github.com/go-resty/resty/v2"
)

var (
    action string
    username string
    password string
    group string
)

type PostResult struct{
    status int
    msg string
}

func getUserToken(user string, pwd string) (string, error){
    //realm := os.Getenv("REALM_NAME")
    serverURL := os.Getenv("SERVER_URL")
    //clientID := os.Getenv("CLIENT_ID")
    //clientSecretEncoded := os.Getenv("CLIENT_SECRET_KEY")

    /*
    clientSecret, err := base64.StdEncoding.DecodeString(clientSecretEncoded)
    if err != nil{
        return nil, errors.New("decode CLIENT_SECRET_KEY  error")
    }
    */
    client := gocloak.NewClient(serverURL)
    restyClient := client.RestyClient()
    //restyClient.SetDebug(true)
    restyClient.SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: true })

    /*
    jwt, err := client.Login(clientID, string(clientSecret), realm, user, pwd)
    fmt.Println(jwt)
    fmt.Println(err.Error())
    jwt, err = client.LoginClient(clientID, string(clientSecret), realm)
    fmt.Println(jwt)
    fmt.Println(err.Error())
    jwt, err := client.LoginAdmin(user, pwd, realm)
    fmt.Println(jwt)
    fmt.Println(err)
    jwt.access_token
    */
    return "fake token", nil
}
// 2. create user

func createMinioUser(token string, username string, password string, group string) error{
    apiURL := os.Getenv("API_URL")

    data := fmt.Sprintf(`{"username": "%s", "password": "%s", "group": "%s"}`,
                        username, password, group)

    client := resty.New()
    resp, err := client.
                SetTimeout(3 * time.Second).
                R().
                EnableTrace().
                SetHeader("Content-Type", "application/json").
                SetHeader("Authorization", token).
                SetBody(data).
                Post(apiURL)
    if err == nil && resp.StatusCode() == 200{
        var result PostResult
        err := json.Unmarshal(resp.Body(), &result)
        if err != nil{
            return errors.New("failed to parse response")
        }else{
            if result.status == 0{
                return nil
            }else{
                return errors.New("failed to create minio user")
            }
        }
    }else{
        return errors.New("create minio user request error")
    }

}


func usage() {
    fmt.Fprintf(os.Stderr, `manage minio user
Usage: minio_user [-h] [-a] [-u user] [-p password] [-g group]

Options:
`)
    flag.PrintDefaults()
}

func init(){
    flag.StringVar(&action, "a", "", "add/update/delete minio user")
    flag.StringVar(&username, "u", "", "user name")
    flag.StringVar(&password, "p", "", "password")
    flag.StringVar(&group, "g", "", "user group")
    flag.Usage = usage
}
func main() {
    flag.Parse()
    cuser := os.Getenv("ADMIN_USER")
    encoded := os.Getenv("USER_PASSD")
    decoded, err := base64.StdEncoding.DecodeString(encoded)
    if err != nil{
        fmt.Println("decode USER_PASSD error")
        os.Exit(1)
    }
    token, err := getUserToken(cuser, string(decoded))
    if err != nil{
        fmt.Println("failed to fetch token")
        os.Exit(2)
    }
    //fmt.Println(token)
    if action == "add"{
        err = createMinioUser(token, username, password, group)
        if err != nil{
            fmt.Println(err.Error())
            os.Exit(3)
        }else{
            fmt.Println("success to create minio user")
        }

    }
}
