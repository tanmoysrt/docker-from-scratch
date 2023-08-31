package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {
    switch os.Args[1]{
    case "run":
        run()
    case "child":
        child()
    default:
        panic("bad command")
    }

}

func run(){
    fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())
    cmd := exec.Command("/proc/self/exe", append([]string{"child"},os.Args[2:]...)...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
        Unshareflags: syscall.CLONE_NEWNS,        
    }
    err := cmd.Start()
    if err != nil {
        fmt.Println(err)
    }
    fmt.Printf("PID %d\n",cmd.Process.Pid)
	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error waiting for command:", err)
		return
	}
}

func child(){
    fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

    // cg()

    
    syscall.Sethostname([]byte("container"))
    syscall.Chroot("/home/tanmoy/Desktop/oss/container_custom/alpine-fs")
    syscall.Chdir("/")
    syscall.Mount("proc", "proc", "proc", 0, "")

    cmd := exec.Command(os.Args[2], os.Args[3:]...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
        fmt.Println(err)
    }

    syscall.Unmount("/proc", 0)
}


func Cg(){
    cgroups := "/sys/fs/cgroup/"
    pids := filepath.Join(cgroups, "pids")
    err := os.Mkdir(filepath.Join(pids, "liz"), 0755)
    if err != nil && !os.IsExist(err){
        panic(err)
    }
    must(ioutil.WriteFile(filepath.Join(pids, "liz/pids.max"), []byte("20"), 0700))
    must(ioutil.WriteFile(filepath.Join(pids, "liz/notify_on_release"), []byte("1"), 0700))
    must(ioutil.WriteFile(filepath.Join(pids, "liz/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}

func must(err error){
    if err != nil {
        panic(err)
    }
}