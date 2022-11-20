### IOS build instructions
On iOS we need to obtain Apple Developer ID (which cost $99/year) in order
to have certificate and ability to upload the app to Apple Store.

The certificate can be obtained in a different way:
- you may start Xcode and create temporary personal account which will
contain a valid certificate for 7 days
  - this certificate will be imported into keychain and you can inspect it
- but you may obtain certificate from any valid authority and export it to
  keychain
  - in keychain, locate you certificate and grant access to Xcode and codesign
    app
  - once you obtain your p12 certificate file you may convert it into pem files
  - you may use openssl to print content of your certificate (Subject line) like
```
openssl x509 -inform pem -in -noout -text -in <your-cert.pem>
```
- finally you will need to use CN field of you certificate to pass it to 
`fyne release ... -cert`

Finally, call
```
make build_ios
```
to build your iOS build.

Notes:
changed /Users/vk/tmp/fyne-2.2.3/cmd/fyne/internal/mobile/build_iosapp.go to
add CODE_SIGNING_ALLOWED=NO to xcrun command, see
https://developer.apple.com/forums/thread/69950

#### Additional notes
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

https://ioscodesigning.io/creating-code-signing-files/

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
