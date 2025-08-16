package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServiceFactory(t *testing.T) {
	factory := NewServiceFactory(nil)

	assert.NotNil(t, factory)
	assert.NotNil(t, factory.ContainerService)
	assert.NotNil(t, factory.ImageService)
	assert.NotNil(t, factory.VolumeService)
	assert.NotNil(t, factory.NetworkService)
	assert.NotNil(t, factory.DockerInfoService)
}

func TestServiceFactory_ServiceInstances(t *testing.T) {
	factory := NewServiceFactory(nil)

	containerService1 := factory.ContainerService
	containerService2 := factory.ContainerService

	assert.Equal(t, containerService1, containerService2)
	assert.NotEqual(t, factory.ContainerService, factory.ImageService)
	assert.NotEqual(t, factory.ContainerService, factory.VolumeService)
	assert.NotEqual(t, factory.ContainerService, factory.NetworkService)
	assert.NotEqual(t, factory.ContainerService, factory.DockerInfoService)
}

func TestServiceFactory_NilClient(t *testing.T) {
	factory := NewServiceFactory(nil)

	assert.NotNil(t, factory)
	assert.NotNil(t, factory.ContainerService)
	assert.NotNil(t, factory.ImageService)
	assert.NotNil(t, factory.VolumeService)
	assert.NotNil(t, factory.NetworkService)
	assert.NotNil(t, factory.DockerInfoService)
}

func TestServiceFactory_ServiceCreation(t *testing.T) {
	factory := NewServiceFactory(nil)

	containerService := factory.ContainerService
	assert.NotNil(t, containerService)

	imageService := factory.ImageService
	assert.NotNil(t, imageService)

	volumeService := factory.VolumeService
	assert.NotNil(t, volumeService)

	networkService := factory.NetworkService
	assert.NotNil(t, networkService)

	dockerInfoService := factory.DockerInfoService
	assert.NotNil(t, dockerInfoService)
}

func TestServiceFactory_InterfaceCompliance(t *testing.T) {
	factory := NewServiceFactory(nil)

	containerService := factory.ContainerService
	assert.NotNil(t, containerService)

	imageService := factory.ImageService
	assert.NotNil(t, imageService)

	volumeService := factory.VolumeService
	assert.NotNil(t, volumeService)

	networkService := factory.NetworkService
	assert.NotNil(t, networkService)

	dockerInfoService := factory.DockerInfoService
	assert.NotNil(t, dockerInfoService)
}
