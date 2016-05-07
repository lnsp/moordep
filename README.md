# Moordep

A webhook listener in Go.

# Getting started
Before you download the software, be sure that you have

- a webhook service, that sends POST requests (like Docker Hub)
- a valid and signed SSL certificate
- command line access to the server where you host your stuff

If you are ready, go ahead and run
```bash
$ go get github.com/lnsp/moordep
```
**Success!** You just installed **moordep**.

# Configuration
I will provide you with the following configuration file with the path `$CFGPATH`.
It should look something like this.
```json
{
    "host": "",
    "port": 8080,
    "token": "your-private-token",
    "hooks": {
        "deploy": "./deploy-my-server.sh"
    }
}
```
Change it for your needs and save it.

# Starting up
To simplify running the server, I will use
`$CERTPATH` for the certificate file and `$KEYPATH` for the key file.

Run the following command to start up the server:
```bash
$ $GOPATH/bin/moordep -cert $CERTPATH -key $KEYPATH -config $CFGPATH
```

If you do not get an error, everything should be fine. I recommend running the job detached from your terminal using
commands like `disown`.

# Testing
You can use **cURL** to test your moordep installation.

```bash
$ curl -X POST "https://[YOUR_HOSTNAME]:[PORT]/?hook=[YOUR_HOOK]&token=[YOUR_TOKEN]"
```

Replace the strings inside the brackets to match your configuration.