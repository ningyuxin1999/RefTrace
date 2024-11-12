from setuptools import setup, find_packages
from setuptools.dist import Distribution

class BinaryDistribution(Distribution):
    def has_ext_modules(self):
        return True

setup(
    packages=find_packages(),
    package_data={
        'reftrace.bindings': ['libreftrace.so'],
    },
    distclass=BinaryDistribution,
)