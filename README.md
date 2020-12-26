# Traefik Pilot Plugin: duplicateheader

## TL;DR

Duplicates a header to one or more headers.

## Description

It can be used to get the real client ip (`X-Real-Ip`) when routing your traffic through Cloudflare or any other CDN.  
If (one of) the destination headers do not exist, it will create the header automatically.  
If the source header does not exist, it will do nothing.

## Examples

```yaml
pilot:
  token: "xxxx"
experimental:
  plugins:
    duplicateheader:
      moduleName: "github.com/containeroo/duplicateheader"
      version: "v1.0.8"
```

```yaml
middlewares:
  my-traefik-plugin-duplicateheader:
    plugin:
      duplicateheader:
        source: Cf-Connecting-Ip
        destination:
          - X-Real-Ip
          - X-Forwarded-For
```

**This plugin is not affiliated with Cloudflare.**
