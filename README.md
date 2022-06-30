# seidl

`seidl` is a lightweight [pint](https://pint.suse.com/) query utility designed for easy usage. It displays the current SUSE publiccloud images according to customizable filter rules.

In aims at complementing the [public-cloud-info-client](https://github.com/SUSE-Enceladus/public-cloud-info-client) by the feature to display all current not-deleted and not-deprecated images in a nice table on the console.

![Screenshot of seidl in action when querying the current Azure and GCP images](usage.png)

Image information are fetched from `https://susepubliccloudinfo.suse.com/v1/[microsoft|amazon|google]/images.json` and then parsed accordingly.

`seidl` is named after a small beer quantity in Austria because it's slightly smaller than a pint (typically 0.3 liters).

## Usage

    seidl -h                           # Print help
    
    seidl gce                          # Query current GCE images
    seidl --region eu-west-1 aws       # Query current AWS images
    seidl azure                        # Query current Azure images

You can also filter results using the `-f` arguments. It supports a comma separated list of strings, where every entry of that list must be present in the image name

    seidl -f sles,15-sp3 gce           # List GCE images, filter for SLES 15-SP3 images
    seidl -f sles,15-sp3 azure gce     # List Azure and GCE images, filter for SLES 15-SP3 images

Additional commands are:

    seidl --list-aws-regions           # List available AWS regions
    seidl --list-az-envs               # List possible Azure environments
    seidl --noformat ...               # Don't format the output as a table

## Install

Checkout the [Releases](https://github.com/grisu48/seidl/releases) page of this repository.

## openSUSE

`seidl` is part of the openSUSE Tumbleweed packages and can be installed there via `zypper in seidl`.

I build packages for openSUSE Leap in my [ph03nix:tools](https://build.opensuse.org/project/show/home:ph03nix:tools) repository. This is my personal "stable" repository (don't use my home repo, weird stuff happens there!).  
For install instructions follow [this link](https://software.opensuse.org//download.html?project=home%3Aph03nix%3Atools&package=seidl) or try the following, if you understand what you are doing:

    zypper addrepo https://download.opensuse.org/repositories/home:ph03nix:tools/15.4/home:ph03nix:tools.repo
    zypper refresh
    zypper install seidl

## Build

The lazy way:

    make

    make seidl install                 # the very lazy way. Build and install to ~/bin

Canonical go way:

    go build ./...

Requirements: Pure `go`, no external dependencies required here.
