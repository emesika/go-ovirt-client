package ovirtclient

import (
	"math/rand"
	"sync"

	"github.com/google/uuid"
)

// MockClient provides in-memory client functions, and additionally provides the ability to inject
// information.
type MockClient interface {
	Client

	// GenerateUUID generates a UUID for testing purposes.
	GenerateUUID() string
}

type mockClient struct {
	logger                            Logger
	url                               string
	lock                              *sync.Mutex
	nonSecureRandom                   *rand.Rand
	vms                               map[string]*vm
	storageDomains                    map[string]*storageDomain
	disks                             map[string]*diskWithData
	clusters                          map[string]*cluster
	hosts                             map[string]*host
	templates                         map[TemplateID]*template
	nics                              map[string]*nic
	vnicProfiles                      map[string]*vnicProfile
	networks                          map[string]*network
	dataCenters                       map[string]*datacenterWithClusters
	vmDiskAttachmentsByVM             map[string]map[string]*diskAttachment
	vmDiskAttachmentsByDisk           map[string]*diskAttachment
	templateDiskAttachmentsByTemplate map[TemplateID][]*templateDiskAttachment
	templateDiskAttachmentsByDisk     map[string]*templateDiskAttachment
	tags                              map[string]*tag
}

func (m *mockClient) GetURL() string {
	return m.url
}

func (m *mockClient) GenerateUUID() string {
	return uuid.NewString()
}
