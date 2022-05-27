# dut
basic dhcp/http server for DUT

To build a u-root initramfs:
u-root -defaultsh=""  ../uinit ~/go/src/github.com/u-root/cpu/cmds/cpud/ ~/go/src/github.com/u-root/u-root/cmds/core/init

To run it locally, nad have it run a cpu server remotely:

./uinit -pubkey ~/.ssh/cpu_rsa.pub  -m cpu
