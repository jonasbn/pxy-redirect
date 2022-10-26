# pxy-redirect

## Introduction

This is service I created for deployment on [DigitalOcean][DO]. The service has been replace by a serverless function [pxy-redirect-ow-function][FUNCTION]

It supports some documentation I have written up, which resulted in a bug. The bug could be observed when displaying a Markdown page on GitHub. which exceeded the size limit of what could be rendered a Markdown.

The service is available at: [pxy.fi]

## Description

The service takes a HTTP request, parses it and redirects to the designated documentation page on the [releases.llvm.org][LLVM] website.

The acceptable URLs being:

- `/<version>/<fragment>`

So the first part has to be a number and the second part a string corresponding to the equivalens on the [releases.llvm.org][LLVM] website.

Example:

`/5/rsanitize-address`

Would be parsed and altered and redirected to:

`https://releases.llvm.org/5.0.0/tools/clang/docs/DiagnosticsReference.html#rsanitize-address`

## Diagnostics

This is a collection of errors which can be emitted from the service. Not all are visible to the end user and not all error scenarios are documented.

This section and documentation is primarily aimed and what can be recovered from.

### Unable to assemble URL (400)

This is the most common error it will provide additional information as to why the request was regarded as a bad request.

#### insufficient parts in provided url

The the request does not contain enough parts to assemble the redirect target URL.

#### first part of url is not a number

The first part of the URL should be a number (integer), which is translated to a version number.

#### second part of url is not a string

The second part of the URL should be a string.

### Unable to parse received URL error (500)

The URL could not be parsed.

### `ListenAndServe` fatal error
  
This error is emitted if the service is unable to start. The error message will contain details.

The error is not visible to the end-user as such, but the service will be unavailble.

## Resources and References

- [pxy-redirect-ow-function][FUNCTION]
- [My TIL collection: clang diagnostic flags](https://github.com/jonasbn/til/blob/master/clang/diagnostic_flags.md) (GitHub)
- [My TIL collection: clang diagnostic flags](http://jonasbn.github.io/til/clang/diagnostic_flags.html) (website)
- [clang diagnostic flags matrix generator](https://github.com/jonasbn/clang-diagnostic-flags-matrix)
- [llvm releases documentation site][LLVM]
- [DigitalOcean][DO]

[DO]: https://www.digitalocean.com/
[LLVM]: https://releases.llvm.org/
[pxy.fi]: https://pxy.fi/
[FUNCTION]: https://github.com/jonasbn/pxy-redirect-ow-function
