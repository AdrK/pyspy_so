from setuptools import setup

setup(name="pyspy_pyapi",
        version="0.1",
        description="Provides API for Pyroscope pyspy agent",
        url="https://github.com/AdrK/pyspy_so/",
        author="Adrian Kurylak",
        author_email="adrian.kurylak@gmail.com",
        license="MIT",
        packages=["pyspy_pyapi"],
        include_package_data=True,
        package_data={'': ['./libpyspy.so', './libpyspy.h']})
