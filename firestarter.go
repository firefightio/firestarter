package main

import (
  "github.com/fsouza/go-dockerclient"
  "log"
  "bytes"
  "strings"
)

func logfail(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

func run(client *docker.Client, options docker.CreateContainerOptions, host_config docker.HostConfig) {
  //Create the container
  container, err := client.CreateContainer(options)
  logfail(err)

  var container_id = container.ID

  // Start the container
  err = client.StartContainer(container_id, &host_config)
  logfail(err)

  // Wait for the container to finish
  _, err = client.WaitContainer(container_id)
  logfail(err)

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
  logfail(err)

  // Print the output to the terminal
  log.Println(buf.String())

  // Remove the container
  err = client.RemoveContainer(docker.RemoveContainerOptions{
    ID: container_id,
    Force: true,
  })
  logfail(err)
}

func main() {
  // Attach to local docker instance
  endpoint := "unix:///docker.sock"
  client, err := docker.NewClient(endpoint)
  logfail(err)

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

  run(client, options, host_config)
}
