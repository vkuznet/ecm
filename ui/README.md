### ECM UI interface
This area contains code for ECM UI based on [fyne.io](https://fyne.io/)

#### Build
To build please do
```
# generic build
make

# specific architecutre build (see Makefile), e.g.
make build_linux

# android package
make build_android

# ios package
make build_ios
```
Please note, to properly build mobile package you'll need to obtain cloud provider
credentials for you app and put them into `credentials.env`. 
