from time import sleep
import os
from pyspy_pyapi.pyspy_pyapi import start_spy
import threading
from multiprocessing import Process


def work(n):
    i = 0
    while i < n:
        i += 1


def fast_function():
    k = 0
    while True:
        work(30000)
        k += 1
        if k == 10000:
            print("Done working fast")
            return


def slow_function():
    k = 0
    while True:
        work(50000)
        k += 1
        if k == 10000:
            print("Done working slow")
            return


if __name__ == "__main__":
    p = Process(target=fast_function)
    p.start()

    spy = Process(target=start_spy, args=("test name", p.pid, "http://192.168.5.16:4040"))
    spy.start()
    print("Own pid: ", os.getpid())
    print("Started spy for pid: ", p.pid)
    print("Spy pid: ", spy.pid)

    #start_spy("test name", pr[1].pid, "http://192.168.5.16:4040")

    sleep(20)
    p.join()
    p.terminate()
    spy.terminate()

