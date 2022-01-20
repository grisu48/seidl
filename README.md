# seidl

`seidl` is a lightweight [pint](https://pint.suse.com/) client, designed for easy usage to query the current images of SUSE publiccloud images.

In aims at complementing the [public-cloud-info-client](https://github.com/SUSE-Enceladus/public-cloud-info-client) by the feature to display all current not-deleted and not-deprecated images in a nice table on the console.

![Screenshot of seidl in action when querying the current Google images](seidl.png)

## Usage

    ./seidl gce                          # Query current GCE images
    ./seidl aws --region eu-west-1       # Query current AWS images
    ./seidl azure                        # Query current Azure images

## Installation

This is a standalone python script with the following requirements:

    requests
    json

It should work out of the box on most machines.