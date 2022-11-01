import os
from os.path import join, dirname
from dotenv import load_dotenv

dotenv_path = join(dirname(__file__), '../../.env')
load_dotenv(dotenv_path)

ENV = os.environ.get("ENV")
DICT_PATH = join(dirname(__file__), '../../mecab-ipadic-neologd')
