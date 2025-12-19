# go-dav-os
Hobby project to dig deeper into how an OS works by writing the bare essentials of a kernel in Go. Only the kernel lives here (gccgo, 32-bit protected mode); BIOS and bootloader are handled by battle-tested tools (GRUB with a Multiboot header). No reinvention of those pieces.

## What’s inside
- Boot: `boot/boot.s` exposes the Multiboot header and `_start`, sets up a 16 KB stack, disables interrupts, and jumps into `kernel.Main`.
- Kernel: `kernel/` in Go, freestanding build with gccgo. Minimal IDT with stubs for #GP, #DF, and a test handler on `int 0x80`.
- Terminal: `terminal/` writes to VGA text mode 80x25, manages cursor and scroll.
- Keyboard: `keyboard/` reads from PS/2 and maps keys with the Italian layout.
- Tiny shell: interactive prompt with `help`, `clear`, `about`, plus a test `TriggerInt80()` at boot to show the return path.

## Project status
- Experimental, single-core, no paging, no scheduler, no filesystem, no “real” drivers yet.
- Runs in 32-bit protected mode, meant for QEMU/GRUB. No UEFI.
- Go runtime pared down: stubs for the functions gccgo expects in freestanding mode (write barrier, panic handlers, etc.).

## Dependencies
- Via Docker (recommended): Docker with `--platform=linux/amd64`.
- Native (if you want to do it manually): cross toolchain `i686-elf-{binutils,gccgo}`, `grub-mkrescue`, `xorriso`, `mtools`, `qemu-system-i386`.

## Build and run (Docker)
```bash
docker build --platform=linux/amd64 -t go-dav-os-toolchain .
docker run --rm --platform=linux/amd64 \
  -v "$PWD":/work -w /work go-dav-os-toolchain \
  make            # builds build/dav-go-os.iso
qemu-system-i386 -cdrom build/dav-go-os.iso
```
Quick targets from the Makefile:
- `make docker-build-only` builds the image and the ISO.
- `make run` (outside Docker) runs QEMU on an existing ISO.

## Build natively
Assuming an `i686-elf-*` toolchain is installed:
```bash
make
qemu-system-i386 -cdrom build/dav-go-os.iso
```
To force cross binaries: `make CROSS=i686-elf`.

## What you’ll see on screen
- On boot a test `int 0x80` fires, then the prompt `> ` shows up.
- `help` lists commands; `clear` wipes the screen; `about` prints kernel info.
- On #GP or #DF, the ASM stubs drop a character to VGA and halt.

## Folder layout
- `boot/`: Multiboot header, `_start`, gccgo runtime stubs, IDT hooks.
- `kernel/`: main logic and IDT setup.
- `terminal/`: VGA text driver.
- `keyboard/`: PS/2 input with IT layout.
- `iso/`: `grub.cfg` for the ISO.

## Final note
Personal, open-source, work-in-progress. If you try it or contribute, any feedback is welcome. I’m building pieces as I learn them—the goal is understanding, not chasing modern-OS feature lists.
