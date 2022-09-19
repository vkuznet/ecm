### How to install ECM on Android platform
- install Android Studio and install SDK with NDK libs
  - go to Preferences -> Appearnace & benavior -> Android SDK 
  - choose SDK platform you want to use, e.g. Android 8.0
  - go to SDK tools tab and select everything with SDK/NDK, CMake,
  Emulator, Good Web/Play services
- make a soft link to your `adb` tool
```
ln -s /path/Library/Android/sdk/platform-tools/adb ./adb
```
- create new keystore
```
keytool -genkey -v -keystore ecm.keystore -alias ecm-key -keyalg RSA -keysize 2048 -validity 10000

# and you'll see the following output
# please enter relevant information for provided fields
Enter keystore password:
Re-enter new password:
What is your first and last name?
  [Unknown]:  First Last
What is the name of your organizational unit?
  [Unknown]:  Bla
What is the name of your organization?
  [Unknown]:  Org
What is the name of your City or Locality?
  [Unknown]:  City
What is the name of your State or Province?
  [Unknown]:  State
What is the two-letter country code for this unit?
  [Unknown]:  US
Is CN=First Last, OU=Bla, O=Org, L=City, ST=State, C=US correct?
  [no]:  yes

Generating 2,048 bit RSA key pair and self-signed certificate (SHA256withRSA) with a validity of 10,000 days
	for: CN=First Last, OU=Bla, O=Org, L=City, ST=State, C=US
[Storing ecm.keystore]
```
At this step you'll get `ecm.keystore` in your local area.

- build new Android bundle
```
rm ecm.aab
make build_android
```
At this step you will get `ecm.aab` file

- Convert aab (Android Application Bundle) file to apk (Android Package) using
`bundetool`, see [3, 4]
```
bundletool build-apks --bundle=$PWD/ecm.aab --output=$PWD/ecm.apks --mode=universal

# extract apk file
unzip -p ecm.apks universal.apk > ecm.apk
```
- now you can plug your Android phone and use the following command
to transfer your apk file to the phone:
```
# plug the phone
# check if you can see devices
./adb devices

# deploy ECM applocation 
./adb install ecm.apk

# check application logs
./adb logcat | grep fyne
```

### How to sign your app manually.
Generate a private key using keytool. For example, see [1]

```
$ keytool -genkey -v -keystore ecm.keystore -alias ecm-key -keyalg RSA -keysize 2048 -validity 10000
```

This example prompts you for passwords for the keystore and key, and to provide the Distinguished Name fields for your key. It then generates the keystore as a file called ecm.keystore. The keystore contains a single key, valid for 10000 days. The alias is a name that you will use later when signing your app.

Compile your app in release mode to obtain an unsigned APK.
Sign your app with your private key using jarsigner:
```
$ jarsigner -verbose -sigalg SHA1withRSA -digestalg SHA1 -keystore ecm.keystore ecm-app.apk ecm-key
```

This example prompts you for passwords for the keystore and key. It then modifies the APK in-place to sign it. Note that you can sign an APK multiple times with different keys.

Verify that your APK is signed. For example:
```
$ jarsigner -verify -verbose -certs ecm-app.apk
```

Align the final APK package using zipalign.
```
$ zipalign -v 4 ecm-unaligned.apk ecm.apk
```
zipalign ensures that all uncompressed data starts with a particular byte alignment relative to the start of the file, which reduces the amount of RAM consumed by an app.

### How to convert aab to apk
Refer bundletool command and this post [2]

For Debug apk command,

```
bundletool build-apks --bundle=$PWD/ecm.aab --output=$PWD/ecm.apks
```

For Release apk command,
```
# it assumes we create
# keystore.pwd file with password used for keystore
# key.pwd file with password used for keystore key
bundletool build-apks --bundle=$PWD/ecm.aab --output=$PWD/ecm.apks
--ks=$PWD/ecm.keystore
--ks-pass=file:$PWD/keystore.pwd
--ks-key-alias=ecm-key
--key-pass=file:$PWD/key.pwd
```

### References
1. https://stackoverflow.com/questions/25975320/create-android-keystory-private-key-command-line
2. https://stackoverflow.com/questions/53040047/generate-apk-file-from-aab-file-android-app-bundle
3. https://developer.android.com/studio/command-line/bundletool
4. https://www.wikiwand.com/en/APK_%28file_format%29
