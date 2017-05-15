from fabric.api import *
from fabric.contrib.files import exists
from fabric.state import output
from fabric.contrib.project import rsync_project
from fabric.contrib.console import confirm
import yaml

output['running'] = False

env.hosts = ["ubuntu@lantrn.xyz"]

def deploy():
    build()
    upload()
    print ("Done.")

def build():
    print ("Building...")
    local("mkdir -p build")
    local("env GOOS=linux GOARCH=amd64 go build -o build/api .")

def upload():
    print ("Uploading...")
    uploading_directory = "/etc/lantern"

    run("mkdir -p "+uploading_directory)

    with settings(hide('warnings', 'running', 'stdout')):
        rsync_project(remote_dir=uploading_directory, local_dir="build/api", delete=True)