# QRCode API

Just simple QRCode API Service

## Usage

### Text

![TEXT](https://qrcode.woosum.net/api/v1/qrcode?content=HELLO)

<https://qrcode.woosum.net/api/v1/qrcode?content=HELLO>

### URL

with content:

![URL](https://qrcode.woosum.net/api/v1/qrcode?content=https://github.com)

<https://qrcode.woosum.net/api/v1/qrcode?content=https://github.com>

with url:

![URL](https://qrcode.woosum.net/api/v1/qrcode?url=github.com)

<https://qrcode.woosum.net/api/v1/qrcode?url=github.com>

### Join WIFI

![WIFI](https://qrcode.woosum.net/api/v1/qrcode?ssid=MySSID&auth=WPA&pass=mypassword)

<https://qrcode.woosum.net/api/v1/qrcode?ssid=MySSID&auth=WPA&pass=mypassword>

### Contact

![Contact](https://qrcode.woosum.net/api/v1/contact?name[last]=Choe&name[first]=Cheng%20Dae)

<https://qrcode.woosum.net/api/v1/contact?name[last]=Choe&name[first]=Cheng20Dae>

#### with vcard

    curl -X POST https://qrcode.woosum.net/api/v1/vcard \
        -H "content-type: text/vcard" \
        -H "accept: image/png" \
        -d "BEGIN:VCARD
    VERSION:4.0
    N:lastname;firstname;;;
    END:VCARD"

    HTTP/1.1 200 OK
    Content-Type: image/png

### Event

    curl -X POST https://qrcode.woosum.net/api/v1/vevent \
        -H "content-type: text/vevent" \
        -H "accept: image/png" \
        -d "BEGIN:VEVENT
    UID:19970901T130000Z-123401@host.com
    DTSTAMP:19970901T1300Z
    DTSTART:19970903T163000Z
    DTEND:19970903T190000Z
    SUMMARY:Annual Employee Review
    CLASS:PRIVATE
    CATEGORIES:BUSINESS,HUMAN RESOURCES
    END:VEVENT"

    HTTP/1.1 200 OK
    Content-Type: image/png

## more code formsts

<https://github.com/zxing/zxing/wiki/Barcode-Contents>
