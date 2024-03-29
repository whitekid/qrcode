import "@typespec/rest";
import "@typespec/openapi3";

using TypeSpec.Http;

namespace api {
    @service({
        title: "qrcode generate api",
        version: "v1",
    })
    @server("https://qrcodeapi.woosum.net/api", "")
    @summary("simple qrcode generate service")
    @route("v1")
    @tag("v1")
    namespace v1 {
        model QRCode {
            @header contentType: "image/png" | "image/jpeg" | "image/gif";

            @summary("image content")
            @body
            qrcode: bytes;
        }

        model CommonParams {
            @summary("image width")
            @query
            w?: numeric = 200;

            @summary("image height")
            @query
            h?: numeric = 200;
            @header accept?: string = "image/png";
        }

        @error
        model Error {
            @statusCode statusCode: 400 | 415 | 429 | 500;
        }

        @route("qrcode")
        interface QRcode {
            @summary("generate text, url and ssid qrcode")
            @doc("doc text, url and ssid qrcode")
            @get
            generate(
                @doc("any text content")
                @query
                @maxLength(1024)
                content?: string,

                @doc("url like https://example.com")
                @query
                url?: string,

                @doc("ssid")
                @query
                @maxLength(20)
                ssid?: string,

                @doc("ssid only: auth type")
                @query
                auth?: "WEP" | "WPA" | "WPA2",

                @doc("ssid only: password")
                @query
                @secret
                @maxLength(20)
                pass?: string,

                @doc("ssid only: hidden ssid")
                @query
                hidden?: boolean,

                @doc("ssid only: EAP method")
                @query
                @maxLength(20)
                eap?: string,

                @doc("ssid only")
                @query
                @maxLength(20)
                anon?: string,

                @doc("ssid only")
                @query
                @maxLength(20)
                ident?: string,

                @doc("ssid only")
                @query
                @maxLength(20)
                ph2?: string,
                ...CommonParams
            ): QRCode | Error;
        }

        @route("contact")
        interface Contact {
            @summary("generate contact qrcode")
            @get
            generate(
                @query
                @maxLength(100)
                "name[first]"?: string,

                @query
                @maxLength(100)
                "name[list]"?: string,

                @query
                @maxLength(100)
                "name[last]"?: string,

                @query
                @maxLength(100)
                company?: string,

                @query
                @maxLength(100)
                department?: string,

                @query
                @maxLength(100)
                title?: string,

                @query
                @format("email")
                @maxLength(100)
                email?: string,

                @query
                @format("email")
                @maxLength(100)
                "email[home]"?: string,

                @query
                @format("email")
                @maxLength(100)
                "email[work]"?: string,

                @query
                @maxLength(100)
                tel?: string,

                @query
                @maxLength(100)
                "tel[home]"?: string,

                @query
                @maxLength(100)
                "tel[work]"?: string,

                @query
                @maxLength(100)
                mobile?: string,

                @query
                @maxLength(100)
                pager?: string,

                @query
                @maxLength(100)
                "fax[home]"?: string,

                @query
                @maxLength(100)
                "fax[work]"?: string,

                @query
                @maxLength(100)
                "addr[home][postcode]"?: string,

                @query
                @maxLength(100)
                "addr[home][country]"?: string,

                @query
                @maxLength(100)
                "addr[home][province]"?: string,

                @query
                @maxLength(100)
                "addr[home][city]"?: string,

                @query
                @maxLength(100)
                "addr[home][street]"?: string,

                @query
                @maxLength(100)
                "addr[home][street2]"?: string,

                @query
                @maxLength(100)
                "addr[work][postcode]"?: string,

                @query
                @maxLength(100)
                "addr[work][country]"?: string,

                @query
                @maxLength(100)
                "addr[work][province]"?: string,

                @query
                @maxLength(100)
                "addr[work][city]"?: string,

                @query
                @maxLength(100)
                "addr[work][street]"?: string,

                @query
                @maxLength(100)
                "addr[work][street2]"?: string,

                @query
                @maxLength(100)
                note?: string,
                ...CommonParams
            ): QRCode | Error;
        }

        @route("vcard")
        interface VCard {
            @summary("generate vcard qrcode")
            @post
            generate(
                @body vcard: bytes,
                @header contentType: "text/vcard",
                ...CommonParams
            ): QRCode | Error;
        }

        @route("vevent")
        interface VEvent {
            @summary("generate vevent qrcode")
            @post
            generate(
                @body vevent: bytes,
                @header contentType: "text/vevent",
                ...CommonParams
            ): QRCode | Error;
        }
    }
}
