//
// Copyright (C) IBM Corporation 2016, Livio Soares <lsoares@us.ibm.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package retrieve_and_rank

import (
	"os"
	"testing"
	"time"

	"github.com/liviosoares/go-watson-sdk/watson"
)

func TestCreateCluster(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	cluster, err := c.CreateCluster("go-test", 0)
	if err != nil {
		t.Errorf("CreateCluster() failed %#v\n", err)
		return
	}
	t.Logf("%+v\n", cluster)

	for {
		c, err := c.GetCluster(cluster.Id)
		if err != nil {
			t.Errorf("GetCluster() failed %s %s\n", cluster.Id, err.Error())
			return
		}
		if c.Status == "READY" {
			break
		}
		t.Logf("%+v\n", c)
		time.Sleep(5 * time.Second)
	}
}

func TestListClusters(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	clusters, err := c.ListClusters()
	if err != nil {
		t.Errorf("ListClusters() failed %#v\n", err)
		return
	}
	t.Logf("%+v\n", clusters)
}

func TestGetCluster(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	clusters, err := c.ListClusters()
	if err != nil {
		t.Errorf("ListClusters() failed %#v\n", err)
		return
	}
	if len(clusters.Clusters) == 0 {
		t.Errorf("TestGetCluster() no clusters to list; ListClusters() returned empty list\n")
		return
	}
	cluster, err := c.GetCluster(clusters.Clusters[0].Id)
	if err != nil {
		t.Errorf("GetCluster() failed %#v\n", err)
		return
	}
	t.Logf("%+v\n", cluster)
}

func TestUploadConfig(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	clusters, err := c.ListClusters()
	if err != nil {
		t.Errorf("ListClusters() failed %#v\n", err)
		return
	}
	if len(clusters.Clusters) == 0 {
		t.Errorf("TestListConfigs() no clusters to list; ListClusters() returned empty list\n")
		return
	}
	id := ""
	for i := range clusters.Clusters {
		if clusters.Clusters[i].Name == "go-test" {
			id = clusters.Clusters[i].Id
		}
	}
	if len(id) == 0 {
		t.Errorf("TestListConfigs(); could not find cluster to delete named \"go-test\". clusters: %+v\n", clusters.Clusters)
		return
	}
	f, err := os.Open("test_data/solr_config.zip")
	if err != nil {
		t.Errorf("TestUploadConfig(); could not open zip config file \"%s\" %s\n", "test_data/solr_config.zip", err.Error())
		return
	}

	err = c.UploadConfig(id, "config-test", f)
	if err != nil {
		t.Errorf("UploadConfig() failed %#v\n", err)
		return
	}
}

func TestListConfigs(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	clusters, err := c.ListClusters()
	if err != nil {
		t.Errorf("ListClusters() failed %#v\n", err)
		return
	}
	if len(clusters.Clusters) == 0 {
		t.Errorf("TestListConfigs() no clusters to list; ListClusters() returned empty list\n")
		return
	}
	id := ""
	for i := range clusters.Clusters {
		if clusters.Clusters[i].Name == "go-test" {
			id = clusters.Clusters[i].Id
		}
	}
	if len(id) == 0 {
		t.Errorf("TestListConfigs(); could not find cluster to delete named \"go-test\". clusters: %+v\n", clusters.Clusters)
		return
	}
	configs, err := c.ListConfigs(id)
	if err != nil {
		t.Errorf("ListConfigs() failed %#v\n", err)
		return
	}
	t.Logf("%+v\n", configs)
}

func TestGetConfig(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	clusters, err := c.ListClusters()
	if err != nil {
		t.Errorf("ListClusters() failed %#v\n", err)
		return
	}
	if len(clusters.Clusters) == 0 {
		t.Errorf("TestListConfigs() no clusters to list; ListClusters() returned empty list\n")
		return
	}
	id := ""
	for i := range clusters.Clusters {
		if clusters.Clusters[i].Name == "go-test" {
			id = clusters.Clusters[i].Id
		}
	}
	if len(id) == 0 {
		t.Errorf("TestListConfigs(); could not find cluster to delete named \"go-test\". clusters: %+v\n", clusters.Clusters)
		return
	}
	data, err := c.GetConfig(id, "config-test")
	if err != nil {
		t.Errorf("GetConfig() failed %#v\n", err)
		return
	}
	if len(data) == 0 {
		t.Errorf("GetConfig() returned empty configuration file\n", err)
		return
	}
	// t.Logf("%+v\n", len(data))
}

func TestCreateCollection(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	clusters, err := c.ListClusters()
	if err != nil {
		t.Errorf("ListClusters() failed %#v\n", err)
		return
	}
	if len(clusters.Clusters) == 0 {
		t.Errorf("TestListConfigs() no clusters to list; ListClusters() returned empty list\n")
		return
	}
	id := ""
	for i := range clusters.Clusters {
		if clusters.Clusters[i].Name == "go-test" {
			id = clusters.Clusters[i].Id
		}
	}
	if len(id) == 0 {
		t.Errorf("TestListConfigs(); could not find cluster to delete named \"go-test\". clusters: %+v\n", clusters.Clusters)
		return
	}
	data, err := c.CreateCollection(id, "test_collection", "config-test", map[string]interface{}{"wt": "json"})
	if err != nil {
		t.Errorf("CreateCollection() failed %#v\n", err)
		return
	}
	t.Logf("%+v\n", string(data))
}

func TestListCollections(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	clusters, err := c.ListClusters()
	if err != nil {
		t.Errorf("ListClusters() failed %#v\n", err)
		return
	}
	if len(clusters.Clusters) == 0 {
		t.Errorf("TestListConfigs() no clusters to list; ListClusters() returned empty list\n")
		return
	}
	id := ""
	for i := range clusters.Clusters {
		if clusters.Clusters[i].Name == "go-test" {
			id = clusters.Clusters[i].Id
		}
	}
	if len(id) == 0 {
		t.Errorf("TestListConfigs(); could not find cluster to delete named \"go-test\". clusters: %+v\n", clusters.Clusters)
		return
	}
	data, err := c.ListCollections(id, map[string]interface{}{"wt": "json"})
	if err != nil {
		t.Errorf("ListCollections() failed %#v\n", err)
		return
	}
	t.Logf("%+v\n", string(data))
}

func TestUpdateCollection(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	clusters, err := c.ListClusters()
	if err != nil {
		t.Errorf("ListClusters() failed %#v\n", err)
		return
	}
	if len(clusters.Clusters) == 0 {
		t.Errorf("TestListConfigs() no clusters to list; ListClusters() returned empty list\n")
		return
	}
	id := ""
	for i := range clusters.Clusters {
		if clusters.Clusters[i].Name == "go-test" {
			id = clusters.Clusters[i].Id
		}
	}
	if len(id) == 0 {
		t.Errorf("TestListConfigs(); could not find cluster to delete named \"go-test\". clusters: %+v\n", clusters.Clusters)
		return
	}

	f, err := os.Open("test_data/data.json")
	if err != nil {
		t.Errorf("TestUpdateCollection() error: could not open zip config file \"%s\" %s\n", "test_data/solr_config.zip", err.Error())
		return
	}

	resp, err := c.Update(id, "test_collection", "application/json", f, map[string]interface{}{"wt": "json"})
	if err != nil {
		t.Errorf("UpdateIndex() failed %#v\n", err)
		return
	}

	t.Logf("%+v\n", string(resp))
}

func TestSearch(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	clusters, err := c.ListClusters()
	if err != nil {
		t.Errorf("ListClusters() failed %#v\n", err)
		return
	}
	if len(clusters.Clusters) == 0 {
		t.Errorf("TestListConfigs() no clusters to list; ListClusters() returned empty list\n")
		return
	}
	id := ""
	for i := range clusters.Clusters {
		if clusters.Clusters[i].Name == "go-test" {
			id = clusters.Clusters[i].Id
		}
	}
	if len(id) == 0 {
		t.Errorf("TestListConfigs(); could not find cluster to delete named \"go-test\". clusters: %+v\n", clusters.Clusters)
		return
	}

	resp, err := c.Search(id, "test_collection", "*", map[string]interface{}{"wt": "json", "rows": 2})
	if err != nil {
		t.Errorf("UpdateIndex() failed %#v\n", err)
		return
	}

	t.Logf("%+v\n", string(resp))
}

func TestListRankers(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	rankers, err := c.ListRankers()
	if err != nil {
		t.Errorf("ListRankers() failed %#v %s\n", err, err.Error())
		return
	}
	t.Log(rankers)
}

func TestDeleteCluster(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	clusters, err := c.ListClusters()
	if err != nil {
		t.Errorf("ListClusters() failed %#v\n", err)
		return
	}
	if len(clusters.Clusters) == 0 {
		t.Errorf("TestGetCluster() no clusters to list; ListClusters() returned empty list\n")
		return
	}
	id := ""
	for i := range clusters.Clusters {
		if clusters.Clusters[i].Name == "go-test" {
			id = clusters.Clusters[i].Id
		}
	}
	if len(id) == 0 {
		t.Errorf("TestDeleteCluster(); could not find cluster to delete named \"go-test\". clusters: %+v\n", clusters.Clusters)
		return
	}
	err = c.DeleteCluster(id)
	if err != nil {
		t.Errorf("DeleteCluster() failed %#v\n", err)
		return
	}
}
