package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// DiskImage represents a DiskImage for launching instances from
type DiskImage struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Version      string `json:"version,omitempty"`
	State        string `json:"state,omitempty"`
	Distribution string `json:"distribution,omitempty"`
	Description  string `json:"description,omitempty"`
	Label        string `json:"label,omitempty"`
}

// ListDiskImages return all disk image in system
func (c *Client) ListDiskImages() ([]DiskImage, error) {
	resp, err := c.SendGetRequest("/v2/disk_images")
	if err != nil {
		return nil, decodeERROR(err)
	}

	diskImages := make([]DiskImage, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&diskImages); err != nil {
		return nil, err
	}

	return diskImages, nil
}

// GetDiskImage get one disk image using the id
func (c *Client) GetDiskImage(id string) (*DiskImage, error) {
	resp, err := c.SendGetRequest(fmt.Sprintf("/v2/disk_images/%s", id))
	if err != nil {
		return nil, decodeERROR(err)
	}

	diskImage := &DiskImage{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&diskImage); err != nil {
		return nil, err
	}

	return diskImage, nil
}

// FindDiskImage finds a disk image by either part of the ID or part of the name
func (c *Client) FindDiskImage(search string) (*DiskImage, error) {
	templateList, err := c.ListDiskImages()
	if err != nil {
		return nil, decodeERROR(err)
	}

	exactMatch := false
	partialMatchesCount := 0
	result := DiskImage{}

	for _, value := range templateList {
		if value.Name == search || value.ID == search {
			exactMatch = true
			result = value
		} else if strings.Contains(value.Name, search) || strings.Contains(value.ID, search) {
			if !exactMatch {
				result = value
				partialMatchesCount++
			}
		}
	}

	if exactMatch || partialMatchesCount == 1 {
		return &result, nil
	} else if partialMatchesCount > 1 {
		err := fmt.Errorf("unable to find %s because there were multiple matches", search)
		return nil, MultipleMatchesError.wrap(err)
	} else {
		err := fmt.Errorf("unable to find %s, zero matches", search)
		return nil, ZeroMatchesError.wrap(err)
	}
}
