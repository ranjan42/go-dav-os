/* kernel/scheduler/switch.s */

.section .text
/* Updated mangled name for github.com/dmarro89/go-dav-os/kernel/scheduler.CpuSwitch */
.global github_0com_1dmarro89_1go_x2ddav_x2dos_1kernel_1scheduler.CpuSwitch
.type   github_0com_1dmarro89_1go_x2ddav_x2dos_1kernel_1scheduler.CpuSwitch, @function

github_0com_1dmarro89_1go_x2ddav_x2dos_1kernel_1scheduler.CpuSwitch:
    pushl %ebp
    pushl %ebx
    pushl %esi
    pushl %edi

    movl 20(%esp), %eax  /* old *uint32 */
    movl 24(%esp), %edx  /* new uint32  */

    movl %esp, (%eax)
    movl %edx, %esp

    popl %edi
    popl %esi
    popl %ebx
    popl %ebp

    ret
