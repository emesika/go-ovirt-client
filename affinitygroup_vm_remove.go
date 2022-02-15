package ovirtclient

import (
    "fmt"
)

func (o *oVirtClient) RemoveVMFromAffinityGroup(
    clusterID ClusterID,
    vmID string,
    agID AffinityGroupID,
    retries ...RetryStrategy,
) error {
    retries = defaultRetries(retries, defaultWriteTimeouts())
    return retry(
        fmt.Sprintf("adding VM %s to affinity group %s", vmID, agID),
        o.logger,
        retries,
        func() error {
            _, err := o.conn.
                SystemService().
                ClustersService().
                ClusterService(string(clusterID)).
                AffinityGroupsService().
                GroupService(string(agID)).
                VmsService().
                VmService(vmID).
                Remove().
                Send()
            return err
        },
    )
}

func (m *mockClient) RemoveVMFromAffinityGroup(
    clusterID ClusterID,
    vmID string,
    agID AffinityGroupID,
    retries ...RetryStrategy,
) error {
    //TODO implement me
    panic("implement me")
}