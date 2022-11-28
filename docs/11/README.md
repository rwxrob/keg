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

I did find this reference to WSL1 or WSL2:

```
How to Run Podman on Windows {via WSL and Linux} - Knowledge Base by ...
Prerequisites. A machine running Windows 10 or 11 (For other systems, see how
to Install Podman on Ubuntu, Install Podman on macOS).; Administrator
privileges. Install Linux on Windows. Installing Linux distributions on Windows
requires the WSL (WSL1 or WSL2) feature. In this tutorial, we will use WSL2
because it offers superior performance compared to WSL1 and full system
compatibility. phoenixnap.com/kb/podman-windows
```

Maybe I was mistaken.

* Windows \| Podman Desktop  
  https://podman-desktop.io/docs/Installation/windows-install
* Podman expands to the Desktop \| Red Hat Developer  
  https://developers.redhat.com/articles/2022/10/24/podman-expands-desktop
* Podman in Action (Early access) \| Red Hat Developer  
  https://developers.redhat.com/e-books/podman-action-early-access
* Advancing the next phase of container development with Podman Desktop  
  https://www.redhat.com/en/blog/advancing-next-phase-container-development-podman-desktop
