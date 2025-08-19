package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServiceFactory(t *testing.T) {
	factory := NewServiceFactory(nil)

	assert.NotNil(t, factory)
	assert.Nil(t, factory.ContainerService)
	assert.Nil(t, factory.ImageService)
	assert.Nil(t, factory.VolumeService)
	assert.Nil(t, factory.NetworkService)
	assert.Nil(t, factory.DockerInfoService)
}

func TestServiceFactory_ServiceInstances(t *testing.T) {
	factory := NewServiceFactory(nil)

	containerService1 := factory.ContainerService
	containerService2 := factory.ContainerService

	assert.Equal(t, containerService1, containerService2)
	assert.Equal(t, factory.ContainerService, factory.ImageService)      // Both should be nil
	assert.Equal(t, factory.ContainerService, factory.VolumeService)     // All should be nil
	assert.Equal(t, factory.ContainerService, factory.NetworkService)    // All should be nil
	assert.Equal(t, factory.ContainerService, factory.DockerInfoService) // All should be nil
}

func TestServiceFactory_NilClient(t *testing.T) {
	factory := NewServiceFactory(nil)

	assert.NotNil(t, factory)
	assert.Nil(t, factory.ContainerService)
	assert.Nil(t, factory.ImageService)
	assert.Nil(t, factory.VolumeService)
	assert.Nil(t, factory.NetworkService)
	assert.Nil(t, factory.DockerInfoService)
}

func TestServiceFactory_ServiceCreation(t *testing.T) {
	factory := NewServiceFactory(nil)

	containerService := factory.ContainerService
	assert.Nil(t, containerService)

	imageService := factory.ImageService
	assert.Nil(t, imageService)

	volumeService := factory.VolumeService
	assert.Nil(t, volumeService)

	networkService := factory.NetworkService
	assert.Nil(t, networkService)

	dockerInfoService := factory.DockerInfoService
	assert.Nil(t, dockerInfoService)
}

func TestServiceFactory_InterfaceCompliance(t *testing.T) {
	factory := NewServiceFactory(nil)

	containerService := factory.ContainerService
	assert.Nil(t, containerService)

	imageService := factory.ImageService
	assert.Nil(t, imageService)

	volumeService := factory.VolumeService
	assert.Nil(t, volumeService)

	networkService := factory.NetworkService
	assert.Nil(t, networkService)

	dockerInfoService := factory.DockerInfoService
	assert.Nil(t, dockerInfoService)
}
