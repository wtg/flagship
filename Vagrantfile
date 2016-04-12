# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure(2) do |config|
  config.vm.box = "ubuntu/trusty64"
  config.vm.network "forwarded_port", guest: 5432, host: 5432

  config.vm.provision "shell", inline: <<-SHELL
    sudo apt-get update -qq
    sudo apt-get install -qqy postgresql postgresql-client postgresql-contrib libpq-dev
    sudo -u postgres psql -c "CREATE USER vagrant WITH PASSWORD 'vagrant';"
    sudo -u postgres psql -c "CREATE DATABASE vagrant OWNER vagrant;"
    sudo sed -i -e "s/#listen_addresses = 'localhost'/listen_addresses = '*'/g" /etc/postgresql/9.3/main/postgresql.conf
    echo "host all all 0.0.0.0/0 md5" | sudo tee -a /etc/postgresql/9.3/main/pg_hba.conf
    sudo service postgresql restart
  SHELL
end
