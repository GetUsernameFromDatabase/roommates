# Certificates

Local development certificate for HTTPS. Not used when `TRUSTED_PLATFORM` is
set, the platform (Cloudflare, Fly.io) then holds the public certificate.

The browser must trust it: with a "proceed anyway" exception Chrome still
refuses to register the service worker, so push and the offline page break.

## Install mkcert

<https://github.com/FiloSottile/mkcert#installation> (e.g. `winget install FiloSottile.mkcert`),
then trust its local CA once:

```sh
mkcert -install
```

## Generate the certificate

From `app/certificates/`:

```sh
mkcert -cert-file server.pem -key-file server.key localhost 127.0.0.1 ::1
```

To test from a phone on the LAN, add the PC's LAN IP to the names and install
the root CA from `mkcert -CAROOT` (`rootCA.pem`) on the phone.
