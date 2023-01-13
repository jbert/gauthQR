# GauthQR

This is a tool to read an image containing QR code and output a line for
the `gauth.csv` used by the gauth tool (https://github.com/pcarrier/gauth/).

This can be used to export OTP codes from google authenticator without root
on your android device.

My phone didn't let me screenshot the exported QR code, but using another
phone (or webcam) works well enough.

## Dependencies

To extract the QR data, we need `zbarimg` to be in your path. On ubuntu:

```
apt install zbar-tools
```

## TODO

Output rich URLs for formats which don't match the gauth short line
constraints.
