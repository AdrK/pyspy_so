import pathlib
import ctypes
import os

pyspy_lib_name = "libpyspy.so"
pyspy_lib_path = os.path.join(pathlib.Path(__file__).parent.resolve(), pyspy_lib_name)
pyspy = ctypes.cdll.LoadLibrary(pyspy_lib_path)

def start_spy(app_name, pid, server_address):
    print("Loading library from path: ", pyspy_lib_path)
    b_app_name = app_name.encode("UTF-8")
    b_server_address = server_address.encode("UTF-8")
    b_spyname = "pyspy".encode("UTF-8")
    pyspy.Start(ctypes.c_char_p(b_app_name), pid, ctypes.c_char_p(b_spyname), ctypes.c_char_p(b_server_address))

