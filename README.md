# QRCode API

Just simple QRCode API Service

## Usage

### Text

![TEXT](https://qrcodeapi.woosum.net/v1/qrcode?content=HELLO)

<https://qrcodeapi.woosum.net/v1/qrcode?content=HELLO>

### URL

with content:

![URL](https://qrcodeapi.woosum.net/v1/qrcode?content=https://github.com)

<https://qrcodeapi.woosum.net/v1/qrcode?content=https://github.com>

with url:

![URL](https://qrcodeapi.woosum.net/v1/qrcode?url=github.com)

<https://qrcodeapi.woosum.net/v1/qrcode?url=github.com>

### Join WIFI

![WIFI](https://qrcodeapi.woosum.net/v1/qrcode?ssid=MySSID&auth=WPA&pass=mypassword)

<https://qrcodeapi.woosum.net/v1/qrcode?ssid=MySSID&auth=WPA&pass=mypassword>

### Contact

![Contact](https://qrcodeapi.woosum.net/v1/contact?name[last]=Choe&name[first]=Cheng%20Dae)

<https://qrcodeapi.woosum.net/v1/contact?name[last]=Choe&name[first]=Cheng20Dae>

#### with vcard

    POST https://qrcodeapi.woosum.net/contact
    content-type: text/vcard

    BEGIN:VCARD
    VERSION:4.0
    N:lastname;firstname;;;
    END:VCARD

    HTTP/1.1 200 OK
    Content-Type: image/png

## more code formsts

<https://github.com/zxing/zxing/wiki/Barcode-Contents>
