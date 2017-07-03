
```
          ___         __              _   _              ___      _                         
__/\__   / _ \___    / _| ___  _ __  | |_| |__   ___    / _ \__ _| | __ ___  ___   _  __/\__
\    /  / /_\/ _ \  | |_ / _ \| '__| | __| '_ \ / _ \  / /_\/ _` | |/ _` \ \/ / | | | \    /
/_  _\ / /_\\ (_) | |  _| (_) | |    | |_| | | |  __/ / /_\\ (_| | | (_| |>  <| |_| | /_  _\
  \/   \____/\___/  |_|  \___/|_|     \__|_| |_|\___| \____/\__,_|_|\__,_/_/\_\\__, |   \/  
                                                                               |___/        
```

An abstract multiplayer RTS inspired on Galcon written in Golang.


How to run
----------

Build the server and the client:

```bash
go build ./cmd/server
go build ./cmd/client
```

Run the server:

```bash
./client
```

Connect clients to the server:

```bash
./client -addr IPADDR:PORT
```

If the server is running locally, simply do:

```bash
./client
```
