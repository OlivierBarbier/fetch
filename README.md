# fetch

Pour mettre un timeout sur golang

timeout := time.Duration(5 * time.Second)
client := http.Client{
    Timeout: timeout,
}
client.Get(url)

http://stackoverflow.com/questions/16895294/how-to-set-timeout-for-http-get-requests-in-golang
