### ECM UI interface
This area contains code for ECM UI based on [fyne.io](https://fyne.io/)

#### Helpful hints
Here we provide helpful hints about ECM UI setup and builds. To build ECM
for mobile please follow the following steps:
- build appropriate application package for your mobile
  - for Andoird please use `make build_android` while for iOS `make build_iphone`
- run Android Studio or Xcode
  - for Android
    - run emulator and check devices `./adb devices`
    - deploy ECM applocation `./adb install ecm.apk`
    - check application log `./adb logcat | grep fyne`
    - on macOS we write sync config to `$HOME/Library/Preferences/fyne/io.github.USERNAME`
