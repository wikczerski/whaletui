package docker

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
	"github.com/wikczerski/whaletui/internal/docker/utils"
)

// convertToContainers converts Docker API containers to our Container type
func (c *Client) convertToContainers(containers []container.Summary) []Container {
	result := make([]Container, 0, len(containers))
	for i := range containers {
		cont := &containers[i]
		result = append(result, c.convertContainer(cont))
	}
	return result
}

// convertContainer converts a single Docker API container to our Container type
func (c *Client) convertContainer(cont *container.Summary) Container {
	ports := utils.FormatContainerPorts(cont.Ports)
	return Container{
		ID:      cont.ID[:12],
		Name:    cont.Names[0][1:], // Remove leading slash
		Image:   cont.Image,
		Status:  cont.Status,
		State:   cont.State,
		Created: time.Unix(cont.Created, 0),
		Ports:   ports,
		Size:    "", // Size is not available in ContainerList
	}
}

// convertToImages converts Docker API images to our Image type
func (c *Client) convertToImages(images []image.Summary) []Image {
	result := make([]Image, 0, len(images))
	for i := range images {
		img := &images[i]
		result = append(result, c.convertImage(img))
	}
	return result
}

// convertImage converts a single Docker API image to our Image type
func (c *Client) convertImage(img *image.Summary) Image {
	repo, tag := utils.ParseImageRepository(img.RepoTags)
	size := utils.FormatSize(img.Size)
	return Image{
		ID:         img.ID[7:19], // Remove "sha256:" prefix and truncate
		Repository: repo,
		Tag:        tag,
		Size:       size,
		Created:    time.Unix(img.Created, 0),
		Containers: int(img.Containers),
	}
}

// createVolumeFromAPI creates a Volume from the API response
func (c *Client) createVolumeFromAPI(vol *volume.Volume) Volume {
	created := time.Time{}
	if vol.CreatedAt != "" {
		created, _ = time.Parse(time.RFC3339, vol.CreatedAt)
	}

	return Volume{
		Name:       vol.Name,
		Driver:     vol.Driver,
		Mountpoint: vol.Mountpoint,
		Created:    created,
		Size:       "", // Size is not available in VolumeList
	}
}

// readAndFormatLogs reads logs from the response and formats them
func (c *Client) readAndFormatLogs(logs io.ReadCloser) string {
	var logLines []string
	buffer := make([]byte, 1024)

	c.readLogsIntoBuffer(logs, buffer, &logLines)

	return strings.Join(logLines, "")
}

// readLogsIntoBuffer reads logs into the buffer and formats them
func (c *Client) readLogsIntoBuffer(logs io.ReadCloser, buffer []byte, logLines *[]string) {
	for {
		n, err := logs.Read(buffer)
		if n > 0 {
			line := string(buffer[:n])
			formattedLine := c.formatLogLine(line)
			*logLines = append(*logLines, formattedLine)
		}
		if err != nil {
			break
		}
	}
}

// formatLogLine formats a single log line by removing timestamp prefix if present
func (c *Client) formatLogLine(line string) string {
	if len(line) >= 8 {
		return line[8:] // Remove timestamp prefix
	}
	return line
}

// decodeStatsResponse decodes the stats response body into a map
func (c *Client) decodeStatsResponse(body io.ReadCloser) (map[string]any, error) {
	var statsJSON map[string]any
	if err := json.NewDecoder(body).Decode(&statsJSON); err != nil {
		return nil, err
	}
	return statsJSON, nil
}

// readExecOutput reads the output from the hijacked response
func (c *Client) readExecOutput(hijackedResp types.HijackedResponse) string {
	var output strings.Builder
	buffer := make([]byte, 1024)

	for {
		n, err := hijackedResp.Reader.Read(buffer)
		if n > 0 {
			output.Write(buffer[:n])
		}
		if err != nil {
			break
		}
	}

	return output.String()
}

// readServiceLogs reads and formats the service logs
func (c *Client) readServiceLogs(response io.ReadCloser) string {
	output := &strings.Builder{}
	buffer := make([]byte, 1024)

	for {
		n, err := response.Read(buffer)
		if n > 0 {
			output.Write(buffer[:n])
		}
		if err != nil {
			break
		}
	}

	return output.String()
}
