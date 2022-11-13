import os
from os.path import join, dirname
import re
from dotenv import load_dotenv

dotenv_path = join(dirname(__file__), '../.env')
load_dotenv(dotenv_path)

ENV = os.environ.get("ENV")


DB_USER = os.environ.get("DB_USER")
DB_PASS = os.environ.get("DB_PASS")
DB_HOST = os.environ.get("DB_HOST")
DB_PORT = os.environ.get("DB_PORT")
DB_USER = os.environ.get("DB_USER")
DB_NAME = os.environ.get("DB_NAME")
SSL_ROOT_CERT = os.environ.get("SSL_ROOT_CERT")
SSL_CERT = os.environ.get("SSL_CERT")
SSL_KEY = os.environ.get("SSL_KEY")
