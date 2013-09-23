$provisioning_script = <<SCRIPT
GOLANG_TARBALL=go1.2rc1.linux-amd64.tar.gz
if [ ! -d /usr/local/go ]; then
        echo "Downloading Go: $GOLANG_TARBALL...";
        wget -q http://go.googlecode.com/files/$GOLANG_TARBALL 
        echo "Extracting $GOLANG_TARBALL into /usr/local"
        tar -xzf ~vagrant/$GOLANG_TARBALL -C /usr/local

        echo "Setting up some aliases and environmental variables..."
        echo "export PATH=\$PATH:/usr/local/go/bin" >> ~vagrant/.bashrc

        echo "Fixing permissions...";
#        chown vagrant.vagrant ~vagrant/.bash_aliases
#        chown vagrant.vagrant ~vagrant/.bashrc

        echo "Cleanup.";
        rm ~vagrant/$GOLANG_TARBALL
fi;

echo "DONE!";
SCRIPT

Vagrant::Config.run do |config|
#  config.vm.boot_mode = :gui
  config.vm.box = "ubuntu-server-12.04"
  config.vm.box_url = "http://files.vagrantup.com/precise64.box"
  config.vm.network :hostonly, "10.250.250.250"
  config.vm.forward_port 80, 8080

  config.vm.share_folder( "vagrant-root", "/vagrant", ".", :extra => 'dmode=777,fmode=777' )
  config.vm.share_folder( "project-root", "/home/vagrant/project" , "~/Sites/googleGo" )
  config.vm.provision "shell", inline: $provisioning_script

end
