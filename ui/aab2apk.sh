#!/bin/bash
rm -f $PWD/ecm.apks
bundletool build-apks --bundle=$PWD/ecm.aab --output=$PWD/ecm.apks --mode=universal
unzip -p ecm.apks universal.apk > ecm.apk
rm -f $PWD/ecm.apks
