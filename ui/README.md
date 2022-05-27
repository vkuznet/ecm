
- run Android Studio
- run emulator and check devices
  - ./adb devices
- deploy APP
  - ./adb install ecm_ui.apk
  - check app log
    - ./adb logcat | grep fyne
