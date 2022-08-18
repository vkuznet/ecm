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
    - how to codesign with provisioning profile in xcode see
    [here](https://steemit.com/xcode/@ktsteemit/xcode-free-provisioning)
    and [here](https://ioscodesigning.com/generating-code-signing-files/)
    and
    [here](https://www.testdevlab.com/blog/2019/07/24/xcode-provisioning-profile-automation-for-ci/)
       - attach existing Apple device to computer
       - start xcode and create new project
       - choose to use your device for that project
       - click on top name of your project which will lead you to its settings
       - in settings you'll find that you may use your development AppleID team
         to codesign the project

#### A note about maxOS build
The process involves the following steps:
```
fyne release --os darwin --appID io.github.vkuznet.ECM --appVersion 1.0 --appBuild 1 --cert W123XYZ123XYZ --category utilities --profile "iOS Team Provisioning Profile: io.github.vkuznet.ECM"
### /usr/bin/codesign -vfs W123XYZ123XYZ --entitlement entitlements.plist ecm.app
ecm.app: signed app bundle with Mach-O thin (x86_64) [io.github.vkuznet.ECM]
### /usr/bin/productbuild --component ecm.app /Applications/ --product ecm.app/Contents/Info.plist ecm-unsigned.pkg
### /usr/bin/productsign --sign W123XYZ123XYZ ecm-unsigned.pkg ecm.pkg
```
The last command here should use CA identity which can be found from 
`security -v find-identity` command.
