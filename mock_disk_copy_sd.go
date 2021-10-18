package ovirtclient

import (
	"time"
)


func (m *mockClient) CopyDiskToStorageDomain(id string, storageDomainID string, retries ...RetryStrategy) (result Disk, err error){

	progress, err := m.StartCopyDiskToStorageDomain(id, storageDomainID, retries...)
	if err != nil {
		return progress.Disk(), err
	}
	return progress.Wait(retries...)
}

func (m *mockClient) StartCopyDiskToStorageDomain(id string, storageDomainID string, _ ...RetryStrategy) (
	DiskUpdate,
	error,
) {
	m.lock.Lock()
	defer m.lock.Unlock()

	disk, ok := m.disks[id]
	if !ok {
		return nil, newError(ENotFound, "disk with ID %s not found", id)
	}
	if err := disk.Lock(); err != nil {
		return nil, err
	}
	update := &mockDiskCopy{
		client: m,
		disk:   disk,
		done:   make(chan struct{}),
	}
	defer update.do()
	return update, nil
}

type mockDiskCopy struct {
	client *mockClient
	disk   *diskWithData
	done   chan struct{}
}

func (c *mockDiskCopy) Disk() Disk {
	c.client.lock.Lock()
	defer c.client.lock.Unlock()

	return c.disk
}

func (c *mockDiskCopy) Wait(_ ...RetryStrategy) (Disk, error) {
	<-c.done

	return c.disk, nil
}

func (c *mockDiskCopy) do() {
	// Sleep to trigger potential race conditions / improper status handling.
	time.Sleep(time.Second)

	c.client.disks[c.disk.ID()] = c.disk
	c.disk.Unlock()

	close(c.done)
}
