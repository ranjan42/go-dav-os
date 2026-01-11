import subprocess
import time
import sys
import os

def check_log_for(target, log_file, timeout=5):
    start = time.time()
    while time.time() - start < timeout:
        if os.path.exists(log_file):
            with open(log_file, 'r', errors='ignore') as f:
                content = f.read()
                if target in content:
                    return True
        time.sleep(0.5)
    return False

def main():
    iso_path = "build/dav-go-os.iso"
    log_file = "qemu.log"
    
    if os.path.exists(log_file):
        os.remove(log_file)
        
    print(f"Starting QEMU verification for {iso_path}...")
    
    # Run QEMU with debugcon logging to file and monitor on stdio
    cmd = [
        "qemu-system-i386",
        "-cdrom", iso_path,
        "-debugcon", f"file:{log_file}",
        "-monitor", "stdio",
        "-display", "none",
        "-nographic" 
    ]
    
    # Start QEMU process
    process = subprocess.Popen(
        cmd,
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE, # Capture stdout (monitor output) to avoid clutter
        stderr=subprocess.PIPE,
        start_new_session=True # Detach group
    )
    
    try:
        # 1. Wait for Boot Prompt "DavOS"
        print("Waiting for boot prompt...")
        if not check_log_for("DavOS", log_file, timeout=10):
            print("ERROR: Timeout waiting for DavOS prompt.")
            sys.exit(1)
        print("Boot successful.")
        
        # 2. Test Shell Command "help"
        # Since we are using -monitor stdio, writing to stdin interacts with the QEMU monitor, NOT the guest OS serial/keyboard.
        # To send keys to the guest, we use the monitor command 'sendkey'.
        
        print("Sending 'help' command via QEMU monitor...")
        # Send keys: h, e, l, p, ret
        keys = ['h', 'e', 'l', 'p', 'ret']
        for k in keys:
            cmd_str = f"sendkey {k}\n"
            process.stdin.write(cmd_str.encode())
            process.stdin.flush()
            time.sleep(0.1) # Debounce
            
        # 3. Wait for Command Output "Commands:"
        print("Waiting for help output...")
        if not check_log_for("Commands:", log_file, timeout=5):
            print("ERROR: Timeout waiting for help command text.")
            sys.exit(1)
            
        print("Test Passed: 'help' command executed successfully.")
        
    finally:
        process.kill()

if __name__ == "__main__":
    main()
