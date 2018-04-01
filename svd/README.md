# SVD

### A SPIFFE ID Verification Daemon

SVD provides a simple way to check if an SVID
can connect to a service. By connecting to svd using an SVID (an
x509 cert over TLS) you can verify if the svid is correct for the SPIFFE
URI specified for SVD. By default, svd runs as an htts server
on port 8443.

#### Running SVD

To run svd, simply ensure that all the important files exist in the same
directory, and run

    svd

##### Important Files
In the directory where svd lives there are several important files.

File | Contents
---- | --------
URI | Contains the SPIFFE URI for the service to be checked, something like "spiffe://foo.com/path/service"
svd.crt, svd.key | The cert and key for the TLS server (x.509 and pkcs#12, respectively)
ROOT.crt | the root CA cert for SPIFFE verication (x.509)
INT.crt | the intermediated cert for SPIFFE verification (x.509)

#### General Design
svd is based on the "net/http" web server in the golang standard library.
A simple handler uses the "spiffe-go" library to verify the SVID provided
on the incoming connection.

#### Building svd

You will first need to install the go-spiffe module from github:

    go get -u -v github.com/spiffe/go-spiffe

Then download the source and run:

    go build

Everything should now be ready to run.