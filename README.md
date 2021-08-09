# goVirt: an easy-to-use overlay for the oVirt Go SDK

<p align="center"><strong>⚠⚠⚠ This library is early in development. ⚠⚠⚠</strong></p>

This library is early in development, and the API may change at any time until version 1.0.0. We hope to stabilize the API soon, providing the core functionality on an as-needed basis. If you need an API integrated, please open an issue. 

This library provides an easy-to-use overlay for the automatically generated [Go SDK for oVirt](https://github.com/ovirt/go-ovirt). It does *not* replace the Go SDK. It implements the functions of the SDK only partially and is primarily used by the [oVirt Terraform provider](https://github.com/ovirt/terraform-provider-ovirt/).

## Using this library

To use this library you will have to include it as a Go module dependency:

```
go get github.com/ovirt/go-ovirt-client
```

You can then create a client instance like this:

```go
package main

import (
    "github.com/ovirt/go-client-log"
    "github.com/ovirt/go-ovirt-client"
)

func main() {
    // Create a logger that logs to the standard Go log here:
    logger := ovirtclientlog.NewGoLogLogger(nil)

    // Create a ovirtclient.TLSProvider implementation. This allows for simple
    // TLS configuration.
	tls := ovirtclient.TLS()

	// Add certificates from an in-memory byte slice. Certificates must be in PEM format.
	tls.CACertsFromMemory(caCerts)
	
	// Add certificates from a single file. Certificates must be in PEM format.
	tls.CACertsFromFile("/path/to/file.pem")

	// Add certificates from a directory. Optionally, regular expressions can be passed that must match the file
	// names.
	tls.CACertsFromDir("/path/to/certs", regexp.MustCompile(`\.pem`)) 

	// Add system certificates
	tls.CACertsFromSystem()

	// Disable certificate verification. This is a bad idea, please don't do this.
	tls.Insecure()

    // Create a new goVirt instance:
    client, err := ovirtclient.New(
        // URL to your oVirt engine API here:
        "https://your-ovirt-engine/ovirt-engine/api/",
        // Username here:
        "admin@internal",
        // Password here:
        "password-here",
        // Pass the TLS provider here:
        tls,
        // Pass the logger here:
        logger,
        // Pass in extra settings here. Must implement the ovirtclient.ExtraSettings interface.
        nil,
    )
    if err != nil {
        // Handle error, here in a really crude way:
        panic(err)
    }
    // Use client. Please use the code completion in your IDE to
    // discover the functions. Each is well documented.
    upload, err := client.StartImageUpload(
        //...
    )
    //....
}
```

## Test helper

This library also provides a test helper for integration testing against the oVirt engine. It allows for automatically discovering a usable storage domain, host, clusters, etc:

```go
package main

import (
    "os"
    "testing"

    ovirtclient "github.com/ovirt/go-ovirt-client"
    ovirtclientlog "github.com/ovirt/go-ovirt-client-log"
)

func TestSomething(t *testing.T) {
    // Create a logger that logs to the standard Go log here:
    logger := ovirtclientlog.NewTestLogger(t)
    // Set to true to use in-memory mock, see below
    mock := false
    
    tls := ovirtclient.TLS()
    
    if caFile := os.Getenv("OVIRT_CAFILE"); caFile != "" {
    	tls.CACertsFromFile(caFile)
    }
    if caBundle := os.Getenv("OVIRT_CABUNDLE"); caBundle != "" {
        tls.CACertsFromMemory([]byte(caBundle))
	}
	if os.Getenv("OVIRT_INSECURE") != "" {
		tls.Insecure()
    }
    
    // Create the test helper
    helper, err := ovirtclient.NewTestHelper(
        os.Getenv("OVIRT_URL"),
        os.Getenv("OVIRT_USER"),
        os.Getenv("OVIRT_PASSWORD"),
        tls,
        os.Getenv("OVIRT_CLUSTER_ID"),
        os.Getenv("OVIRT_BLANK_TEMPLATE_ID"),
        os.Getenv("OVIRT_STORAGE_DOMAIN_ID"),
        mock,
        logger,
    )
    if err != nil {
        t.Fatal(err)
    }
    // Fetch the cluster ID for testing
    clusterID := helper.GetClusterID()
    //...
}
```

**Tip:** You can use any logger that satisfies the `Logger` interface described in [go-ovirt-client-log](https://github.com/oVirt/go-ovirt-client-log)

## Retries

This library attempts to retry API calls that can be retried if possible. Each function has a sensible retry policy. However, you may want to customize the retries by passing one or more retry flags. The following retry flags are supported:

- `ovirtclient.ContextStrategy(ctx)`: this strategy will stop retries when the context parameter is canceled.
- `ovirtclient.ExponentialBackoff(factor)`: this strategy adds a wait time after each time, which is increased by the given factor on each try. The default is a backoff with a factor of 2.
- `ovirtclient.AutoRetry()`: this strategy will cancel retries if the error in question is a permanent error. This is enabled by default.
- `ovirtclient.MaxTries(tries)`: this strategy will abort retries if a maximum number of tries is reached. On complex calls the retries are counted per underlying API call.
- `ovirtclient.Timeout(duration)`: this strategy will abort retries if a certain time has been elapsed for the higher level call.
- `ovirtclient.CallTimeout(duration)`: this strategy will abort retries if a certain underlying API call takes longer than the specified duration. 

## Mock client

This library also provides a mock oVirt client that doesn't need working oVirt engine to function. It stores all information in-memory and simulates a working oVirt system. You can instantiate the mock client like so:

```go
client := ovirtclient.NewMock()
```

We recommend using the `ovirtclient.Client` interface as a means to declare it as a dependency in your factory so you can pass both the mock and the real connection as a parameter:

```go
func NewMyoVirtUsingUtility(
    client ovirtclient.Client,
) *myOVirtUsingUtility {
    return &myOVirtUsingUtility{
        client: client,
    }
}
``` 

## FAQ

### Why doesn't the library return the underlying oVirt SDK objects?

It's a painful decision we made. We want to encourage anyone who needs a certain function to submit a PR instead of simply relying on the SDK objects. This will lead to some overhead when a new function needs to be added, but leads to cleaner code in the end and makes this library more comprehensive. It also makes it possible to create the mock client, which would not be possibly if we had to simulate all parts of the oVirt engine.

If you need to access the oVirt SDK client you can do so from the `ovirtclient.New()` function:

```go
client, err := ovirtclient.New(
    //...
)
if err != nil {
    //...
}
sdkClient := client.GetSDKClient()
```

You can also get a properly preconfigured HTTP client if you need it:

```go
httpClient := client.GetHTTPClient()
```

**⚠ Warning:** If you code relies on the SDK or HTTP clients you will not be able to use the mock functionality described above for testing.