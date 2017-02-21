
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

Build it:

```bash
go build .
```

Run the server:

```bash
./go-for-the-galaxy -server
```

Connect clients to the server:

```bash
./go-for-the-galaxy -client -addr $SERVER_ADDR:7771
```

If the server is running locally, simply do:

```bash
./go-for-the-galaxy -client
```
