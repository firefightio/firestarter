package main

import (
  "github.com/fsouza/go-dockerclient"
  "log"
  "bytes"
  "strings"
)

func main() {
  // Attach to local docker instance
  endpoint := "unix:///docker.sock"
  client, err := docker.NewClient(endpoint)
  if err != nil {
    log.Fatal(err)
  }

  // Define command in string array format
  cmd := "echo\n'Hello World'"
  cmd_array := strings.Split(cmd, "\n")

  // Define container
  container_config := docker.Config{
    Image: "dockerfile/ubuntu",
    Cmd:   cmd_array,
  }

  options := docker.CreateContainerOptions{
    Config: &container_config,
  }

  // Define the container config
  host_config := docker.HostConfig{}

  // ------------------------------------------------------------------
  // This section will eventually be it's own goroutine
  // ------------------------------------------------------------------

  //Create the container
  container, err := client.CreateContainer(options)
  if err != nil {
    log.Fatal(err)
  }

  var container_id = container.ID

  // Start the container
  err = client.StartContainer(container_id, &host_config)
  if err != nil {
    log.Fatal(err)
  }

  // Wait for the container to finish
  _, err = client.WaitContainer(container_id)
  if err != nil {
    log.Fatal(err)
  }

  // Create a stream to read from docker container
  var buf bytes.Buffer

  // Read the log from the container
  err = client.AttachToContainer(docker.AttachToContainerOptions{
    Container:    container_id,
    OutputStream: &buf,
    Logs:         true,
    Stdout:       true,
    Stderr:       true,
  })
  
  if err != nil {
    log.Fatal(err)
  }

  // Print the output to the terminal
  log.Println(buf.String())

  // Remove the container
  err = client.RemoveContainer(docker.RemoveContainerOptions{
    ID: container_id,
    Force: true,
  })
  
  if err != nil {
    log.Fatal(err)
  }
}