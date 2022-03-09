Timelapse
=========

This is a personal project that I use to create timelapses at home.

There are 2 main components: the `tlagent` which captures the pictures and pushes them over the network 
to the `tlserver`, which stores the frames and generates the timelapse video on demand.

## Timelapse agent

The timelapse agent takes the pictures that will make the timelapse video.
It runs as a daemon, and keeps taking pictures every interval of time.

Currently, the agent only works on a Raspberry Pi with the camera module, but I may expand to more devices.

```shell
tlagent capture --server-url http://192.168.0.1:8990 --interval 5m
```

## Timelapse server

The timelapse server receives the pictures from agents and stores them at a given location.
It exposes an API on port `8990`

```shell
tlserver --storage-path ~/Pictures/timelapse
```
