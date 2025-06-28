# PO System APIs

An experiment wherein I add a Go API to speed up the list POs page.

## Developing the app

This app is run quite simply, just do:

```
go run ./go-app/main.go
```

However, you will need to do some setup first. First, make sure you are logged into the [gcloud cli](https://cloud.google.com/sdk/docs/install):

```
gcloud auth login
```

Then, set your project

```
gcloud config set project cdac(-demo)-purchaseorder
```

Finally, create the app default credentials:

```
gcloud auth application-default login
```

After these steps are complete, run `go run ./go-app/main.go` and try calling a few APIs to see if things are connected.

## Deployment

The app is deployed with scripts in the [purchase-order-system](https://github.com/cdachurch/purchase-order-system) repository using npm run scripts. Check out the package.json for the scripts, and make sure you have this repo and that repo checked out in the same root directory for the best chance at success!
