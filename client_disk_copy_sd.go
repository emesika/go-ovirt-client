package ovirtclient

import (
    "fmt"
    ovirtsdk "github.com/ovirt/go-ovirt"
    "sync"
)

func (o *oVirtClient) CopyDiskToStorageDomain(
    id string,
    storageDomainID string,
    retries ...RetryStrategy) (result Disk, err error) {
    retries = defaultRetries(retries, defaultLongTimeouts())
    progress,err := o.StartCopyDiskToStorageDomain(id,storageDomainID, retries...)
    if err != nil {
        return progress.Disk(), err
    }
    return progress.Wait(retries...)
}

func (o *oVirtClient) StartCopyDiskToStorageDomain(id string,storageDomainID string,retries ...RetryStrategy) (	DiskUpdate,
    error,){

    o.logger.Infof("Starting copy disk to different storage domain.")
    retries = defaultRetries(retries, defaultWriteTimeouts())
    sdkDisk := ovirtsdk.NewDiskBuilder().Id(id)

    correlationID := fmt.Sprintf("disk_copy_%s", generateRandomID(5, o.nonSecRand))

    var disk Disk

    err := retry(
        fmt.Sprintf("copying disk %s", id),
        o.logger,
        retries,
        func() error {
            _, err := o.conn.
                SystemService().
                DisksService().
                DiskService(id).
                Copy().
                StorageDomain(ovirtsdk.NewStorageDomainBuilder().Id(storageDomainID).MustBuild()).
                Disk(sdkDisk.MustBuild()).
                Query("correlation_id", correlationID).
                Send()
            if err != nil {
                return err
            }

            disk, err = convertSDKDisk(sdkDisk.MustBuild(), o)
            if err != nil {
                return wrap(err, EUnidentified, "failed to convert SDK disk object")
            }
            return nil
        },
    )
    if err != nil {
        return nil, err
    }
    return &diskWait{
        client:        o,
        disk:          disk,
        correlationID: correlationID,
        lock:          &sync.Mutex{},
    }, nil




}


