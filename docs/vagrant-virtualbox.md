# Create VM with vagrant and virtualbox
## Download virtualbox and vagrant
 * Visit https://www.virtualbox.org/wiki/Downloads for virtualbox
 * Visit https://www.vagrantup.com/downloads for vagrant
 
## Create vm using vagrant
 * Edit Vagranfile using editor - ubuntu
```bash
$ vi Vagrantfile
VAGRANTFILE_API_VERSION = "2"
IMAGE_NAME = "generic/ubuntu2004"

$script = <<SCRIPT
sudo mkdir /root/.ssh
sudo chmod 600 /root/.ssh
sudo cp /vagrant/id_rsa.pub /root/.ssh/authorized_keys
sudo sed -i 's/dhcp4: .*/dhcp4: no/g' /etc/netplan/01-netcfg.yaml
sudo sed -i 's/dhcp6: .*/dhcp4: no/g' /etc/netplan/01-netcfg.yaml
sudo echo '      gateway4: 192.168.77.1' >> /etc/netplan/50-vagrant.yaml
#sudo sed -i 's/^#DNS.*/DNS=8.8.8.8/g' /etc/systemd/resolved.conf

sudo reboot
SCRIPT

$override_disk_size ||= false
$disk_size ||= "30GB"



CLUSTER = {
    "ubuntu2004-190" => { :ip => "192.168.77.190", :cpus => 4, :memory => 4096, :script => $script },
    "ubuntu2004-191" => { :ip => "192.168.77.191", :cpus => 4, :memory => 4096, :script => $script },
    "ubuntu2004-192" => { :ip => "192.168.77.192", :cpus => 4, :memory => 4096, :script => $script },
    "ubuntu2004-193" => { :ip => "192.168.77.193", :cpus => 4, :memory => 4096, :script => $script },
    "ubuntu2004-194" => { :ip => "192.168.77.194", :cpus => 4, :memory => 4096, :script => $script },
    "ubuntu2004-195" => { :ip => "192.168.77.195", :cpus => 4, :memory => 4096, :script => $script },
}

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
    CLUSTER.each_with_index do |(hostname, info), index|

        if ($override_disk_size)
            unless Vagrant.has_plugin?("vagrant-disksize")
                system "vagrant plugin install vagrant-disksize"
            end
            config.disksize.size = $disk_size
        end

        config.vm.synced_folder ".", "/vagrant", disabled: false


        config.vm.define hostname do |cfg|
            cfg.vm.provider :virtualbox do |vb, override|
                config.vm.box = IMAGE_NAME
                override.vm.network :public_network, bridge: "eno3", ip: "#{info[:ip]}"
                override.vm.hostname = hostname
                vb.name = hostname
                vb.customize ["modifyvm", :id, "--memory", info[:memory], "--cpus", info[:cpus], "--hwvirtex", "on"]
            end # end provider
            # inline shell scripts
            cfg.vm.provision :shell do |s|
                s.inline = info[:script]
            end # end inline shell scripts
        end # end config
    end # end cluster
end
```

 * Edit Vagranfile using editor - centos
```bash
$ vi Vagrantfile
# -*- mode: ruby -*- # vi: set ft=ruby :

VAGRANTFILE_API_VERSION = "2"

ENV['VAGRANT_EXPERIMENTAL'] = "disks"

$script = <<SCRIPT
sudo mkdir /root/.ssh
sudo chmod 600 /root/.ssh
sudo cp /vagrant/id_rsa.pub /root/.ssh/authorized_keys
sudo sed -i 's/^ONBOOT.*/ONBOOT=no/g' /etc/sysconfig/network-scripts/ifcfg-eth0
sudo sed -i 's/^nameserver.*/nameserver 8.8.8.8/g' /etc/resolv.conf
sudo chmod 666 /etc/sysconfig/network
sudo echo 'GATEWAY=192.168.77.1' >> /etc/sysconfig/network
sudo chmod 644 /etc/sysconfig/network
sudo chmod 666 /etc/resolv.conf
sudo echo 'nameserver 8.8.8.8' >> /etc/resolv.conf
sudo chmod 644 /etc/resolv.conf
sudo chmod 666 /etc/sysconfig/network-scripts/ifcfg-eth1
sudo echo 'DNS1=8.8.8.8' >> /etc/sysconfig/network-scripts/ifcfg-eth1
sudo chmod 644 /etc/sysconfig/network-scripts/ifcfg-eth1
sudo yum install -y cloud-utils-growpart
sudo growpart /dev/sda 1
sudo xfs_growfs /dev/sda1
sudo reboot
SCRIPT

SUPPORTED_OS = {
  "ubuntu1604"          => {box: "generic/ubuntu1604",         user: "vagrant"},
  "ubuntu1804"          => {box: "generic/ubuntu1804",         user: "vagrant"},
  "ubuntu2004"          => {box: "generic/ubuntu2004",         user: "vagrant"},
  "centos"              => {box: "centos/7",                   user: "vagrant"},
  "centos-bento"        => {box: "bento/centos-7.6",           user: "vagrant"},
  "centos8"             => {box: "centos/8",                   user: "vagrant"},
  "centos8-bento"       => {box: "bento/centos-8",             user: "vagrant"},
  "rhel7"               => {box: "generic/rhel7",              user: "vagrant"},
  "rhel8"               => {box: "generic/rhel8",              user: "vagrant"},
}

$os ||= "centos8"
$num_instances ||= 6
$instance_name_prefix ||= "node"

$vm_cpus ||= 4
$vm_memory ||= 4096
$override_disk_size ||= true
$disk_size ||= "30GB"
$vm_gui ||= false
$box = SUPPORTED_OS[$os][:box]

$subnet ||= "192.168.77"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
	config.vm.box = $box
	config.ssh.insert_key = false
	config.vm.box_check_update = false
	#config.vm.provision :shell,run:"always", inline: $script
	config.vm.provision :shell, inline: $script

	if ($override_disk_size)
		unless Vagrant.has_plugin?("vagrant-disksize")
			system "vagrant plugin install vagrant-disksize"
		end
		config.disksize.size = $disk_size
	end

	config.vm.disk :disk, size: "30GB", primary: true

	(1..$num_instances).each do |i|
		config.vm.define "#{$instance_name_prefix}#{i}" do |node|
			node.vm.network "public_network", bridge: "eno3", ip: "#{$subnet}.#{i+222}"
			node.vm.hostname = "node#{i}"
			node.vm.provider "virtualbox" do |v|
				v.cpus = $vm_cpus
				v.memory = $vm_memory
				v.name = "#{$instance_name_prefix}#{i}-#{$subnet}.#{i+222}"
			end
		end
	end
end```

```bash
$ vagrant up
```