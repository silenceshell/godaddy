A small tool for DDNS(Dynamic DNS), to access linux server in lan from Internet via dns name.

It will periodly get external IP on Internet and call Godaddy API to set A record for your domain name. You should have an domain in Godaddy and generate an key/secret for developing.

# Godaddy config

Ref: [Get Started](https://developer.godaddy.com/getstarted)

# Compile

This is optional because you could use my docker image on Dockerhub: `silenceshell/godaddy:0.0.1`.

```
GOOS=linux go build godaddy.go
cp godaddy artifacts
pushd artifacts
docker build . -t silenceshell/godaddy:0.0.2
popd
rm godaddy
```

# Run on kubernetes

```
kubectl run godaddy --image=silenceshell/godaddy:0.0.2 --command --/godaddy ${godaddy key} ${godaddy secret}
```
