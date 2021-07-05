import pathlib
import ctypes
import os
from dataclasses import dataclass

@dataclass
class _PyspySession:
    PYSPY_LIB_NAME = "libpyspy.so"
    PYSPY_LIB_PATH = os.path.join(pathlib.Path(__file__).parent.resolve(), PYSPY_LIB_NAME)
    
    pid: int
    app_name: str
    server_address: str
    pyspy = ctypes.cdll.LoadLibrary(PYSPY_LIB_PATH)

    def start(self):
        c_app_name = ctypes.c_char_p(self.app_name.encode("UTF-8"))
        c_server_address = ctypes.c_char_p(self.server_address.encode("UTF-8"))
        c_spyname = ctypes.c_char_p("pyspy".encode("UTF-8"))
        self.pyspy.Start(c_app_name, self.pid, c_spyname, c_server_address)


class PyroscopePyspy:

    def __init__(self) -> None:
        self.sessions = []

    def start(self, app_name, pid, server_address) -> None:
        session = _PyspySession(pid, app_name, server_address)
        session.start()
        self.sessions.append(session)
