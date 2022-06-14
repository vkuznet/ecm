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

#### Helpful hints
Here we provide helpful hints about ECM UI setup and builds. To build ECM
for mobile please follow the following steps:
- build appropriate application package for your mobile
  - for Andoird please use `make build_android` while for iOS `make build_iphone`
- run Android Studio
  - run emulator and check devices `./adb devices`
  - deploy ECM applocation `./adb install ecm.apk`
  - check application log `./adb logcat | grep fyne`
  - on macOS we write sync config to `$HOME/Library/Preferences/fyne/io.github.vkuznet`
- on iOS we need to run Xcode
  - obtain developer [certificate](https://help.apple.com/xcode/mac/current/#/dev154b28f09?sub=dev6dab365c2)
    - create new developer account in Xcode using AppleID
    - click on Manage certificates, Control click on your AppleID and export
      certificate
    - import your developer certificate into keychain
    - finally, in keychain click on your Apple Development certificate and grab
    its ID, e.g. WXYZ123XYZ
