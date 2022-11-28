# Podman Desktop now requires WSL2

While I cannot find a reference to the instructions from the beginning of October 2022 saying that Podman Desktop automatically installs WSL1 (not 2), it is clear that now it requires WSL2.

```
NOTE: Administrator access is required for both these prerequisites.
1. Hyper-V should be enabled
2. Windows Subsystem for Linux v2 (WSL2) should be installed.
```

It's worth noting that they have a Chocolatey install now as well.

```
choco install podman-desktop
```

* Windows \| Podman Desktop  
  https://podman-desktop.io/docs/Installation/windows-install
* Podman expands to the Desktop \| Red Hat Developer  
  https://developers.redhat.com/articles/2022/10/24/podman-expands-desktop
