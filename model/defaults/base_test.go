package defaults_test

import (
	"github.com/banzaicloud/banzai-types/components"
	"github.com/banzaicloud/banzai-types/components/amazon"
	"github.com/banzaicloud/banzai-types/components/azure"
	"github.com/banzaicloud/banzai-types/components/google"
	"github.com/banzaicloud/banzai-types/constants"
	"github.com/banzaicloud/pipeline/model/defaults"
	"github.com/banzaicloud/pipeline/utils"
	"testing"
)

func TestTableName(t *testing.T) {

	tableName := defaults.GKEProfile.TableName(defaults.GKEProfile{})
	if defaults.DefaultGoogleProfileTablaName != tableName {
		t.Errorf("Expected table name: %s, got: %s", defaults.DefaultGoogleProfileTablaName, tableName)
	}

}

func TestGetType(t *testing.T) {

	cases := []struct {
		name         string
		profile      defaults.ClusterProfile
		expectedType string
	}{
		{"type gke", &defaults.GKEProfile{}, constants.Google},
		{"type aks", &defaults.AKSProfile{}, constants.Azure},
		{"type aws", &defaults.AWSProfile{}, constants.Amazon},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			currentType := tc.profile.GetType()
			if tc.expectedType != currentType {
				t.Errorf("Expected cloud type: %s, got: %s", tc.expectedType, currentType)
			}
		})
	}

}

func TestUpdateWithoutSave(t *testing.T) {

	testCases := []struct {
		name           string
		basicProfile   defaults.ClusterProfile
		request        *components.ClusterProfileRequest
		expectedResult defaults.ClusterProfile
	}{
		{"full request GKE", &defaults.GKEProfile{}, fullRequestGKE, &fullGKE},
		{"just master update GKE", &defaults.GKEProfile{}, masterRequestGKE, &masterGKE},
		{"just node update GKE", &defaults.GKEProfile{}, nodeRequestGKE, &nodeGKE},
		{"just basic update GKE", &defaults.GKEProfile{}, emptyRequestGKE, &emptyGKE},

		{"full request AKS", &defaults.AKSProfile{}, fullRequestAKS, &fullAKS},
		{"just basic update AKS", &defaults.AKSProfile{}, emptyRequestAKS, &emptyAKS},

		{"full request AWS", &defaults.AWSProfile{}, fullRequestAWS, &fullAWS},
		{"just master update AWS", &defaults.AWSProfile{}, masterRequestAWS, &masterAWS},
		{"just node update AWS", &defaults.AWSProfile{}, nodeRequestAWS, &nodeAWS},
		{"just basic update AWS", &defaults.AWSProfile{}, emptyRequestAWS, &emptyAWS},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			err := tc.basicProfile.UpdateProfile(tc.request, false)

			if err != nil {
				t.Errorf("Expected error <nil>, got: %s", err.Error())
			}

			if err := utils.IsDifferent(tc.expectedResult, tc.basicProfile); err == nil {
				t.Errorf("Expected result: %#v, got: %#v", tc.expectedResult, tc.basicProfile)
			}

		})

	}

}

const (
	name               = "TestProfile"
	location           = "TestLocation"
	nodeInstanceType   = "TestNodeInstance"
	masterInstanceType = "TestMasterInstance"
	masterImage        = "TestMasterImage"
	nodeImage          = "TestMasterImage"
	version            = "TestVersion"
	nodeCount          = 1
	agentName          = "TestAgent"
	k8sVersion         = "TestKubernetesVersion"
	minCount           = 1
	maxCount           = 2
	spotPrice          = "0.2"
	serviceAccount     = "TestServiceAccount"
)

var (
	fullRequestGKE = &components.ClusterProfileRequest{
		Name:     name,
		Location: location,
		Cloud:    constants.Google,
		Properties: struct {
			Amazon *amazon.ClusterProfileAmazon `json:"amazon,omitempty"`
			Azure  *azure.ClusterProfileAzure   `json:"azure,omitempty"`
			Google *google.ClusterProfileGoogle `json:"google,omitempty"`
		}{
			Google: &google.ClusterProfileGoogle{
				Master: &google.Master{
					Version: version,
				},
				NodeVersion: version,
				NodePools: map[string]*google.NodePool{
					agentName: {
						Count:            nodeCount,
						NodeInstanceType: nodeInstanceType,
						ServiceAccount:   serviceAccount,
					},
				},
			},
		},
	}

	fullRequestAKS = &components.ClusterProfileRequest{
		Name:     name,
		Location: location,
		Cloud:    constants.Azure,
		Properties: struct {
			Amazon *amazon.ClusterProfileAmazon `json:"amazon,omitempty"`
			Azure  *azure.ClusterProfileAzure   `json:"azure,omitempty"`
			Google *google.ClusterProfileGoogle `json:"google,omitempty"`
		}{
			Azure: &azure.ClusterProfileAzure{
				KubernetesVersion: k8sVersion,
				NodePools: map[string]*azure.NodePoolCreate{
					agentName: {
						Count:            nodeCount,
						NodeInstanceType: nodeInstanceType,
					},
				},
			},
		},
	}

	fullRequestAWS = &components.ClusterProfileRequest{
		Name:     name,
		Location: location,
		Cloud:    constants.Amazon,
		Properties: struct {
			Amazon *amazon.ClusterProfileAmazon `json:"amazon,omitempty"`
			Azure  *azure.ClusterProfileAzure   `json:"azure,omitempty"`
			Google *google.ClusterProfileGoogle `json:"google,omitempty"`
		}{
			Amazon: &amazon.ClusterProfileAmazon{
				Master: &amazon.AmazonProfileMaster{
					InstanceType: masterInstanceType,
					Image:        masterImage,
				},
				NodePools: map[string]*amazon.NodePool{
					agentName: {
						InstanceType: nodeInstanceType,
						SpotPrice:    spotPrice,
						MinCount:     minCount,
						MaxCount:     maxCount,
						Image:        nodeImage,
					},
				},
			},
		},
	}

	fullGKE = defaults.GKEProfile{
		DefaultModel:  defaults.DefaultModel{Name: name},
		Location:      location,
		NodeVersion:   version,
		MasterVersion: version,
		NodePools: []*defaults.GKENodePoolProfile{
			{
				Count:            nodeCount,
				NodeInstanceType: nodeInstanceType,
				NodeName:         agentName,
			},
		},
	}

	fullAKS = defaults.AKSProfile{
		DefaultModel:      defaults.DefaultModel{Name: name},
		Location:          location,
		KubernetesVersion: k8sVersion,
		NodePools: []*defaults.AKSNodePoolProfile{
			{
				NodeInstanceType: nodeInstanceType,
				Count:            nodeCount,
				NodeName:         agentName,
			},
		},
	}

	fullAWS = defaults.AWSProfile{
		DefaultModel:       defaults.DefaultModel{Name: name},
		Location:           location,
		MasterInstanceType: masterInstanceType,
		MasterImage:        masterImage,
		NodePools: []*defaults.AWSNodePoolProfile{
			{
				InstanceType: nodeInstanceType,
				NodeName:     agentName,
				SpotPrice:    spotPrice,
				MinCount:     minCount,
				MaxCount:     maxCount,
				Image:        nodeImage,
			},
		},
	}
)

var (
	masterRequestGKE = &components.ClusterProfileRequest{
		Name:     name,
		Location: location,
		Cloud:    constants.Google,
		Properties: struct {
			Amazon *amazon.ClusterProfileAmazon `json:"amazon,omitempty"`
			Azure  *azure.ClusterProfileAzure   `json:"azure,omitempty"`
			Google *google.ClusterProfileGoogle `json:"google,omitempty"`
		}{
			Google: &google.ClusterProfileGoogle{
				Master: &google.Master{
					Version: version,
				},
			},
		},
	}

	masterRequestAWS = &components.ClusterProfileRequest{
		Name:     name,
		Location: location,
		Cloud:    constants.Amazon,
		Properties: struct {
			Amazon *amazon.ClusterProfileAmazon `json:"amazon,omitempty"`
			Azure  *azure.ClusterProfileAzure   `json:"azure,omitempty"`
			Google *google.ClusterProfileGoogle `json:"google,omitempty"`
		}{
			Amazon: &amazon.ClusterProfileAmazon{
				Master: &amazon.AmazonProfileMaster{
					InstanceType: masterInstanceType,
					Image:        masterImage,
				},
			},
		},
	}

	masterGKE = defaults.GKEProfile{
		DefaultModel:  defaults.DefaultModel{Name: name},
		Location:      location,
		MasterVersion: version,
	}

	masterAWS = defaults.AWSProfile{
		DefaultModel:       defaults.DefaultModel{Name: name},
		Location:           location,
		MasterInstanceType: masterInstanceType,
		MasterImage:        masterImage,
	}
)

var (
	nodeRequestGKE = &components.ClusterProfileRequest{
		Name:     name,
		Location: location,
		Cloud:    constants.Google,
		Properties: struct {
			Amazon *amazon.ClusterProfileAmazon `json:"amazon,omitempty"`
			Azure  *azure.ClusterProfileAzure   `json:"azure,omitempty"`
			Google *google.ClusterProfileGoogle `json:"google,omitempty"`
		}{
			Google: &google.ClusterProfileGoogle{
				NodeVersion: version,
				NodePools: map[string]*google.NodePool{
					agentName: {
						Count:            nodeCount,
						NodeInstanceType: nodeInstanceType,
					},
				},
			},
		},
	}

	nodeRequestAWS = &components.ClusterProfileRequest{
		Name:     name,
		Location: location,
		Cloud:    constants.Amazon,
		Properties: struct {
			Amazon *amazon.ClusterProfileAmazon `json:"amazon,omitempty"`
			Azure  *azure.ClusterProfileAzure   `json:"azure,omitempty"`
			Google *google.ClusterProfileGoogle `json:"google,omitempty"`
		}{
			Amazon: &amazon.ClusterProfileAmazon{
				NodePools: map[string]*amazon.NodePool{
					agentName: {
						InstanceType: nodeInstanceType,
						SpotPrice:    spotPrice,
						MinCount:     minCount,
						MaxCount:     maxCount,
						Image:        nodeImage,
					},
				},
			},
		},
	}

	nodeGKE = defaults.GKEProfile{
		DefaultModel: defaults.DefaultModel{Name: name},
		Location:     location,
		NodeVersion:  version,
		NodePools: []*defaults.GKENodePoolProfile{
			{
				Count:            nodeCount,
				NodeInstanceType: nodeInstanceType,
				NodeName:         agentName,
			},
		},
	}

	nodeAWS = defaults.AWSProfile{
		DefaultModel: defaults.DefaultModel{Name: name},
		Location:     location,
		NodePools: []*defaults.AWSNodePoolProfile{
			{
				InstanceType: nodeInstanceType,
				NodeName:     agentName,
				SpotPrice:    spotPrice,
				MinCount:     minCount,
				MaxCount:     maxCount,
				Image:        nodeImage,
			},
		},
	}
)

var (
	emptyRequestGKE = &components.ClusterProfileRequest{
		Name:     name,
		Location: location,
		Cloud:    constants.Google,
	}

	emptyRequestAKS = &components.ClusterProfileRequest{
		Name:     name,
		Location: location,
		Cloud:    constants.Azure,
	}

	emptyRequestAWS = &components.ClusterProfileRequest{
		Name:     name,
		Location: location,
		Cloud:    constants.Amazon,
	}

	emptyGKE = defaults.GKEProfile{
		DefaultModel: defaults.DefaultModel{Name: name},
		Location:     location,
	}

	emptyAKS = defaults.AKSProfile{
		DefaultModel: defaults.DefaultModel{Name: name},
		Location:     location,
	}

	emptyAWS = defaults.AWSProfile{
		DefaultModel: defaults.DefaultModel{Name: name},
		Location:     location,
	}
)
