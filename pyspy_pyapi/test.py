from multiprocessing import Process
from pyspy_pyapi import start_spy
from time import sleep
import os
import signal
import threading


def work(n):
    i = 0
    while i < n:
        i += 1


def fast_function():
    k = 0
    while True:
        work(30000)
        k += 1


def slow_function():
    k = 0
    while True:
        work(50000)
        k += 1


def killer(p, timeout):
    sleep(timeout)
    #gpid = os.getpgid(p.pid)
    #print("Terminating pid: ", p.pid, " gpid: ", gpid)
    #os.killpg(gpid, signal.SIGTERM)
    print("Terminating pid: ", p.pid)
    os.kill(p.pid, signal.SIGTERM)


if __name__ == "__main__":
    pr = []
    pr.append(Process(target=fast_function))
    pr.append(Process(target=slow_function))
    
    [p.start() for p in pr]
    [threading.Thread(target=killer, args=(p, 10)).start() for p in pr]

    for p in pr:
        gpid = os.getpgid(p.pid)
        print("Started pid: ", p.pid, " gpid: ", gpid)
        
    main_pid = os.getpid()
    main_gpid = os.getpgid(p.pid)
    print("Main pid: ", main_pid, " gpid: ", main_gpid)
    
    start_spy("test name", main_pid, "http://192.168.5.16:4040")
