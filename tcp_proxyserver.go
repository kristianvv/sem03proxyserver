package main

import (
    "log"
    "net"
    "io"
)

func main() {
    proxyServer, err := net.Listen("tcp", "172.17.0.3:8081")
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("bound to %s", proxyServer.Addr().String())

    for {
        conn, err := proxyServer.Accept()
        if err != nil {
            log.Println(err)
            continue
        }
        go func(client net.Conn) {
            defer client.Close()

            server, err := net.Dial("tcp", "172.17.0.4:8080")
            if err != nil {
                log.Println(err)
                return
            }
            defer server.Close()

            err = proxy(client, server)
            if err != nil && err != io.EOF {
                log.Println(err)
            }
        }(conn)
    }
}

func proxy(client, server net.Conn) error {
    go func() {
        defer server.Close()
        defer client.Close()

        _, err := io.Copy(server, client)
        if err != nil {
            log.Println(err)
        }
    }()

    _, err := io.Copy(client, server)
    if err != nil {
        log.Println(err)
    }

    return err
}
