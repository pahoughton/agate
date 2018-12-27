# 2018-12-23 (cc) <paul4hough@gmail.com>
#

Vagrant.configure("2") do |config|
  # config.vm.box_check_update = false

  demo = 'demo'
  config.vm.define demo do |bcfg|
    bcfg.vm.box = "centos/7"
    # bcfg.vm.box = "ubuntu/xenial64"

    bcfg.vm.hostname = demo
    bcfg.vm.network    'private_network', ip: '10.0.0.7'
    bcfg.vm.network "forwarded_port", guest: 9090, host: 9090
    bcfg.vm.network "forwarded_port", guest: 9093, host: 9093
    bcfg.vm.network "forwarded_port", guest: 5003, host: 5003
    bcfg.vm.provider   'virtualbox' do |vb|
      vb.name      = demo
      vb.cpus      = 1
      vb.memory    = 1024
      vb.customize ['modifyvm', :id, '--natdnshostresolver1', 'on']
      vb.customize ['modifyvm', :id, '--natdnspassdomain1', 'on']
      vb.customize ['modifyvm', :id, '--usb', 'off']
    end
  end

end
