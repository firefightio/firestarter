# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

# Source config file if it exists
CONFIG= "config.rb"

# Set logical defaults for options in config
$num_servers = 1
$vb_mem = 512
$vb_cpus = 1
$subnet = "10.100.150"
$update_channel = "alpha"

if File.exist?(CONFIG)
  require_relative CONFIG
end

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = "coreos-%s" % $update_channel
  config.vm.box_version = ">= 308.0.1"
  config.vm.box_url = "http://%s.release.core-os.net/amd64-usr/current/coreos_production_vagrant.json" % $update_channel  

  config.vm.provider "virtualbox" do |v|
    v.memory = $vb_mem
    v.cpus = $vb_cpus
  end

  # Define and provision test runners
  (1..$num_servers).each do |server_num|
    config.vm.define vm_name = "firestarter-%02d" % server_num do |server|
      # Set ip address and forward shipyard port
      server.vm.network "private_network", ip: "%s.1%02d" % [$subnet, server_num]

      # Provisioning vagrant box with docker
      server.vm.provision "docker",
        images: ["dockerfile/ubuntu", "buildingbananas/firestarter"]
    end
  end
end