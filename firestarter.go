// Notes:
//
// - Currently run will fail if the image has not already been pulled. This is
//   desired behavior. There will eventually be an option to supply sources 
//   which can be pulled from.
//

package main

import (
  "github.com/fsouza/go-dockerclient"
  "github.com/streadway/amqp"
  "log"
  "bytes"
  "strings"
  "os"
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

func addDockerClient(endpoint string) *docker.Client {
  client, err := docker.NewClient(endpoint)
  logfail(err)
  return client
}

func main() {
  // Create a client for the local docker instance
  endpoint := "unix:///docker.sock"
  client := addDockerClient(endpoint)

  // Attach to RabbitMQ
  connection, err := amqp.Dial("amqp://" + os.Getenv("RABBITMQ_USER") + 
                     ":" + os.Getenv("RABBITMQ_PASS") + "@rabbitmq:5672")
  logfail(err)

  // Define command in string array format
  cmd := "echo\n'Hello World'"
  cmd_array := strings.Split(cmd, "\n")

  // Define container
  container_config := docker.Config{
    Image: "fedora",
    Cmd:   cmd_array,
  }

  options := docker.CreateContainerOptions{
    Config: &container_config,
  }

  // Define the container config
  host_config := docker.HostConfig{}

  // Run the container
  run(client, options, host_config)

  // Close the RabbitMQ connection before exiting
  defer connection.Close()
}
