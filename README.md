Whspr
=====

_Shhhh it's a secret._

### Background

Whspr is a secure messaging protocol. More of a just for fun project, it could be used when secure information needs to be transfered.

I'm sure there are plenty of great solutions out there, but I wanted someting that would alow me to feel like a 1337h4x0r.

## To use

1. Generate a x509 certficate (RSA, 2048bit, extensions don't matter), with CommonName = your screen name.

1.) Distribute your public certificate to others you'd like to communicate with

1.) Agree on a pre shared key

1.) Agree on a host machine and port

1.) Relevent configs go into config.go file.

```
git clone https://github.com/c0nrad/whspr.git
cd whspr
go run client.go 
```

## Infastructure

To use the service, all you need is an echo server. The server accepts data on a TCP connection, and replays it to all other open connections.

In my model I assume the server can't be trusted.

Everyone that uses the service is a client.

## Security

To use the service, you need to be in a party that has the same PSK, and everyone in that party needs to have the public keys of each other member. 

Messages and transfered over the wire in json objects like such:

```javascript
{
  IV: "eTJ4UzREVVFjblZ4b2NvTWl5MFE4dz09",

  Data: "K2N5U3FFRFR5cnk1STdVZW1tWnVzZzliQ3VGMU1rQ3RxeDFlQjVEYUJpOD0=",

  Signature: "Vm5ycFFyanN2T1JtcDJQWWxXV2lLbHNGWVprL2M5NUQ1ZHMwUU9uclZTa1Bqc1BCcW1kdlRqc3pPdW9UdXNsbDFqRTlSZEo2ejQrNFUvMGRra2krajZkTG5yUzk1NVJ6U0M3dm5DRDhveXVPb3oxSzZWVi9kY1RsejJ4UnFMakQ0bGwrUnQrdnJJSVRWQ3RxVktpeEJib0JrREpFcDA3Mmx0RlRJYnhxaHNEM0Z0Nm5yYUdZK3dsMit6Z2xmbWhDUEx2UU93K2FqSzN5allUckFmamNJYmtrTGJ0ejJoS0JLMzNZaU1seU02YWdETTcxNFFjOCtVVTUrM3R0Zi9iOURFRnJhejR1SFRpTi8zektsYmtoKzZINjJPRDE0cWpqdzlWY1pkd2F6NVV6dUJaK2VoU0ZwOTVUdVlDbzBFa2R3WFVtUEpEWnhsZHRtT3Q4cjJLMGpnPT0=",

  Name:"WXpCdWNtRms="
}

### IV 

base64(IV)

The AES IV is rolling. I use the last 16 bytes of each message as the next message. It's also initialized to a random set of bytes at startup.

### Data

AES256-CBC("hello world", secret_preshared_key, IV)

### Signature

PKCS1v1.5 of SHA256 using RSA private key. 

### Name

TO BE REMOVED SOON. Instead just check signature for each party you trust. Name will NOT be transfered over wire. Just helpful for debugging.
